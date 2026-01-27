// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2026 SCANOSS.COM
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

// Handle all interaction with the db_version table

package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/jmoiron/sqlx"
)

// DBVersionModel provides database access for the db_version table.
type DBVersionModel struct {
	db *sqlx.DB
}

// DBVersion represents the database version information from the db_version table.
type DBVersion struct {
	PackageName   string `db:"package_name"`
	SchemaVersion string `db:"schema_version"`
	CreatedAt     string `db:"created_at"`
	DBRelease     string `db:"db_release"`
}

// NewDBVersionModel creates a new instance of the DBVersion Model.
func NewDBVersionModel(db *sqlx.DB) *DBVersionModel {
	return &DBVersionModel{db: db}
}

// GetCurrentVersion retrieves the current database schema version.
// Returns ErrTableNotFound if the db_version table does not exist.
// This check supports backward compatibility with databases that predate the db_version table.
func (m *DBVersionModel) GetCurrentVersion(ctx context.Context) (DBVersion, error) {
	s := ctxzap.Extract(ctx).Sugar()
	if !tableExists(ctx, m.db, "db_version") {
		s.Debug("db_version table does not exist")
		return DBVersion{}, ErrTableNotFound
	}
	var dbVersion DBVersion
	err := m.db.QueryRowxContext(ctx,
		"SELECT package_name, schema_version, created_at, db_release FROM db_version LIMIT 1").StructScan(&dbVersion)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.Debug("No version found in db_version table")
			return DBVersion{}, nil
		}
		return DBVersion{}, fmt.Errorf("failed to query db_version table: %w", err)
	}
	if t, err := time.Parse(time.RFC3339, dbVersion.CreatedAt); err == nil {
		dbVersion.CreatedAt = t.Format(time.DateOnly)
	}
	return dbVersion, nil
}
