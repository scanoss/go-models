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

package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/scanoss/go-models/pkg/models"
	"github.com/scanoss/go-models/pkg/types"
	purlutils "github.com/scanoss/go-purl-helper/pkg"
)

// ErrProjectNotFound is returned when no project row is found for the given PURL.
var ErrProjectNotFound = errors.New("project not found")

// ProjectService exposes project-table lookups.
type ProjectService struct {
	models *models.Models
}

// NewProjectService creates a new ProjectService instance.
func NewProjectService(models *models.Models) *ProjectService {
	return &ProjectService{models: models}
}

// GetProject retrieves a project row from the projects table for the given
// PURL, joining mines (for purl_type and source-mine details) and licenses.
// It is intended as a fallback when no resolved component is available.
// Returns ErrProjectNotFound if no project matches the (purl_name, purl_type)
// pair. When multiple projects match, the first row is returned.
func (ps *ProjectService) GetProject(ctx context.Context, purl string) (types.Project, error) {
	if len(purl) == 0 {
		return types.Project{}, errors.New("please specify a valid purl to query")
	}
	packageURL, err := purlutils.PurlFromString(purl)
	if err != nil {
		return types.Project{}, fmt.Errorf("failed to parse purl: %w", err)
	}
	purlName, err := purlutils.PurlNameFromString(purl)
	if err != nil {
		return types.Project{}, fmt.Errorf("failed to extract purl name: %w", err)
	}
	row, err := ps.models.Projects.GetProjectByPurlName(ctx, purlName, packageURL.Type)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.Project{}, ErrProjectNotFound
		}
		return types.Project{}, err
	}
	return types.Project{
		MineID:              row.MineID,
		PurlName:            row.PurlName,
		PurlType:            row.PurlType,
		Vendor:              row.Vendor,
		Component:           row.Component,
		License:             row.License,
		LicenseID:           row.LicenseID,
		SourceMineID:        row.SourceMineID,
		SourcePurlName:      row.SourcePurlName,
		SourceVendor:        row.SourceVendor,
		SourceComponent:     row.SourceComponent,
		SourceMineName:      row.SourceMineName,
		SourcePurlType:      row.SourcePurlType,
		SourceRepositoryURL: row.SourceRepositoryURL,
	}, nil
}
