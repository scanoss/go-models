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

package models

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

// ErrTableNotFound is returned when a required database table does not exist.
var ErrTableNotFound = errors.New("table not found")

// tableExists checks if a table exists in an SQLite database.
// This is used for backward compatibility with databases that predate certain tables.
func tableExists(ctx context.Context, db *sqlx.DB, tableName string) bool {
	var exists bool
	err := db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM sqlite_master WHERE type='table' AND name=?)",
		tableName).Scan(&exists)
	return err == nil && exists
}
