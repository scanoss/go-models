// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2018-2023 SCANOSS.COM
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
	"github.com/jmoiron/sqlx"
)

// Models provides unified access to all SCANOSS data models.
// It maintains database connections and provides access to individual model instances.
type Models struct {
	AllUrls  *AllUrlsModel
	Projects *ProjectModel
	Versions *VersionModel
	Licenses *LicenseModel
	Mines    *MineModel
}

// NewModels creates a new instance of the unified SCANOSS models database wrapper.
// It initializes all individual models and sets up their dependencies.
func NewModels(db *sqlx.DB) *Models {
	models := &Models{
		AllUrls:  NewAllURLModel(db),
		Projects: NewProjectModel(db),
		Versions: NewVersionModel(db),
		Licenses: NewLicenseModel(db),
		Mines:    NewMineModel(db),
	}

	return models
}
