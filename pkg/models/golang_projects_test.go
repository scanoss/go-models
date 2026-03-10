// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2026 SCANOSS.COM
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

	pkggodevclient "github.com/guseggert/pkggodev-client"
	"github.com/scanoss/go-models/internal/testutils"
	zlog "github.com/scanoss/zap-logging-helper/pkg/logger"
)

const (
	testGolangPurlName = "github.com/scanoss/papi"
	testGolangV001     = "v0.0.1"
	testGolangV002     = "v0.0.2"
	testGolangLicMIT   = "MIT"
)

func TestGolangProjectUrlsSearch(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	ctx := context.Background()
	db := testutils.SqliteSetup(t) // Setup SQL Lite DB
	defer testutils.CloseDB(t, db)
	testutils.LoadSQLDataFile(t, db, "../../internal/testutils/mock/golang_projects.sql")
	testutils.LoadSQLDataFile(t, db, "../../internal/testutils/mock/mines.sql")
	testutils.LoadSQLDataFile(t, db, "../../internal/testutils/mock/licenses.sql")
	testutils.LoadSQLDataFile(t, db, "../../internal/testutils/mock/versions.sql")
	golangProjModel := NewGolangProjectModel(db)

	url, err := golangProjModel.GetGolangUrlsByPurlNameType(ctx, "google.golang.org/grpc", "golang", "")
	if err != nil {
		t.Errorf("FAILED: golang_projects.GetUrlsByPurlName() error = %v", err)
	}
	if len(url.PurlName) == 0 {
		t.Errorf("FAILED: golang_projects.GetGoLangURLByPurlString() No URLs returned from query")
	}

	url, err = golangProjModel.GetGolangUrlsByPurlNameType(ctx, "NONEXISTENT", "none", "")
	if err != nil {
		t.Errorf("FAILED: golang_projects.GetGolangUrlsByPurlNameType() error = %v", err)
	}
	if len(url.PurlName) > 0 {
		t.Errorf("FAILED: golang_projects.GetGolangUrlsByPurlNameType() URLs found when none should be: %v", golangProjModel)
	}
	fmt.Printf("No Urls: %v\n", url)

	_, err = golangProjModel.GetGolangUrlsByPurlNameType(ctx, "NONEXISTENT", "", "")
	if err == nil {
		t.Errorf("FAILED: golang_projects.GetGolangUrlsByPurlNameType() error = did not get an error")
	} else {
		fmt.Printf("Got expected error = %v\n", err)
	}
	_, err = golangProjModel.GetGolangUrlsByPurlNameType(ctx, "", "", "")
	if err == nil {
		t.Errorf("FAILED: golang_projects.GetURLsByPurlString() error = did not get an error")
	} else {
		fmt.Printf("Got expected error = %v\n", err)
	}
	_, err = golangProjModel.GetGoLangURLByPurlString(ctx, "", "")
	if err == nil {
		t.Errorf("FAILED: golang_projects.GetURLsByPurlString() error = did not get an error")
	} else {
		fmt.Printf("Got expected error = %v\n", err)
	}
	_, err = golangProjModel.GetGoLangURLByPurlString(ctx, "rubbish-purl", "")
	if err == nil {
		t.Errorf("FAILED: golang_projects.GetGoLangURLByPurlString() error = did not get an error")
	} else {
		fmt.Printf("Got expected error = %v\n", err)
	}
	url, err = golangProjModel.GetGoLangURLByPurlString(ctx, "pkg:golang/google.golang.org/grpc", "")
	if err != nil {
		t.Errorf("FAILED: golang_projects.GetGoLangURLByPurlString() error = %v", err)
	}
	if len(url.PurlName) == 0 {
		t.Errorf("FAILED: golang_projects.GetGoLangURLByPurlString() No URLs returned from query")
	}
	fmt.Printf("Golang URL: %v\n", url)
}

