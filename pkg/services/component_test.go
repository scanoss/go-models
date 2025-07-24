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

	q := database.NewDBSelectContext(s, db, conn, false)
	modelsDB := models.NewDB(ctx, s, conn, q)

	service := NewComponentService(ctx, s, modelsDB)

	if service == nil {
		t.Fatal("NewComponentService returned nil")
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

	q := database.NewDBSelectContext(s, db, conn, false)
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

	q := database.NewDBSelectContext(s, db, conn, false)
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

	q := database.NewDBSelectContext(s, db, conn, false)
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

func TestPickOneUrl(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	ctx := ctxzap.ToContext(context.Background(), zlog.L)
	s := ctxzap.Extract(ctx).Sugar()
	db := testutils.SqliteSetup(t)
	defer testutils.CloseDB(t, db)
	conn := testutils.SqliteConn(t, ctx, db)
	defer testutils.CloseConn(t, conn)
	testutils.LoadMockSQLData(t, db, "../../internal/testutils/mock")

	q := database.NewDBSelectContext(s, db, conn, false)
	modelsDB := models.NewDB(ctx, s, conn, q)
	service := NewComponentService(ctx, s, modelsDB)

	tests := []struct {
		name          string
		urls          []models.AllURL
		component     string
		purlType      string
		requirement   string
		expectedVer   string
		shouldError   bool
		expectedEmpty bool
	}{
		{
			name:          "empty urls",
			urls:          []models.AllURL{},
			component:     "",
			purlType:      "",
			requirement:   "",
			expectedVer:   "",
			shouldError:   false,
			expectedEmpty: true,
		},
		{
			name: "multiple versions - picks highest",
			urls: []models.AllURL{
				{
					Component: "lodash",
					Version:   "1.0.0",
					SemVer:    "1.0.0",
					PurlName:  "lodash",
					MineID:    2,
				},
				{
					Component: "lodash",
					Version:   "2.0.0",
					SemVer:    "2.0.0",
					PurlName:  "lodash",
					MineID:    2,
				},
			},
			component:     "lodash",
			purlType:      "npm",
			requirement:   "",
			expectedVer:   "2.0.0",
			shouldError:   false,
			expectedEmpty: false,
		},
		{
			name: "version constraints filtering",
			urls: []models.AllURL{
				{
					Component: "lodash",
					Version:   "v1.0.0",
					SemVer:    "1.0.0",
					PurlName:  "lodash",
					MineID:    2,
				},
				{
					Component: "lodash",
					Version:   "v2.0.0",
					SemVer:    "2.0.0",
					PurlName:  "lodash",
					MineID:    2,
				},
			},
			component:     "lodash",
			purlType:      "npm",
			requirement:   "^v1.0.0",
			expectedVer:   "v1.0.0",
			shouldError:   false,
			expectedEmpty: false,
		},
		{
			name: "no versions after filter",
			urls: []models.AllURL{
				{
					Component: "lodash",
					Version:   "1.0.0",
					SemVer:    "1.0.0",
					PurlName:  "lodash",
					MineID:    2,
				},
			},
			component:     "lodash",
			purlType:      "npm",
			requirement:   "^2.0.0",
			expectedVer:   "",
			shouldError:   false,
			expectedEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.pickOneUrl(tt.urls, tt.component, tt.purlType, tt.requirement)

			if tt.shouldError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.expectedEmpty {
				if len(result.Version) != 0 {
					t.Errorf("expected empty result but got version: %s", result.Version)
				}
			} else {
				if result.Version != tt.expectedVer {
					t.Errorf("expected version %s, got %s", tt.expectedVer, result.Version)
				}
			}
		})
	}
}
