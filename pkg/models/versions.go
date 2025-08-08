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

// Handle all interaction with the versions table

package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/jmoiron/sqlx"
)

type VersionModel struct {
	db *sqlx.DB
}

type Version struct {
	ID          int32  `db:"id"`
	VersionName string `db:"version_name"`
	SemVer      string `db:"semver"`
}

// NewVersionModel creates a new instance of the Version Model.
func NewVersionModel(db *sqlx.DB) *VersionModel {
	return &VersionModel{db: db}
}

// GetVersionByName gets the given version from the versions table.
func (m *VersionModel) GetVersionByName(ctx context.Context, name string) (Version, error) {
	s := ctxzap.Extract(ctx).Sugar()
	if len(name) == 0 {
		s.Error("Please specify a valid Version Name to query")
		return Version{}, errors.New("please specify a valid Version Name to query")
	}
	var version Version
	err := m.db.QueryRowxContext(ctx,
		"SELECT id, version_name, semver FROM versions"+
			" WHERE version_name = $1",
		name).StructScan(&version)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.Errorf("Error: Failed to query versions table for %v: %v", name, err)
		return Version{}, fmt.Errorf("failed to query the versions table: %v", err)
	}

	return version, nil
}
