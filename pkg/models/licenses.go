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
	"encoding/json"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
)

type LicenseModel struct {
	db *sqlx.DB
}

// SeeAlsoArray represents an array of strings that can be stored as JSON in the database
type SeeAlsoArray []string

// Scan implements the sql.Scanner interface for database deserialization
func (s *SeeAlsoArray) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}

	switch v := value.(type) {
	case string:
		if v == "" {
			*s = []string{}
			return nil
		}
		// Check if it's PostgreSQL array format
		if strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}") {
			return s.parsePostgreSQLArray(v)
		}
		// Otherwise try JSON
		return json.Unmarshal([]byte(v), s)
	case []byte:
		if len(v) == 0 {
			*s = []string{}
			return nil
		}
		str := string(v)
		// Check if it's PostgreSQL array format
		if strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}") {
			return s.parsePostgreSQLArray(str)
		}
		// Otherwise try JSON
		return json.Unmarshal(v, s)
	default:
		return fmt.Errorf("cannot scan %T into SeeAlsoArray", value)
	}
}

// parsePostgreSQLArray parses PostgreSQL array format like {val1,val2,val3}
func (s *SeeAlsoArray) parsePostgreSQLArray(str string) error {
	// Remove curly braces
	content := str[1 : len(str)-1]
	if content == "" {
		*s = []string{}
		return nil
	}
	
	// Split by comma and clean up each element
	parts := strings.Split(content, ",")
	result := make([]string, 0, len(parts))
	
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	
	*s = result
	return nil
}

// Value implements the driver.Valuer interface for database serialization
func (s SeeAlsoArray) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

type License struct {
	ID          int32  `db:"id"`
	LicenseName string `db:"license_name"`
	SPDX        string `db:"spdx_id"`
	IsSpdx      bool   `db:"is_spdx"`
}

type SPDXLicenseDetail struct {
	ID          string  `db:"id"`
	Type  string `db:"type"`
	Reference        string `db:"reference"`
	IsDeprecatedLicenseId      bool   `db:"isdeprecatedlicenseid"`
	DetailsURL         string `db:"detailsurl"`
	ReferenceNumber        string `db:"referencenumber"`
	Name        string `db:"name"`
	SeeAlso       SeeAlsoArray `db:"seealso"`
	IsOsiApproved        bool   `db:"isosiapproved"`
}

var bannedLicPrefixes = []string{"see ", "\"", "'", "-", "*", ".", "/", "?", "@", "\\", ";", ",", "`", "$"} // unwanted license prefixes
var bannedLicSuffixes = []string{".md", ".txt", ".html"}                                                    // unwanted license suffixes
var whiteSpaceRegex = regexp.MustCompile(`\s+`)                                                             // generic whitespace regex

// NewLicenseModel create a new instance of the License Model.
func NewLicenseModel(db *sqlx.DB) *LicenseModel {
	return &LicenseModel{db: db}
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

// GetSPDXLicenseDetails get spdx license details.
func (m* LicenseModel) GetSPDXLicenseDetails(ctx context.Context, spdxId string) (SPDXLicenseDetail, error) {
	s := ctxzap.Extract(ctx).Sugar()
	if spdxId == "" {
		s.Error("Please specify a valid SPDX ID to query")
		return SPDXLicenseDetail{}, errors.New("please specify a valid SPDX ID to query")
	}
	s.Debugf("Getting SPDX License Details for %v", spdxId)
	spdxIdToLower := strings.ToLower(spdxId)
	var license SPDXLicenseDetail
	err := m.db.QueryRowxContext(ctx,
		"SELECT * FROM spdx_license_data WHERE LOWER(id) = LOWER($1)", spdxIdToLower).StructScan(&license)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.Errorf("Error: Failed to query spdx_license_data table for %v: %#v", spdxId, err)
		return SPDXLicenseDetail{}, fmt.Errorf("failed to query the spdx_license_data table: %v", err)
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

