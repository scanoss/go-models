// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2025 SCANOSS.COM
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

package testutils

import (
	"context"
	"testing"

	_ "modernc.org/sqlite"
)

func TestSqliteSetupAndCleanup(t *testing.T) {
	ctx := context.Background()

	db := SqliteSetup(t)
	conn := SqliteConn(t, ctx, db)

	CloseConn(t, conn)
	CloseDB(t, db)
}

func TestLoadValidSQLData(t *testing.T) {
	db := SqliteSetup(t)
	defer CloseDB(t, db)

	LoadSQLDataFile(t, db, "mock/mines.sql")
}

func TestLoadMockSQLData(t *testing.T) {
	db := SqliteSetup(t)
	defer CloseDB(t, db)

	LoadMockSQLData(t, db, "mock")
}