func TestGolangProjectsSearchVersion(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	ctx := context.Background()
	db := testutils.SqliteSetup(t) // Setup SQL Lite DB
	testutils.LoadSQLDataFile(t, db, "../../internal/testutils/mock/golang_projects.sql")
	testutils.LoadSQLDataFile(t, db, "../../internal/testutils/mock/mines.sql")
	testutils.LoadSQLDataFile(t, db, "../../internal/testutils/mock/licenses.sql")
	testutils.LoadSQLDataFile(t, db, "../../internal/testutils/mock/versions.sql")
	golangProjModel := NewGolangProjectModel(db)

	url, err := golangProjModel.GetGolangUrlsByPurlNameTypeVersion(ctx, "google.golang.org/grpc", "golang", "1.19.0")
	if err != nil {
		t.Errorf("FAILED: golang_projects.GetGolangUrlsByPurlNameTypeVersion() error = %v", err)
	}
	if len(url.PurlName) == 0 {
		t.Errorf("FAILED: golang_projects.GetGolangUrlsByPurlNameTypeVersion() No URLs returned from query")
	}
	fmt.Printf("Golang URL Version: %#v\n", url)

	url, err = golangProjModel.GetGoLangURLByPurlString(ctx, "pkg:golang/google.golang.org/grpc@v1.19.0", "")
	if err != nil {
		t.Errorf("FAILED: golang_projects.GetGoLangURLByPurlString() error = failed to find purl by version string")
	}
	fmt.Printf("Golang URL Version: %#v\n", url)

	_, err = golangProjModel.GetGolangUrlsByPurlNameTypeVersion(ctx, "", "", "")
	if err == nil {
		t.Errorf("FAILED: golang_projects.GetGolangUrlsByPurlNameTypeVersion() error = did not get an error")
	} else {
		fmt.Printf("Got expected error = %v\n", err)
	}
	_, err = golangProjModel.GetGolangUrlsByPurlNameTypeVersion(ctx, "NONEXISTENT", "", "")
	if err == nil {
		t.Errorf("FAILED: golang_projects.GetGolangUrlsByPurlNameTypeVersion() error = did not get an error")
	} else {
		fmt.Printf("Got expected error = %v\n", err)
	}
	_, err = golangProjModel.GetGolangUrlsByPurlNameTypeVersion(ctx, "NONEXISTENT", "NONEXISTENT", "")
	if err == nil {
		t.Errorf("FAILED: golang_projects.GetGolangUrlsByPurlNameTypeVersion() error = did not get an error")
	} else {
		fmt.Printf("Got expected error = %v\n", err)
	}

	url, err = golangProjModel.GetGoLangURLByPurlString(ctx, "pkg:golang/google.golang.org/grpc", "22.22.22") // Shouldn't exist
	if err != nil {
		t.Errorf("FAILED: golang_projects.GetGoLangURLByPurlString() error = failed to find purl by version string")
	}
	url, err = golangProjModel.GetGoLangURLByPurlString(ctx, "pkg:golang/google.golang.org/grpc", "=v1.19.0")
	if err != nil {
		t.Errorf("FAILED: golang_projects.GetGoLangURLByPurlString() error = %v", err)
	}
	if len(url.PurlName) == 0 {
		t.Errorf("FAILED: golang_projects.GetGoLangURLByPurlString() No URLs returned from query")
	}
	fmt.Printf("Golang URL: %v\n", url)
	url, err = golangProjModel.GetGoLangURLByPurlString(ctx, "pkg:golang/google.golang.org/grpc", "==v1.19.0")
	if err != nil {
		t.Errorf("FAILED: golang_projects.GetGoLangURLByPurlString() error = %v", err)
	}
	if len(url.PurlName) == 0 {
		t.Errorf("FAILED: golang_projects.GetGoLangURLByPurlString() No URLs returned from query")
	}
	fmt.Printf("Golang URL: %v\n", url)

	url, err = golangProjModel.GetGoLangURLByPurlString(ctx, "pkg:golang/google.golang.org/grpc@1.7.0", "") // Should be missing license
	if err != nil {
		t.Errorf("FAILED: golang_projects.GetGoLangURLByPurlString() error = %v", err)
	}
	if len(url.License) == 0 {
		t.Errorf("FAILED: golang_projects.GetGoLangURLByPurlString() No URL License returned from query")
	}
	fmt.Printf("Golang URL: %v\n", url)
}

