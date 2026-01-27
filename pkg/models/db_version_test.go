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
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/scanoss/go-models/internal/testutils"
	zlog "github.com/scanoss/zap-logging-helper/pkg/logger"
)

func TestDBVersionGetCurrentVersion(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	ctx := ctxzap.ToContext(context.Background(), zlog.L)
	db := testutils.SqliteSetup(t)
	defer testutils.CloseDB(t, db)
	testutils.LoadSQLDataFile(t, db, "../../internal/testutils/mock/db_version.sql")

	model := NewDBVersionModel(db)

	fmt.Println("Testing GetCurrentVersion with data...")
	version, err := model.GetCurrentVersion(ctx)
	if err != nil {
		t.Errorf("DBVersionModel.GetCurrentVersion() error = %v", err)
	}
	if len(version.SchemaVersion) == 0 {
		t.Errorf("DBVersionModel.GetCurrentVersion() returned empty schema version")
	}
	if version.SchemaVersion != "1.0.0" {
		t.Errorf("DBVersionModel.GetCurrentVersion() schema_version = %v, want 1.0.0", version.SchemaVersion)
	}
	if version.PackageName != "base" {
		t.Errorf("DBVersionModel.GetCurrentVersion() package_name = %v, want components", version.PackageName)
	}
	if version.DBRelease != "2026.01" {
		t.Errorf("DBVersionModel.GetCurrentVersion() db_release = %v, want 2026.01", version.DBRelease)
	}
	if version.CreatedAt != "2026-01-15" {
		t.Errorf("DBVersionModel.GetCurrentVersion() created_at = %v, want 2026-01-15", version.CreatedAt)
	}
	fmt.Printf("DBVersion: %#v\n", version)
}

// TestDBVersionGetCurrentVersionNoTable tests querying when the db_version table doesn't exist.
func TestDBVersionGetCurrentVersionNoTable(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	ctx := ctxzap.ToContext(context.Background(), zlog.L)
	db := testutils.SqliteSetup(t)
	defer testutils.CloseDB(t, db)

	model := NewDBVersionModel(db)

	fmt.Println("Testing GetCurrentVersion without table...")
	version, err := model.GetCurrentVersion(ctx)
	if !errors.Is(err, ErrTableNotFound) {
		t.Errorf("DBVersionModel.GetCurrentVersion() expected ErrTableNotFound, got %v", err)
	}
	if len(version.SchemaVersion) > 0 {
		t.Errorf("DBVersionModel.GetCurrentVersion() expected empty version for missing table, got %v", version.SchemaVersion)
	}
	fmt.Printf("DBVersion (empty expected): %#v\n", version)
}

// TestDBVersionGetCurrentVersionEmptyTable tests querying when the db_version table exists but is empty.
func TestDBVersionGetCurrentVersionEmptyTable(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	ctx := ctxzap.ToContext(context.Background(), zlog.L)
	db := testutils.SqliteSetup(t)
	defer testutils.CloseDB(t, db)
	// Create the table but don't insert any data
	testutils.LoadSQLDataFile(t, db, "../../internal/testutils/mock/db_version.sql")
	_, delErr := db.Exec("DELETE FROM db_version")
	if delErr != nil {
		t.Fatalf("failed to clear db_version table: %v", delErr)
	}

	model := NewDBVersionModel(db)

	fmt.Println("Testing GetCurrentVersion with empty table...")
	version, err := model.GetCurrentVersion(ctx)
	if err != nil {
		t.Errorf("DBVersionModel.GetCurrentVersion() error = %v, expected nil for empty table", err)
	}
	if len(version.SchemaVersion) > 0 {
		t.Errorf("DBVersionModel.GetCurrentVersion() expected empty version for empty table, got %v", version.SchemaVersion)
	}
	fmt.Printf("DBVersion (empty expected): %#v\n", version)
}
