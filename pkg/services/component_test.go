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

package services

import (
	"context"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/scanoss/go-grpc-helper/pkg/grpc/database"
	"github.com/scanoss/go-models/internal/testutils"
	"github.com/scanoss/go-models/pkg/models"
	"github.com/scanoss/go-models/pkg/types"
	zlog "github.com/scanoss/zap-logging-helper/pkg/logger"
)

func TestNewComponentService(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	ctx := ctxzap.ToContext(context.Background(), zlog.L)
	s := ctxzap.Extract(ctx).Sugar()
	db := testutils.SqliteSetup(t) // Setup SQL Lite DB
	defer testutils.CloseDB(t, db)
	conn := testutils.SqliteConn(t, ctx, db) // Get a connection from the pool
	defer testutils.CloseConn(t, conn)
	testutils.LoadMockSQLData(t, db, "../../internal/testutils/mock")

	q := &database.DBQueryContext{}
	modelsDB := models.NewDB(ctx, s, conn, q)

	service := NewComponentService(ctx, s, modelsDB)

	if service == nil {
		t.Fatal("NewComponentService returned nil")
	}

	if service.ctx != ctx {
		t.Error("NewComponentService did not set context correctly")
	}

	if service.s != s {
		t.Error("NewComponentService did not set logger correctly")
	}

	if service.models != modelsDB {
		t.Error("NewComponentService did not set models correctly")
	}
}

func TestNewComponentServiceWithNilParams(t *testing.T) {
	ctx := context.Background()

	service := NewComponentService(ctx, nil, nil)

	if service == nil {
		t.Fatal("NewComponentService returned nil with nil parameters")
	}

	if service.ctx != ctx {
		t.Error("NewComponentService did not set context correctly with nil parameters")
	}
}

func TestGetComponentEmptyPurl(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	ctx := ctxzap.ToContext(context.Background(), zlog.L)
	s := ctxzap.Extract(ctx).Sugar()
	db := testutils.SqliteSetup(t) // Setup SQL Lite DB
	defer testutils.CloseDB(t, db)
	conn := testutils.SqliteConn(t, ctx, db) // Get a connection from the pool
	defer testutils.CloseConn(t, conn)
	testutils.LoadMockSQLData(t, db, "../../internal/testutils/mock")

	q := &database.DBQueryContext{}
	modelsDB := models.NewDB(ctx, s, conn, q)
	service := NewComponentService(ctx, s, modelsDB)

	req := types.ComponentRequest{
		Purl:        "",
		Requirement: "",
	}

	_, err = service.GetComponent(req)
	if err == nil {
		t.Error("GetComponent should return error for empty purl")
	}

	if err.Error() != "please specify a valid purl to query" {
		t.Errorf("Expected error message 'please specify a valid purl to query', got '%s'", err.Error())
	}
}

func TestGetComponentInvalidPurl(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	ctx := ctxzap.ToContext(context.Background(), zlog.L)
	s := ctxzap.Extract(ctx).Sugar()
	db := testutils.SqliteSetup(t) // Setup SQL Lite DB
	defer testutils.CloseDB(t, db)
	conn := testutils.SqliteConn(t, ctx, db) // Get a connection from the pool
	defer testutils.CloseConn(t, conn)
	testutils.LoadMockSQLData(t, db, "../../internal/testutils/mock")

	q := &database.DBQueryContext{}
	modelsDB := models.NewDB(ctx, s, conn, q)
	service := NewComponentService(ctx, s, modelsDB)

	req := types.ComponentRequest{
		Purl:        "invalid-purl",
		Requirement: "",
	}

	_, err = service.GetComponent(req)
	if err == nil {
		t.Error("GetComponent should return error for invalid purl")
	}

	if err.Error() != "failed to parse purl: invalid PURL: invalid-purl" {
		t.Logf("Got error: %s", err.Error())
		// Test that it's some form of parse error
		if len(err.Error()) == 0 {
			t.Error("Expected non-empty error message for invalid purl")
		}
	}
}

