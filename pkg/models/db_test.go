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

package models

import (
	"context"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/scanoss/go-grpc-helper/pkg/grpc/database"
	"github.com/scanoss/go-models/internal/testutils"
	zlog "github.com/scanoss/zap-logging-helper/pkg/logger"
)

func TestNewDB(t *testing.T) {
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
	dbWrapper := NewDB(ctx, s, conn, q)

	if dbWrapper.AllUrls == nil {
		t.Error("NewDB did not initialize AllUrls model")
	}

	if dbWrapper.Projects == nil {
		t.Error("NewDB did not initialize Projects model")
	}

	if dbWrapper.Versions == nil {
		t.Error("NewDB did not initialize Versions model")
	}

	if dbWrapper.Licenses == nil {
		t.Error("NewDB did not initialize Licenses model")
	}

	if dbWrapper.Mines == nil {
		t.Error("NewDB did not initialize Mines model")
	}
}
