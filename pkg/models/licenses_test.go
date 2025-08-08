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
	"reflect"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/scanoss/go-models/internal/testutils"
	zlog "github.com/scanoss/zap-logging-helper/pkg/logger"
)

func TestLicensesSearch(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	ctx := ctxzap.ToContext(context.Background(), zlog.L)
	db := testutils.SqliteSetup(t) // Setup SQL Lite DB
	defer testutils.CloseDB(t, db)
	testutils.LoadMockSQLData(t, db, "../../internal/testutils/mock")

	licenseModel := NewLicenseModel(db)

	var name = "MIT"
	fmt.Printf("Searching for license: %v\n", name)
	license, err := licenseModel.GetLicenseByName(ctx, name)
	if err != nil {
		t.Errorf("licenses.GetLicenseByName() error = %v", err)
	}
	if len(license.LicenseName) == 0 {
		t.Errorf("licenses.GetLicenseByName() No license returned from query")
	}
	fmt.Printf("License: %#v\n", license)

	name = ""
	fmt.Printf("Searching for license: %v\n", name)
	license, err = licenseModel.GetLicenseByName(ctx, name)
	if err != nil {
		t.Errorf("licenses.GetLicenseByName() error = %v", err)
	}
	if len(license.LicenseName) > 0 {
		t.Errorf("licenses.GetLicenseByName() License returned when one shouldn't")
	}
	fmt.Printf("License: %#v\n", license)

	name = "Unknown License"
	fmt.Printf("Searching for license: %v\n", name)
	license, err = licenseModel.GetLicenseByName(ctx, name)
	if err != nil {
		t.Errorf("licenses.GetLicenseByName() error = %v", err)
	}
	if len(license.LicenseName) != 0 {
		t.Errorf("licenses.GetLicenseByName() License returned when one shouldn't")
	}
}

func TestLicensesSearchId(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	ctx := ctxzap.ToContext(context.Background(), zlog.L)
	db := testutils.SqliteSetup(t) // Setup SQL Lite DB
	defer testutils.CloseDB(t, db)

	testutils.LoadMockSQLData(t, db, "../../internal/testutils/mock")

	licenseModel := NewLicenseModel(db)

	name := "MIT"
	fmt.Printf("Searching for license: %v\n", name)
	license, err := licenseModel.GetLicenseByName(ctx, name)
	if err != nil {
		t.Errorf("licenseModel.GetLicenseByName() error = %v", err)
	}
	if len(license.LicenseName) == 0 {
		t.Errorf("licenseModel.GetLicenseByName() No license returned from query")
	}
	fmt.Printf("Found License: %#v\n", license)

	id := license.ID
	fmt.Printf("Searching for license by id: %v\n", id)
	license, err = licenseModel.GetLicenseByID(ctx, id)
	if err != nil {
		t.Errorf("licenses.GetLicenseByID() error = %v", err)
	}
	if len(license.LicenseName) == 0 {
		t.Errorf("licenses.GetLicenseByID() No license returned from query")
	}
	fmt.Printf("License: %#v\n", license)

	id = 109
	fmt.Printf("Searching for license by id: %v\n", id)
	license, err = licenseModel.GetLicenseByID(ctx, id)
	if err != nil {
		t.Errorf("licenses.GetLicenseByID() error = %v", err)
	}
	if len(license.LicenseName) == 0 {
		t.Errorf("licenses.GetLicenseByID() No license returned from query")
	}
	fmt.Printf("License: %#v\n", license)

	id = -1
	fmt.Printf("Searching for license by id: %v\n", id)
	_, err = licenseModel.GetLicenseByID(ctx, id)
	if err == nil {
		t.Errorf("licenses.GetLicenseByID() error = did not get an error")
	} else {
		fmt.Printf("Got expected error = %v\n", err)
	}
}

func TestCleanseLicenseName(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "MIT",
			input: "MIT",
			want:  "MIT",
		},
		{
			name:  "Apache 2",
			input: " Apache 2.0 ",
			want:  "Apache 2.0",
		},
		{
			name:  "Apache/MIT",
			input: " Apache 2.0, 	MIT		",
			want:  "Apache 2.0; MIT",
		},
		{
			name:  "Empty String",
			input: "",
			want:  "",
		},
		{
			name:    "Banned prefixes",
			input:   "see something else",
			want:    "",
			wantErr: true,
		},
		{
			name:    "Banned suffixes",
			input:   "license name.html",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CleanseLicenseName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("licenses.CleanseLicenseName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("licenses.CleanseLicenseName() = %v, want %v", got, tt.want)
			}
		})
	}
}
