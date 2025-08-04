// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2018-2022 SCANOSS.COM
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

// Handle all interaction with the licenses table

package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/scanoss/go-grpc-helper/pkg/grpc/database"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
)

type LicenseModel struct {
	q  *database.DBQueryContext
	db *sqlx.DB
}

type License struct {
	ID          int32  `db:"id"`
	LicenseName string `db:"license_name"`
	LicenseID   string `db:"spdx_id"`
	IsSpdx      bool   `db:"is_spdx"`
}

var bannedLicPrefixes = []string{"see ", "\"", "'", "-", "*", ".", "/", "?", "@", "\\", ";", ",", "`", "$"} // unwanted license prefixes
var bannedLicSuffixes = []string{".md", ".txt", ".html"}                                                    // unwanted license suffixes
var whiteSpaceRegex = regexp.MustCompile(`\s+`)                                                             // generic whitespace regex

// NewLicenseModel create a new instance of the License Model.
func NewLicenseModel(q *database.DBQueryContext, db *sqlx.DB) *LicenseModel {
	return &LicenseModel{q: q, db: db}
}

// GetLicenseByID retrieves license data by the given row ID.
func (m *LicenseModel) GetLicenseByID(ctx context.Context, id int32) (License, error) {
	s := ctxzap.Extract(ctx).Sugar()
	if id < 0 {
		s.Error("Please specify a valid License ID to query")
		return License{}, errors.New("please specify a valid License Name to query")
	}
	var license License
	err := m.db.QueryRowxContext(ctx,
		"SELECT id, license_name, spdx_id, is_spdx FROM licenses"+
			" WHERE id = $1",
		id).StructScan(&license)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.Errorf("Error: Failed to query license table for %v: %#v", id, err)
		return License{}, fmt.Errorf("failed to query the license table: %v", err)
	}
	return license, nil
}

// GetLicenseByName retrieves the license details for the given license name.
func (m *LicenseModel) GetLicenseByName(ctx context.Context, name string) (License, error) {
	s := ctxzap.Extract(ctx).Sugar()
	if len(name) == 0 {
		s.Warn("No License Name specified to query")
		return License{}, nil
	}
	var license License
	err := m.db.QueryRowxContext(ctx,
		"SELECT id, license_name, spdx_id, is_spdx FROM licenses"+
			" WHERE license_name = $1",
		name,
	).StructScan(&license)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.Errorf("Failed to query license table for %v: %v", name, err)
		return License{}, fmt.Errorf("failed to query the license table: %v", err)
	}

	return license, nil
}

// CleanseLicenseName cleans up a license name to make it searchable in the licenses table.
func CleanseLicenseName(name string) (string, error) {
	if len(name) > 0 {
		name = strings.TrimSpace(name)     // remove leading/trailing spaces before even starting
		nameLower := strings.ToLower(name) // check banned strings against lowercase
		for _, prefix := range bannedLicPrefixes {
			if strings.HasPrefix(nameLower, prefix) {
				return "", fmt.Errorf("license name has banned prefix: %v", prefix)
			}
		}
		for _, suffix := range bannedLicSuffixes {
			if strings.HasSuffix(nameLower, suffix) {
				return "", fmt.Errorf("license name has banned suffix: %v", suffix)
			}
		}
		clean := whiteSpaceRegex.ReplaceAllString(name, " ")    // gets rid of new lines, tabs, etc.
		cleaner := whiteSpaceRegex.ReplaceAllString(clean, " ") // reduces it down to a single space
		cleanest := strings.ReplaceAll(cleaner, ",", ";")       // swap commas with semicolons
		// zlog.S.Debugf("in: %v clean: %v cleaner: %v cleanest: %v", name, clean, cleaner, cleanest)
		return strings.TrimSpace(cleanest), nil // return the cleansed license name
	}
	return "", nil // empty string, so just return it.
}