func TestGolangProjectsSearchVersionRequirement(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	ctx := context.Background()
	db := testutils.SqliteSetup(t) // Setup SQL Lite DB
	testutils.LoadSQLDataFile(t, db, "../../internal/testutils/mock/golang_projects.sql")
	testutils.LoadSQLDataFile(t, db, "../../internal/testutils/mock/mines.sql")
	testutils.LoadSQLDataFile(t, db, "../../internal/testutils/mock/licenses.sql")
	testutils.LoadSQLDataFile(t, db, "../../internal/testutils/mock/versions.sql")
	golangProjModel := NewGolangProjectModel(db)

	url, err := golangProjModel.GetGoLangURLByPurlString(ctx, "pkg:golang/google.golang.org/grpc", ">0.0.4")
	if err != nil {
		t.Errorf("FAILED: golang_projects.GetUrlsByPurlName() error = %v", err)
	}
	if len(url.PurlName) == 0 {
		t.Errorf("FAILED: golang_projects.GetUrlsByPurlName() No URLs returned from query")
	}
	fmt.Printf("Golang URL Version: %#v\n", url)

	url, err = golangProjModel.GetGoLangURLByPurlString(ctx, "pkg:golang/google.golang.org/grpc", "v0.0.0-201910101010-s3333")
	if err != nil {
		t.Errorf("FAILED: golang_projects.GetUrlsByPurlName() error = %v", err)
	}
	if len(url.PurlName) == 0 {
		t.Errorf("FAILED: golang_projects.GetUrlsByPurlName() No URLs returned from query")
	}
	fmt.Printf("Golang URL Version: %#v\n", url)
}

