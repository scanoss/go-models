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
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/scanoss/go-grpc-helper/pkg/grpc/database"
	"go.uber.org/zap"
)

// DB provides unified access to all SCANOSS data models.
// It maintains database connections and provides access to individual model instances.
type DB struct {
	ctx  context.Context
	s    *zap.SugaredLogger
	conn *sqlx.Conn //TODO: refactor all models and replace with *database.DBQueryContext
	q    *database.DBQueryContext

	AllUrls        *AllUrlsModel
	Projects       *ProjectModel
	Versions       *VersionModel
	Licenses       *LicenseModel
	Mines          *MineModel
	GolangProjects *GolangProjects
}

// NewDB creates a new instance of the unified SCANOSS models database wrapper.
// It initializes all individual models and sets up their dependencies.
func NewDB(ctx context.Context, s *zap.SugaredLogger, conn *sqlx.Conn, q *database.DBQueryContext) *DB {
	db := &DB{
		ctx:  ctx,
		s:    s,
		conn: conn,
		q:    q,
	}

	// Initialize core models
	db.Projects = NewProjectModel(ctx, s, conn)
	db.Versions = NewVersionModel(ctx, s, conn)
	db.Licenses = NewLicenseModel(ctx, s, conn)
	db.Mines = NewMineModel(ctx, s, conn)

	db.GolangProjects = NewGolangProjectModel(ctx, s, conn, q)
	db.AllUrls = NewAllURLModel(ctx, s, q)

	return db
}

// Close closes the database connection and releases any resources.
// This should be called when the database is no longer needed.
func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}
