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

package scanoss

import (
	"context"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/scanoss/go-grpc-helper/pkg/grpc/database"
	"github.com/scanoss/go-models/internal/testutils"
	zlog "github.com/scanoss/zap-logging-helper/pkg/logger"
)

func TestNew(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	ctx := ctxzap.ToContext(context.Background(), zlog.L)
	logger := ctxzap.Extract(ctx).Sugar()
	db := testutils.SqliteSetup(t) // Setup SQL Lite DB
	defer testutils.CloseDB(t, db)
	conn := testutils.SqliteConn(t, ctx, db) // Get a connection from the pool
	defer testutils.CloseConn(t, conn)
	testutils.LoadMockSQLData(t, db, "../../internal/testutils/mock")

	q := &database.DBQueryContext{}

	client := New(ctx, logger, conn, q)

	if client == nil {
		t.Fatal("New returned nil")
	}

	if client.ctx != ctx {
		t.Error("New did not set context correctly")
	}

	if client.logger != logger {
		t.Error("New did not set logger correctly")
	}

	if client.conn != conn {
		t.Error("New did not set connection correctly")
	}

	if client.Models == nil {
		t.Error("New did not initialize Models")
	}

	if client.Component == nil {
		t.Error("New did not initialize Component service")
	}
}

func TestNewWithNilParams(t *testing.T) {
	ctx := context.Background()

	client := New(ctx, nil, nil, nil)

	if client == nil {
		t.Fatal("New returned nil with nil parameters")
	}

	if client.ctx != ctx {
		t.Error("New did not set context correctly with nil parameters")
	}

	if client.Models == nil {
		t.Error("New did not initialize Models with nil parameters")
	}

	if client.Component == nil {
		t.Error("New did not initialize Component service with nil parameters")
	}
}

func TestClientIntegration(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	ctx := ctxzap.ToContext(context.Background(), zlog.L)
	logger := ctxzap.Extract(ctx).Sugar()
	db := testutils.SqliteSetup(t)
	defer testutils.CloseDB(t, db)
	conn := testutils.SqliteConn(t, ctx, db)
	defer testutils.CloseConn(t, conn)

	q := &database.DBQueryContext{}
	client := New(ctx, logger, conn, q)

	if client.Models.AllUrls == nil {
		t.Error("Client integration: AllUrls model not initialized")
	}

	if client.Models.Projects == nil {
		t.Error("Client integration: Projects model not initialized")
	}

	if client.Models.Versions == nil {
		t.Error("Client integration: Versions model not initialized")
	}

	if client.Models.Licenses == nil {
		t.Error("Client integration: Licenses model not initialized")
	}

	if client.Models.Mines == nil {
		t.Error("Client integration: Mines model not initialized")
	}

	if client.Component == nil {
		t.Error("Client integration: Component service not initialized")
	}
}
