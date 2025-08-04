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

	testutils.LoadMockSQLData(t, db, "../../internal/testutils/mock")

	q := database.NewDBSelectContext(s, db, nil, false)
	modelsDB := models.NewModels(q, db)

	service := NewComponentService(modelsDB)

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

	testutils.LoadMockSQLData(t, db, "../../internal/testutils/mock")

	q := database.NewDBSelectContext(s, db, nil, false)
	modelsDB := models.NewModels(q, db)
	service := NewComponentService(modelsDB)

	req := types.ComponentRequest{
		Purl:        "",
		Requirement: "",
	}

	_, err = service.GetComponent(ctx, req)
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

	testutils.LoadMockSQLData(t, db, "../../internal/testutils/mock")

	q := database.NewDBSelectContext(s, db, nil, false)
	modelsDB := models.NewModels(q, db)
	service := NewComponentService(modelsDB)

	req := types.ComponentRequest{
		Purl:        "invalid-purl",
		Requirement: "",
	}

	_, err = service.GetComponent(ctx, req)
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

	testutils.LoadMockSQLData(t, db, "../../internal/testutils/mock")

	q := database.NewDBSelectContext(s, db, nil, false)
	modelsDB := models.NewModels(q, db)
	service := NewComponentService(modelsDB)

	req := types.ComponentRequest{
		Purl:        "pkg:npm/",
		Requirement: "",
	}

	_, err = service.GetComponent(ctx, req)
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

	testutils.LoadMockSQLData(t, db, "../../internal/testutils/mock")

	q := database.NewDBSelectContext(s, db, nil, false)
	modelsDB := models.NewModels(q, db)
	service := NewComponentService(modelsDB)

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
			result, err := service.pickOneUrl(ctx, tt.urls, tt.component, tt.purlType, tt.requirement)

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

func TestGetComponent(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	ctx := ctxzap.ToContext(context.Background(), zlog.L)
	s := ctxzap.Extract(ctx).Sugar()
	db := testutils.SqliteSetup(t)
	defer testutils.CloseDB(t, db)

	testutils.LoadMockSQLData(t, db, "../../internal/testutils/mock")

	q := database.NewDBSelectContext(s, db, nil, false)
	modelsDB := models.NewModels(q, db)
	service := NewComponentService(modelsDB)

	tests := []struct {
		name        string
		purl        string
		requirement string
		expectedVer string
		shouldError bool
	}{
		{
			name:        "empty purl",
			purl:        "",
			requirement: "",
			shouldError: true,
			// Expected error because GetComponent requires a non-empty PURL string
		},
		{
			name:        "invalid purl format",
			purl:        "invalid-purl-format",
			requirement: "",
			shouldError: true,
			// Expected error because the string doesn't follow PURL format (pkg:type/name@version)
		},
		{
			name:        "valid purl with empty name",
			purl:        "pkg:npm/",
			requirement: "",
			shouldError: true,
			// Expected error because PURL has valid format but missing component name after "npm/"
		},
		{
			name:        "exact version match - electron-updater",
			purl:        "pkg:npm/electron-updater@4.0.8",
			requirement: "",
			expectedVer: "4.0.8",
			shouldError: false,
			// Expected version 4.0.8 because PURL specifies exact version and it exists in mock data
		},
		{
			name:        "exact version match - react",
			purl:        "pkg:npm/react@17.0.2",
			requirement: "",
			expectedVer: "17.0.2",
			shouldError: false,
			// Expected version 17.0.2 because PURL specifies exact version and it exists in mock data
		},
		{
			name:        "no version - picks latest electron-updater",
			purl:        "pkg:npm/electron-updater",
			requirement: "",
			expectedVer: "4.6.5",
			shouldError: false,
			// Expected version 4.6.5 because no version specified, so picks highest semver from 223 available versions.
			// grep 'electron-updater' internal/testutils/mock/all_urls.sql | awk '{print  $21 }' | sort
		},
		{
			name:        "no version - picks latest react",
			purl:        "pkg:npm/react",
			requirement: "",
			expectedVer: "18.0.0-beta-fdc1d617a-20211118",
			shouldError: false,
			// Expected: 18.0.0-beta-fdc1d617a-20211118 because it's the highest semver from 715 available versions
		},
		{
			name:        "version constraint - electron-updater ^4.0.0",
			purl:        "pkg:npm/electron-updater",
			requirement: "^4.0.0",
			expectedVer: "4.6.5",
			shouldError: false,
			// Expected: 4.6.5 because ^4.0.0 allows any 4.x version, and 4.6.5 is the highest 4.x available
		},
		{
			name:        "version constraint - uuid ~3.1.0",
			purl:        "pkg:npm/uuid",
			requirement: "~3.1.0",
			expectedVer: "3.1.0",
			shouldError: false,
			// Expected: 3.1.0 because ~3.1.0 allows patch versions (3.1.x), and 3.1.0 is the only match
		},
		{
			name:        "version constraint - react ^15.0.0",
			purl:        "pkg:npm/react",
			requirement: "^15.0.0",
			expectedVer: "15.7.0",
			shouldError: false,
			// Expected: 15.7.0 because ^15.0.0 allows 15.x versions, and 15.7.0 is highest 15.x (excludes 16.x+)
		},
		{
			name:        "non-existent component",
			purl:        "pkg:npm/non-existent-package",
			requirement: "",
			shouldError: true,
			// Expected error because component doesn't exist in mock database
		},
		{
			name:        "valid component but no matching versions",
			purl:        "pkg:npm/electron-updater",
			requirement: "^99.0.0",
			shouldError: true,
			// Expected error because no electron-updater versions match ^99.0.0 constraint
		},
		{
			name:        "exact version from requirement",
			purl:        "pkg:npm/electron-updater",
			requirement: "4.0.8",
			expectedVer: "4.0.8",
			shouldError: false,
			// Expected version 4.0.8 because requirement contains exact version (no PURL version), extracts 4.0.8
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := types.ComponentRequest{
				Purl:        tt.purl,
				Requirement: tt.requirement,
			}

			result, err := service.GetComponent(ctx, req)

			if tt.shouldError {
				if err == nil {
					t.Error("expected error but got none")
					return
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if result.Version != tt.expectedVer {
					t.Errorf("expected version %s, got %s", tt.expectedVer, result.Version)
				}
				if result.Purl != tt.purl {
					t.Errorf("expected purl %s, got %s", tt.purl, result.Purl)
				}
			}
		})
	}
}
