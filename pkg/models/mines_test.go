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
	"fmt"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/scanoss/go-models/internal/testutils"
	zlog "github.com/scanoss/zap-logging-helper/pkg/logger"
	_ "modernc.org/sqlite"
)

func TestMines(t *testing.T) {
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

	mine := NewMineModel(ctx, s, conn)

	var purlType = "maven"
	mineIds, err := mine.GetMineIdsByPurlType(purlType)
	if err != nil {
		t.Errorf("mines.GetMineIdByPurlType() error = %v", err)
	}
	fmt.Printf("Mine ID for %v: %v\n", purlType, mineIds)

	purlType = "gem"
	mineIds, err = mine.GetMineIdsByPurlType(purlType)
	if err != nil {
		t.Errorf("mines.GetMineIdByPurlType() error = %v", err)
	}
	fmt.Printf("Mine ID for %v: %v\n", purlType, mineIds)

	purlType = ""
	mineIds, err = mine.GetMineIdsByPurlType(purlType)
	if err != nil {
		fmt.Printf("Mine ID not found: %v\n", err)
	} else {
		t.Errorf("mines.GetMineIdByPurlType() found for %v = %v", purlType, mineIds)
	}

	purlType = "NONEXISTENT"
	mineIds, err = mine.GetMineIdsByPurlType(purlType)
	if err != nil {
		fmt.Printf("Mine ID not found: %v\n", err)
	} else {
		t.Errorf("mines.GetMineIdByPurlType() found for %v = %v", purlType, mineIds)
	}

	purlType = "npm"
	mineIds, err = mine.GetMineIdsByPurlType(purlType)
	if err != nil {
		t.Errorf("mines.GetMineIdsByPurlType() error = %v", err)
	}
	fmt.Printf("Mine IDs for %v: %v\n", purlType, mineIds)
}

// TestMinesBadSql test queries without creating/loading the mines table.
func TestMinesBadSql(t *testing.T) {
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

	mine := NewMineModel(ctx, s, conn)
	purlType := "NONEXISTENT"
	mineIds, err := mine.GetMineIdsByPurlType(purlType)
	if err != nil {
		fmt.Printf("Mine ID not found: %v\n", err)
	} else {
		t.Errorf("mines.GetMineIdByPurlType() found for %v = %v", purlType, mineIds)
	}
}