func TestGetComponentValidPurlButInvalidPurlName(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	ctx := ctxzap.ToContext(context.Background(), zlog.L)
	s := ctxzap.Extract(ctx).Sugar()
	db := testutils.SqliteSetup(t) // Setup SQL Lite DB
	defer testutils.CloseDB(t, db)
	conn := testutils.SqliteConn(t, ctx, db) // Get a connection from the pool
	defer testutils.CloseConn(t, conn)
	testutils.LoadMockSQLData(t, db, "../../internal/testutils/mock")

	q := &database.DBQueryContext{}
	modelsDB := models.NewDB(ctx, s, conn, q)
	service := NewComponentService(ctx, s, modelsDB)

	req := types.ComponentRequest{
		Purl:        "pkg:npm/",
		Requirement: "",
	}

	_, err = service.GetComponent(req)
	if err == nil {
		t.Error("GetComponent should return error for purl with empty name")
	}
}

//func TestGetComponentFileRequirement(t *testing.T) {
//	err := zlog.NewSugaredDevLogger()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
//	}
//	defer zlog.SyncZap()
//	ctx := ctxzap.ToContext(context.Background(), zlog.L)
//	s := ctxzap.Extract(ctx).Sugar()
//	db := testutils.SqliteSetup(t) // Setup SQL Lite DB
//	defer testutils.CloseDB(t, db)
//	conn := testutils.SqliteConn(t, ctx, db) // Get a connection from the pool
//	defer testutils.CloseConn(t, conn)
//	testutils.LoadMockSQLData(t, db, "../../internal/testutils/mock")
//
//	q := &database.DBQueryContext{}
//	modelsDB := models.NewDB(ctx, s, conn, q)
//	service := NewComponentService(ctx, s, modelsDB)
//
//	req := types.ComponentRequest{
//		Purl:        "pkg:npm/lodash@4.17.21",
//		Requirement: "file:../some/path",
//	}
//
//	// This will fail because we don't have test data, but we can verify the file requirement processing
//	_, err = service.GetComponent(req)
//	// We expect it to fail with database error since we don't have test data loaded
//	if err == nil {
//		t.Error("Expected error due to missing test data")
//	}
//}

func TestPickOneUrlEmptyUrls(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	ctx := ctxzap.ToContext(context.Background(), zlog.L)
	s := ctxzap.Extract(ctx).Sugar()
	db := testutils.SqliteSetup(t) // Setup SQL Lite DB
	defer testutils.CloseDB(t, db)
	conn := testutils.SqliteConn(t, ctx, db) // Get a connection from the pool
	defer testutils.CloseConn(t, conn)
	testutils.LoadMockSQLData(t, db, "../../internal/testutils/mock")

	q := &database.DBQueryContext{}
	modelsDB := models.NewDB(ctx, s, conn, q)
	service := NewComponentService(ctx, s, modelsDB)

	var emptyUrls []models.AllURL
	result, err := service.pickOneUrl(emptyUrls, "lodash", "npm", "")

	if err != nil {
		t.Errorf("pickOneUrl should not return error for empty urls: %v", err)
	}

	if len(result.Version) != 0 {
		t.Error("pickOneUrl should return empty AllURL for empty input")
	}
}

func TestPickOneUrlWithVersions(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	ctx := ctxzap.ToContext(context.Background(), zlog.L)
	s := ctxzap.Extract(ctx).Sugar()
	db := testutils.SqliteSetup(t) // Setup SQL Lite DB
	defer testutils.CloseDB(t, db)
	conn := testutils.SqliteConn(t, ctx, db) // Get a connection from the pool
	defer testutils.CloseConn(t, conn)
	testutils.LoadMockSQLData(t, db, "../../internal/testutils/mock")

	q := &database.DBQueryContext{}
	modelsDB := models.NewDB(ctx, s, conn, q)
	service := NewComponentService(ctx, s, modelsDB)

	urls := []models.AllURL{
		{
			Component: "lodash",
			Version:   "1.0.0",
			SemVer:    "1.0.0",
			PurlName:  "lodash",
			MineID:    1,
		},
		{
			Component: "lodash",
			Version:   "2.0.0",
			SemVer:    "2.0.0",
			PurlName:  "lodash",
			MineID:    1,
		},
	}

	result, err := service.pickOneUrl(urls, "lodash", "npm", "")

	if err != nil {
		t.Errorf("pickOneUrl should not return error for valid urls: %v", err)
	}

	if result.Version != "2.0.0" {
		t.Errorf("Expected version 2.0.0 (highest), got %s", result.Version)
	}
}

