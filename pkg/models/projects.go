// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2018-2022 SCANOSS.COM
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

// Handle all interaction with the projects table

package models

import (
	"context"
	"errors"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/jmoiron/sqlx"
)

type ProjectModel struct {
	db *sqlx.DB
}

type Project struct {
	PurlName     string `db:"purl_name"`
	Component    string `db:"component"`
	License      string `db:"license"`
	LicenseID    string `db:"license_id"`
	IsSpdx       bool   `db:"is_spdx"`
	GitLicense   string `db:"g_license"`
	GitLicenseID string `db:"g_license_id"`
	GitIsSpdx    bool   `db:"g_is_spdx"`
}

// SourcePurlRow is the raw row returned by the source-PURL lookup query:
// a join across projects and mines that exposes the source-mine fields
// needed to assemble a source PURL.
type SourcePurlRow struct {
	SourceMineID   int32  `db:"source_mine_id"`
	SourcePurlName string `db:"source_purl_name"`
	SourceVendor   string `db:"source_vendor"`
	MineName       string `db:"mine_name"`
	PurlType       string `db:"purl_type"`
	RepositoryURL  string `db:"repository_url"`
}

// NewProjectModel creates a new instance of the Project Model.
func NewProjectModel(db *sqlx.DB) *ProjectModel {
	return &ProjectModel{db: db}
}

// GetProjectsByPurlName searches the projects' table for details about Purl Name and Type.
func (m *ProjectModel) GetProjectsByPurlName(ctx context.Context, purlName string, purlType string) ([]Project, error) {
	s := ctxzap.Extract(ctx).Sugar()
	if len(purlName) == 0 {
		s.Error("Please specify a valid Purl Name to query")
		return nil, errors.New("please specify a valid Purl Name to query")
	}
	if len(purlType) == 0 {
		s.Error("Please specify a valid Purl Type to query")
		return nil, errors.New("please specify a valid Purl Type to query")
	}
	var allProjects []Project
	err := m.db.SelectContext(ctx, &allProjects,
		"SELECT purl_name, component,"+
			" l.license_name AS   license, l.spdx_id AS   license_id, l.is_spdx AS   is_spdx,"+
			" g.license_name AS g_license, g.spdx_id AS g_license_id, g.is_spdx AS g_is_spdx"+
			" FROM projects p"+
			" LEFT JOIN mines m ON p.mine_id = m.id"+
			" LEFT JOIN licenses l ON p.license_id = l.id"+
			" LEFT JOIN licenses g ON p.git_license_id = g.id"+
			" WHERE m.purl_type = $1 AND p.purl_name = $2",
		purlType, purlName)
	if err != nil {
		s.Errorf("Failed to query projects table for %v, %v: %v", purlName, purlType, err)
		return nil, fmt.Errorf("failed to query the projects table: %v", err)
	}
	return allProjects, nil
}

// GetSourcePurl looks up the source-mine row used to build a source PURL
// for a component identified by (purlName, purlType). It returns an empty
// row (and nil error) when no match exists.
func (m *ProjectModel) GetSourcePurl(ctx context.Context, purlName string, purlType string) (SourcePurlRow, error) {
	s := ctxzap.Extract(ctx).Sugar()
	if len(purlName) == 0 {
		s.Error("Please specify a valid Purl Name to query")
		return SourcePurlRow{}, errors.New("please specify a valid Purl Name to query")
	}
	if len(purlType) == 0 {
		s.Error("Please specify a valid Purl Type to query")
		return SourcePurlRow{}, errors.New("please specify a valid Purl Type to query")
	}
	rows, err := m.db.QueryxContext(ctx,
		"SELECT p.source_mine_id,"+
			" COALESCE(p.source_purl_name, '') AS source_purl_name,"+
			" COALESCE(p.source_vendor, '')    AS source_vendor,"+
			" COALESCE(sm.mine_name, '')       AS mine_name,"+
			" COALESCE(sm.purl_type, '')       AS purl_type,"+
			" COALESCE(sm.repository_url, '')  AS repository_url"+
			" FROM projects p"+
			" INNER JOIN mines m  ON p.mine_id        = m.id"+
			" INNER JOIN mines sm ON p.source_mine_id = sm.id"+
			" WHERE p.mine_id        IS NOT NULL"+
			"   AND p.source_mine_id IS NOT NULL"+
			"   AND p.purl_name      = $1"+
			"   AND m.purl_type      = $2"+
			" LIMIT 1",
		purlName, purlType)
	defer func() {
		if rows != nil {
			closeErr := rows.Close()
			if closeErr != nil {
				s.Warnf("Problem closing Rows: %v", closeErr)
			}
		}
	}()
	if err != nil {
		s.Errorf("Failed to query source purl for %v, %v: %v", purlName, purlType, err)
		return SourcePurlRow{}, fmt.Errorf("failed to query the projects table: %v", err)
	}
	var row SourcePurlRow
	if rows.Next() {
		if err := rows.StructScan(&row); err != nil {
			s.Errorf("Failed to parse source purl row for %v, %v: %v", purlName, purlType, err)
			return SourcePurlRow{}, fmt.Errorf("failed to parse source purl row: %v", err)
		}
	}
	return row, nil
}

// GetProjectByPurlName searches the projects' table for details about a Purl Name and Mine ID.
func (m *ProjectModel) GetProjectByPurlName(ctx context.Context, purlName string, mineID int32) (Project, error) {
	s := ctxzap.Extract(ctx).Sugar()
	if len(purlName) == 0 {
		s.Error("Please specify a valid Purl Name to query")
		return Project{}, errors.New("please specify a valid Purl Name to query")
	}
	if mineID < 0 {
		s.Error("Please specify a valid Mine ID to query")
		return Project{}, errors.New("please specify a valid Mine ID to query")
	}
	rows, err := m.db.QueryxContext(ctx,
		"SELECT purl_name, component,"+
			" l.license_name AS   license, l.spdx_id AS   license_id, l.is_spdx AS   is_spdx,"+
			" g.license_name AS g_license, g.spdx_id AS g_license_id, g.is_spdx AS g_is_spdx"+
			" FROM projects p"+
			" LEFT JOIN licenses l ON p.license_id = l.id"+
			" LEFT JOIN licenses g ON p.git_license_id = g.id"+
			" WHERE purl_name = $1 AND mine_id = $2",
		purlName, mineID)

	defer func() {
		if rows != nil {
			closeErr := rows.Close()
			if closeErr != nil {
				s.Warnf("Problem closing Rows: %v", closeErr)
			}
		}
	}()

	if err != nil {
		s.Errorf("Error: Failed to query projects table for %v, %v: %v", purlName, mineID, err)
		return Project{}, fmt.Errorf("failed to query the projects table: %v", err)
	}
	var project Project
	if rows.Next() {
		err = rows.StructScan(&project)
		if err != nil {
			s.Errorf("Failed to parse projects table results for %#v: %v", rows, err)
			s.Errorf("Query failed for purl_name = %v, mine_id = %v", purlName, mineID)
			return Project{}, fmt.Errorf("failed to query the projects table: %v", err)
		}
	}
	return project, nil
}
