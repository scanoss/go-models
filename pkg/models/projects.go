// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2018-2026 SCANOSS.COM
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
	"database/sql"
	"errors"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/jmoiron/sqlx"
)

type ProjectModel struct {
	db *sqlx.DB
}

type Project struct {
	MineID              int32   `db:"mine_id"`
	PurlName            string  `db:"purl_name"`
	PurlType            string  `db:"purl_type"`
	Vendor              string  `db:"vendor"`
	Component           string  `db:"component"`
	License             string  `db:"license"`
	LicenseID           string  `db:"license_id"`
	IsSpdx              bool    `db:"is_spdx"`
	GitLicense          string  `db:"g_license"`
	GitLicenseID        string  `db:"g_license_id"`
	GitIsSpdx           bool    `db:"g_is_spdx"`
	SourceMineID        *int32  `db:"source_mine_id"`
	SourcePurlName      *string `db:"source_purl_name"`
	SourceVendor        *string `db:"source_vendor"`
	SourceComponent     *string `db:"source_component"`
	SourceMineName      string  `db:"source_mine_name"`
	SourcePurlType      string  `db:"source_purl_type"`
	SourceRepositoryURL string  `db:"source_repository_url"`
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
		"SELECT p.mine_id, p.purl_name,"+
			" COALESCE(m.purl_type, '')       AS purl_type,"+
			" COALESCE(p.vendor, '')          AS vendor,"+
			" COALESCE(p.component, '')       AS component,"+
			" COALESCE(l.license_name, '')    AS license,"+
			" COALESCE(l.spdx_id, '')         AS license_id,"+
			" COALESCE(l.is_spdx, false)      AS is_spdx,"+
			" COALESCE(g.license_name, '')    AS g_license,"+
			" COALESCE(g.spdx_id, '')         AS g_license_id,"+
			" COALESCE(g.is_spdx, false)      AS g_is_spdx,"+
			" p.source_mine_id,"+
			" p.source_purl_name,"+
			" p.source_vendor,"+
			" p.source_component,"+
			" COALESCE(sm.mine_name, '')      AS source_mine_name,"+
			" COALESCE(sm.purl_type, '')     AS source_purl_type,"+
			" COALESCE(sm.repository_url, '') AS source_repository_url"+
			" FROM projects p"+
			" LEFT JOIN mines m    ON p.mine_id        = m.id"+
			" LEFT JOIN mines sm   ON p.source_mine_id = sm.id"+
			" LEFT JOIN licenses l ON p.license_id     = l.id"+
			" LEFT JOIN licenses g ON p.git_license_id = g.id"+
			" WHERE m.purl_type = $1 AND p.purl_name = $2",
		purlType, purlName)
	if err != nil {
		s.Errorf("Failed to query projects table for %v, %v: %v", purlName, purlType, err)
		return nil, fmt.Errorf("failed to query the projects table: %v", err)
	}
	return allProjects, nil
}

// GetProjectByPurlNameAndMineID searches the projects' table for details about a Purl Name and Mine ID.
func (m *ProjectModel) GetProjectByPurlNameAndMineID(ctx context.Context, purlName string, mineID int32) (Project, error) {
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
		"SELECT p.mine_id, p.purl_name,"+
			" COALESCE(m.purl_type, '')       AS purl_type,"+
			" COALESCE(p.vendor, '')          AS vendor,"+
			" COALESCE(p.component, '')       AS component,"+
			" COALESCE(l.license_name, '')    AS license,"+
			" COALESCE(l.spdx_id, '')         AS license_id,"+
			" COALESCE(l.is_spdx, false)      AS is_spdx,"+
			" COALESCE(g.license_name, '')    AS g_license,"+
			" COALESCE(g.spdx_id, '')         AS g_license_id,"+
			" COALESCE(g.is_spdx, false)      AS g_is_spdx,"+
			" p.source_mine_id,"+
			" p.source_purl_name,"+
			" p.source_vendor,"+
			" p.source_component,"+
			" COALESCE(sm.mine_name, '')      AS source_mine_name,"+
			" COALESCE(sm.purl_type, '')      AS source_purl_type,"+
			" COALESCE(sm.repository_url, '') AS source_repository_url"+
			" FROM projects p"+
			" LEFT JOIN mines m    ON p.mine_id        = m.id"+
			" LEFT JOIN mines sm   ON p.source_mine_id = sm.id"+
			" LEFT JOIN licenses l ON p.license_id     = l.id"+
			" LEFT JOIN licenses g ON p.git_license_id = g.id"+
			" WHERE p.purl_name = $1 AND p.mine_id = $2",
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

// GetProjectByPurlName searches the projects' table for a single project matching
// the given Purl Name and Purl Type (resolved via the mines join). Returns
// sql.ErrNoRows when no match exists.
func (m *ProjectModel) GetProjectByPurlName(ctx context.Context, purlName string, purlType string) (Project, error) {
	s := ctxzap.Extract(ctx).Sugar()
	if len(purlName) == 0 {
		s.Error("Please specify a valid Purl Name to query")
		return Project{}, errors.New("please specify a valid Purl Name to query")
	}
	if len(purlType) == 0 {
		s.Error("Please specify a valid Purl Type to query")
		return Project{}, errors.New("please specify a valid Purl Type to query")
	}
	var project Project
	err := m.db.GetContext(ctx, &project,
		"SELECT p.mine_id, p.purl_name,"+
			" COALESCE(m.purl_type, '')       AS purl_type,"+
			" COALESCE(p.vendor, '')          AS vendor,"+
			" COALESCE(p.component, '')       AS component,"+
			" COALESCE(l.license_name, '')    AS license,"+
			" COALESCE(l.spdx_id, '')         AS license_id,"+
			" COALESCE(l.is_spdx, false)      AS is_spdx,"+
			" COALESCE(g.license_name, '')    AS g_license,"+
			" COALESCE(g.spdx_id, '')         AS g_license_id,"+
			" COALESCE(g.is_spdx, false)      AS g_is_spdx,"+
			" p.source_mine_id,"+
			" p.source_purl_name,"+
			" p.source_vendor,"+
			" p.source_component,"+
			" COALESCE(sm.mine_name, '')      AS source_mine_name,"+
			" COALESCE(sm.purl_type, '')      AS source_purl_type,"+
			" COALESCE(sm.repository_url, '') AS source_repository_url"+
			" FROM projects p"+
			" LEFT JOIN mines m    ON p.mine_id        = m.id"+
			" LEFT JOIN mines sm   ON p.source_mine_id = sm.id"+
			" LEFT JOIN licenses l ON p.license_id     = l.id"+
			" LEFT JOIN licenses g ON p.git_license_id = g.id"+
			" WHERE m.purl_type = $1 AND p.purl_name = $2"+
			" LIMIT 1",
		purlType, purlName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Project{}, err
		}
		s.Errorf("Failed to query projects table for %v, %v: %v", purlName, purlType, err)
		return Project{}, fmt.Errorf("failed to query the projects table: %v", err)
	}
	return project, nil
}