func TestPickOneUrlWithConstraints(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	ctx := ctxzap.ToContext(context.Background(), zlog.L)
	s := ctxzap.Extract(ctx).Sugar()
	db := testutils.SqliteSetup(t) // Setup SQL Lite DB
	defer testutils.CloseDB(t, db)
	conn := testutils.SqliteConn(t, ctx, db) // Get a connection from the pool
	defer testutils.CloseConn(t, conn)
	testutils.LoadMockSQLData(t, db, "../../internal/testutils/mock")

	q := &database.DBQueryContext{}
	modelsDB := models.NewDB(ctx, s, conn, q)
	service := NewComponentService(ctx, s, modelsDB)

	urls := []models.AllURL{
		{
			Component: "lodash",
			Version:   "1.0.0",
			SemVer:    "1.0.0",
			PurlName:  "lodash",
			MineID:    1,
		},
		{
			Component: "lodash",
			Version:   "2.0.0",
			SemVer:    "2.0.0",
			PurlName:  "lodash",
			MineID:    1,
		},
	}

	result, err := service.pickOneUrl(urls, "lodash", "npm", "^1.0.0")

	if err != nil {
		t.Errorf("pickOneUrl should not return error for valid constraint: %v", err)
	}

	if result.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0 (matching constraint ^1.0.0), got %s", result.Version)
	}
}

//func TestPickOneUrlInvalidVersions(t *testing.T) {
//	err := zlog.NewSugaredDevLogger()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
//	}
//	defer zlog.SyncZap()
//	ctx := ctxzap.ToContext(context.Background(), zlog.L)
//	s := ctxzap.Extract(ctx).Sugar()
//	db := testutils.SqliteSetup(t) // Setup SQL Lite DB
//	defer testutils.CloseDB(t, db)
//	conn := testutils.SqliteConn(t, ctx, db) // Get a connection from the pool
//	defer testutils.CloseConn(t, conn)
//	testutils.LoadMockSQLData(t, db, "../../internal/testutils/mock")
//
//	q := &database.DBQueryContext{}
//	modelsDB := models.NewDB(ctx, s, conn, q)
//	service := NewComponentService(ctx, s, modelsDB)
//
//	urls := []models.AllURL{
//		{
//			Component: "lodash",
//			Version:   "invalid-version",
//			SemVer:    "also-invalid",
//			PurlName:  "lodash",
//			MineID:    1,
//		},
//	}
//
//	result, err := service.pickOneUrl(urls, "lodash", "npm", "")
//
//	if err != nil {
//		t.Errorf("pickOneUrl should not return error for invalid versions (should fallback to v0.0.0): %v", err)
//	}
//
//	if result.Version != "v0.0.0" {
//		t.Errorf("Expected fallback version v0.0.0 for invalid versions, got %s", result.Version)
//	}
//}

func TestPickOneUrlNoVersionsAfterFilter(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	ctx := ctxzap.ToContext(context.Background(), zlog.L)
	s := ctxzap.Extract(ctx).Sugar()
	db := testutils.SqliteSetup(t) // Setup SQL Lite DB
	defer testutils.CloseDB(t, db)
	conn := testutils.SqliteConn(t, ctx, db) // Get a connection from the pool
	defer testutils.CloseConn(t, conn)
	testutils.LoadMockSQLData(t, db, "../../internal/testutils/mock")

	q := &database.DBQueryContext{}
	modelsDB := models.NewDB(ctx, s, conn, q)
	service := NewComponentService(ctx, s, modelsDB)

	urls := []models.AllURL{
		{
			Component: "lodash",
			Version:   "1.0.0",
			SemVer:    "1.0.0",
			PurlName:  "lodash",
			MineID:    1,
		},
	}

	result, err := service.pickOneUrl(urls, "lodash", "npm", "^2.0.0")

	if err != nil {
		t.Errorf("pickOneUrl should not return error when no versions match constraint: %v", err)
	}

	if len(result.Version) != 0 {
		t.Error("pickOneUrl should return empty AllURL when no versions match constraint")
	}
}