func TestGolangPkgGoDev(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	ctx := context.Background()
	db := testutils.SqliteSetup(t) // Setup SQL Lite DB
	testutils.LoadSQLDataFile(t, db, "../../internal/testutils/mock/golang_projects.sql")
	testutils.LoadSQLDataFile(t, db, "../../internal/testutils/mock/mines.sql")
	testutils.LoadSQLDataFile(t, db, "../../internal/testutils/mock/licenses.sql")
	testutils.LoadSQLDataFile(t, db, "../../internal/testutils/mock/versions.sql")
	golangProjModel := NewGolangProjectModel(db)

	_, err = golangProjModel.queryPkgGoDev(ctx, "", "")
	if err == nil {
		t.Errorf("FAILED: golang_projects.queryPkgGoDev() error = did not get an error")
	}

	url, err := golangProjModel.getLatestPkgGoDev(ctx, "google.golang.org/grpc", "golang", "v0.0.0-201910101010-s3333")
	if err != nil {
		t.Errorf("FAILED: golang_projects.getLatestPkgGoDev() error = %v", err)
	}
	if len(url.PurlName) == 0 {
		t.Errorf("FAILED: golang_projects.getLatestPkgGoDev() No URLs returned from query")
	}
	fmt.Printf("Golang URL Version: %#v\n", url)

	url, err = golangProjModel.getLatestPkgGoDev(ctx, "github.com/scanoss/papi", "golang", "v0.0.3")
	if err != nil {
		t.Errorf("FAILED: golang_projects.getLatestPkgGoDev() error = %v", err)
	}
	if len(url.PurlName) == 0 {
		t.Errorf("FAILED: golang_projects.getLatestPkgGoDev() No URLs returned from query")
	}
	fmt.Printf("Golang URL Version: %#v\n", url)

	var allURL AllURL
	var license License
	var version Version
	fmt.Printf("SavePkg: %#v - %#v - %#v", allURL, license, version)
	err = golangProjModel.savePkg(ctx, allURL, version, license, nil)
	if err == nil {
		t.Errorf("FAILED: golangProjModel.savePkg() error = did not get an error")
	}
	allURL.PurlName = testGolangPurlName
	err = golangProjModel.savePkg(ctx, allURL, version, license, nil)
	if err == nil {
		t.Errorf("FAILED: golangProjModel.savePkg() error = did not get an error")
	}
	allURL.MineID = 45
	err = golangProjModel.savePkg(ctx, allURL, version, license, nil)
	if err == nil {
		t.Errorf("FAILED: golangProjModel.savePkg() error = did not get an error")
	}
	allURL.Version = testGolangV001
	version.VersionName = testGolangV001
	version.ID = 5958021
	err = golangProjModel.savePkg(ctx, allURL, version, license, nil)
	if err == nil {
		t.Errorf("FAILED: golangProjModel.savePkg() error = did not get an error")
	}
	license.LicenseName = testGolangLicMIT
	license.ID = 5614
	err = golangProjModel.savePkg(ctx, allURL, version, license, nil)
	if err == nil {
		t.Errorf("FAILED: golangProjModel.savePkg() error = did not get an error")
	}
	var comp pkggodevclient.Package
	comp.Package = testGolangPurlName
	comp.IsPackage = true
	comp.IsModule = true
	comp.Version = testGolangV001
	comp.License = testGolangLicMIT
	comp.HasRedistributableLicense = true
	comp.HasStableVersion = true
	comp.HasTaggedVersion = true
	comp.HasValidGoModFile = true
	comp.Repository = testGolangPurlName
	err = golangProjModel.savePkg(ctx, allURL, version, license, &comp)
	if err != nil {
		t.Errorf("FAILED: golangProjModel.savePkg() error = %v", err)
	}
	allURL.Version = testGolangV002
	version.VersionName = testGolangV002
	comp.Version = testGolangV002
	err = golangProjModel.savePkg(ctx, allURL, version, license, &comp)
	if err != nil {
		t.Errorf("FAILED: golangProjModel.savePkg() error = %v", err)
	}
}

func TestGolangProjectsSearchBadSql(t *testing.T) {
	err := zlog.NewSugaredDevLogger()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a sugared logger", err)
	}
	defer zlog.SyncZap()
	ctx := context.Background()
	db := testutils.SqliteSetup(t) // Setup SQL Lite DB
	_, _ = db.Exec("DROP TABLE IF EXISTS golang_projects")
	_, _ = db.Exec("DROP TABLE IF EXISTS mines")
	_, _ = db.Exec("DROP TABLE IF EXISTS licenses")
	_, _ = db.Exec("DROP TABLE IF EXISTS versions")
	golangProjModel := NewGolangProjectModel(db)

	_, err = golangProjModel.GetGoLangURLByPurlString(ctx, "pkg:golang/google.golang.org/grpc", "")
	if err == nil {
		t.Errorf("FAILED: golang_projects.GetGoLangURLByPurlString() error = did not get an error")
	} else {
		fmt.Printf("Got expected error = %v\n", err)
	}
	_, err = golangProjModel.GetGoLangURLByPurlString(ctx, "pkg:golang/google.golang.org/grpc@1.19.0", "")
	if err == nil {
		t.Errorf("FAILED: golang_projects.GetGoLangURLByPurlString() error = did not get an error")
	} else {
		fmt.Printf("Got expected error = %v\n", err)
	}
	_, err = golangProjModel.getLatestPkgGoDev(ctx, "github.com/scanoss/does-not-exist", "golang", "v0.0.99")
	if err == nil {
		t.Errorf("FAILED: golang_projects.getLatestPkgGoDev() error = did not get an error: %v", err)
	} else {
		fmt.Printf("Got expected error = %v\n", err)
	}
}
