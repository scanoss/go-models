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

// Package testutils provides common test utilities for database setup and teardown
// across all SCANOSS Go models packages.
package testutils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// SqliteSetup sets up an in-memory SQLite DB for testing.
func SqliteSetup(t *testing.T) *sqlx.DB {
	db, err := sqlx.Connect("sqlite", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	if db == nil {
		t.Fatal("sqlx.Connect() returned nil database\n")
	}

	return db
}

// SqliteConn sets up a connection to a test DB.
func SqliteConn(t *testing.T, ctx context.Context, db *sqlx.DB) *sqlx.Conn {
	conn, err := db.Connx(ctx) // Get a connection from the pool
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	if conn == nil {
		t.Fatal("db.Connx() returned nil connection")
	}

	return conn
}

// CloseDB closes the specified DB and logs any errors.
func CloseDB(t *testing.T, db *sqlx.DB) {
	if db != nil {
		err := db.Close()
		if err != nil {
			t.Fatalf("Problem closing DB: %v", err)
		}
	}
}

// CloseConn closes the specified DB connection and logs any errors.
func CloseConn(t *testing.T, conn *sqlx.Conn) {
	if conn != nil {
		err := conn.Close()
		if err != nil {
			t.Fatalf("Problem closing DB connection: %v", err)
		}
	}
}

// LoadSQLDataFile loads the specified SQL file into the supplied DB.
func LoadSQLDataFile(t *testing.T, db *sqlx.DB, filename string) {
	fmt.Printf("Loading test data file: %v\n", filename)
	file, err := os.ReadFile(filename)

	if err != nil {
		t.Fatalf("LoadSQLDataFile() cannot read file %s - %v", filename, err)
	}

	if db == nil {
		t.Fatal("LoadSQLDataFile() - DB is null")
	}

	if _, err = db.Exec(string(file)); err != nil {
		t.Fatalf("LoadSQLDataFile() - Problem loading mock data into DB: %v", err)
	}
}

// LoadSQLDataFiles loads a list of test SQL files.
func LoadSQLDataFiles(t *testing.T, db *sqlx.DB, files []string) {
	for _, file := range files {
		LoadSQLDataFile(t, db, file)
	}
}

// LoadMockSQLData loads all the required test SQL files for the models package.
// This is a convenience function that loads the standard models test data files.
func LoadMockSQLData(t *testing.T, db *sqlx.DB, basePath string) {
	files := []string{
		filepath.Join(basePath, "mines.sql"),
		filepath.Join(basePath, "all_urls.sql"),
		filepath.Join(basePath, "projects.sql"),
		filepath.Join(basePath, "licenses.sql"),
		filepath.Join(basePath, "versions.sql"),
	}

	LoadSQLDataFiles(t, db, files)
}
