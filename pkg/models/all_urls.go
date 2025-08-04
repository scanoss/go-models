// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2018-2025 SCANOSS.COM
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

package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"

	"github.com/scanoss/go-grpc-helper/pkg/grpc/database"
)

// AllUrlsModel provides database access for URL information.
type AllUrlsModel struct {
	q *database.DBQueryContext
}

// AllURL represents a row on the AllURL table
type AllURL struct {
	Component string `db:"component"`
	Version   string `db:"version"`
	SemVer    string `db:"semver"`
	License   string `db:"license"`
	LicenseID int32  `db:"license_id"`
	IsSpdx    bool   `db:"is_spdx"`
	PurlName  string `db:"purl_name"`
	MineID    int32  `db:"mine_id"`
	URL       string `db:"-"` // Computed field, not from database
}

// NewAllURLModel creates a new instance of the AllUrlsModel.
func NewAllURLModel(q *database.DBQueryContext) *AllUrlsModel {
	return &AllUrlsModel{
		q: q,
	}
}

// GetURLsByPurlNameType retrieves all component URLs matching the specified PURL name and type.
func (m *AllUrlsModel) GetURLsByPurlNameType(ctx context.Context, purlName, purlType string) ([]AllURL, error) {
	s := ctxzap.Extract(ctx).Sugar()

	if len(purlName) == 0 {
		s.Error("Please specify a valid Purl Name to query")
		return nil, errors.New("please specify a valid Purl Name to query")
	}
	if len(purlType) == 0 {
		s.Errorf("Please specify a valid Purl Type to query: %v", purlName)
		return nil, errors.New("please specify a valid Purl Type to query")
	}

	query := "SELECT component, v.version_name AS version, v.semver AS semver," +
		" l.license_name AS license, l.is_spdx AS is_spdx, u.license_id," +
		" purl_name, mine_id FROM all_urls u" +
		" LEFT JOIN mines m ON u.mine_id = m.id" +
		" LEFT JOIN licenses l ON u.license_id = l.id" +
		" LEFT JOIN versions v ON u.version_id = v.id" +
		" WHERE m.purl_type = $1 AND u.purl_name = $2 ORDER BY date DESC"

	var allUrls []AllURL
	err := m.q.SelectContext(ctx, &allUrls, query, purlType, purlName)
	if err != nil {
		s.Errorf("Failed to query all urls table for %v - %v: %v", purlType, purlName, err)
		return nil, fmt.Errorf("failed to query the all urls table: %v", err)
	}

	s.Debugf("Found %v results for %v, %v.", len(allUrls), purlType, purlName)
	return allUrls, nil
}

// GetURLsByPurlNameTypeVersion retrieves component URLs for a specific PURL name, type, and version.
// Returns all matching results for the exact version.
func (m *AllUrlsModel) GetURLsByPurlNameTypeVersion(ctx context.Context, purlName, purlType, purlVersion string) ([]AllURL, error) {
	s := ctxzap.Extract(ctx).Sugar()
	if len(purlName) == 0 {
		s.Error("Please specify a valid Purl Name to query")
		return nil, errors.New("please specify a valid Purl Name to query")
	}
	if len(purlType) == 0 {
		s.Error("Please specify a valid Purl Type to query")
		return nil, errors.New("please specify a valid Purl Type to query")
	}
	if len(purlVersion) == 0 {
		s.Error("Please specify a valid Purl Version to query")
		return nil, errors.New("please specify a valid Purl Version to query")
	}

	//This query is same as GetURLsByPurlNameType but adds a WHERE clause for versions
	query := "SELECT component, v.version_name AS version, v.semver AS semver," +
		" l.license_name AS license, l.is_spdx AS is_spdx, u.license_id," +
		" purl_name, mine_id FROM all_urls u" +
		" LEFT JOIN mines m ON u.mine_id = m.id" +
		" LEFT JOIN licenses l ON u.license_id = l.id" +
		" LEFT JOIN versions v ON u.version_id = v.id" +
		" WHERE m.purl_type = $1 AND u.purl_name = $2 AND v.version_name = $3 ORDER BY date DESC"

	var allUrls []AllURL
	err := m.q.SelectContext(ctx, &allUrls, query, purlType, purlName, purlVersion)
	if err != nil {
		s.Errorf("Failed to query all urls table for %v - %v - %v: %v", purlType, purlName, purlVersion, err)
		return nil, fmt.Errorf("failed to query the all urls table: %v", err)
	}

	s.Debugf("Found %v results for %v, %v, %v.", len(allUrls), purlType, purlName, purlVersion)
	return allUrls, nil
}
