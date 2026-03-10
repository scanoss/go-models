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
	"database/sql"
	"errors"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	pkggodevclient "github.com/guseggert/pkggodev-client"
	"github.com/jmoiron/sqlx"
	"github.com/package-url/packageurl-go"
	purlutils "github.com/scanoss/go-purl-helper/pkg"
)

type GolangProjects struct {
	db      *sqlx.DB
	ver     *VersionModel
	lic     *LicenseModel
	mine    *MineModel
	project *ProjectModel // TODO Do we add golang component to the projects table?
}

// NewGolangProjectModel creates a new instance of Golang Project Model.
func NewGolangProjectModel(db *sqlx.DB) *GolangProjects {
	return &GolangProjects{db: db,
		ver: NewVersionModel(db), lic: NewLicenseModel(db), mine: NewMineModel(db),
		project: NewProjectModel(db),
	}
}

// GetGoLangURLByPurlString searches the Golang Projects for the specified Purl (and requirement).
func (m *GolangProjects) GetGoLangURLByPurlString(ctx context.Context, purlString, purlReq string) (AllURL, error) {
	s := ctxzap.Extract(ctx).Sugar()
	if len(purlString) == 0 {
		s.Error("Please specify a valid Purl String to query")
		return AllURL{}, errors.New("please specify a valid Purl String to query")
	}
	purl, err := purlutils.PurlFromString(purlString)
	if err != nil {
		return AllURL{}, err
	}
	purlName, err := purlutils.PurlNameFromString(purlString)
	if err != nil {
		return AllURL{}, err
	}
	if len(purl.Version) == 0 && len(purlReq) > 0 { // No version specified, but we might have a specific version in the Requirement
		ver := purlutils.GetVersionFromReq(purlReq)
		if len(ver) > 0 {
			purl.Version = ver
			purlReq = ""
		}
	}
	return m.GetGoLangURLByPurl(ctx, purl, purlName, purlReq)
}

// GetGoLangURLByPurl searches the Golang Projects for the specified Purl Package (and optional requirement).
func (m *GolangProjects) GetGoLangURLByPurl(ctx context.Context, purl packageurl.PackageURL, purlName, purlReq string) (AllURL, error) {
	if len(purl.Version) > 0 {
		return m.GetGolangUrlsByPurlNameTypeVersion(ctx, purlName, purl.Type, purl.Version)
	}
	return m.GetGolangUrlsByPurlNameType(ctx, purlName, purl.Type, purlReq)
}

// GetGolangUrlsByPurlNameType searches Golang Project for the specified Purl by Purl Type (and optional requirement).
func (m *GolangProjects) GetGolangUrlsByPurlNameType(ctx context.Context, purlName, purlType, purlReq string) (AllURL, error) {
	s := ctxzap.Extract(ctx).Sugar()
	if len(purlName) == 0 {
		s.Error("Please specify a valid Purl Name to query")
		return AllURL{}, errors.New("please specify a valid Purl Name to query")
	}
	if len(purlType) == 0 {
		s.Errorf("Please specify a valid Purl Type to query: %v", purlName)
		return AllURL{}, errors.New("please specify a valid Purl Type to query")
	}
	query := "SELECT component, v.version_name AS version, v.semver AS semver," +
		" l.license_name AS license, u.license_id, l.is_spdx AS is_spdx," +
		" purl_name, mine_id FROM golang_projects u" +
		" LEFT JOIN mines m ON u.mine_id = m.id" +
		" LEFT JOIN licenses l ON u.license_id = l.id" +
		" LEFT JOIN versions v ON u.version_id = v.id" +
		" WHERE m.purl_type = $1 AND u.purl_name = $2 AND is_indexed = True" +
		" ORDER BY version_date DESC"
	var allURLs []AllURL
	err := m.db.SelectContext(ctx, &allURLs, query, purlType, purlName)
	if err != nil {
		s.Errorf("Failed to query golang projects table for %v - %v: %v", purlType, purlName, err)
		return AllURL{}, fmt.Errorf("failed to query the golang projects table: %v", err)
	}
	s.Debugf("Found %v results for %v, %v.", len(allURLs), purlType, purlName)
	if len(allURLs) == 0 { // Check pkg.go.dev for the latest data
		s.Debugf("Checking PkgGoDev for live info...")
		allURL, pkgErr := m.getLatestPkgGoDev(ctx, purlName, purlType, "")
		if pkgErr == nil {
			s.Debugf("Retrieved golang data from pkg.go.dev: %#v", allURL)
			allURLs = append(allURLs, allURL)
		} else {
			s.Infof("Ran into an issue looking up pkg.go.dev for: %v. Ignoring", purlName)
		}
	}

	// Pick the most appropriate version to return
	return PickOneUrl(ctx, allURLs, purlName, purlType, purlReq)
}

// GetGolangUrlsByPurlNameTypeVersion searches Golang Projects for specified Purl, Type and Version.
func (m *GolangProjects) GetGolangUrlsByPurlNameTypeVersion(ctx context.Context, purlName, purlType, purlVersion string) (AllURL, error) {
	s := ctxzap.Extract(ctx).Sugar()
	if len(purlName) == 0 {
		s.Error("Please specify a valid Purl Name to query")
		return AllURL{}, errors.New("please specify a valid Purl Name to query")
	}
	if len(purlType) == 0 {
		s.Error("Please specify a valid Purl Type to query")
		return AllURL{}, errors.New("please specify a valid Purl Type to query")
	}
	if len(purlVersion) == 0 {
		s.Error("Please specify a valid Purl Version to query")
		return AllURL{}, errors.New("please specify a valid Purl Version to query")
	}
	query := "SELECT component, v.version_name AS version, v.semver AS semver," +
		" l.license_name AS license, u.license_id, l.is_spdx AS is_spdx," +
		" purl_name, mine_id FROM golang_projects u" +
		" LEFT JOIN mines m ON u.mine_id = m.id" +
		" LEFT JOIN licenses l ON u.license_id = l.id" +
		" LEFT JOIN versions v ON u.version_id = v.id" +
		" WHERE m.purl_type = $1 AND u.purl_name = $2 AND v.version_name = $3 AND is_indexed = True" +
		" ORDER BY version_date DESC"
	var allURLs []AllURL
	err := m.db.SelectContext(ctx, &allURLs, query, purlType, purlName, purlVersion)
	if err != nil {
		s.Errorf("Failed to query golang projects table for %v - %v: %v", purlType, purlName, err)
		return AllURL{}, fmt.Errorf("failed to query the golang projects table: %v", err)
	}
	s.Debugf("Found %v results for %v, %v.", len(allURLs), purlType, purlName)
	if len(allURLs) > 0 { // We found an entry. Let's check if it has license data
		allURL, errPickURL := PickOneUrl(ctx, allURLs, purlName, purlType, "")
		if errPickURL != nil {
			return AllURL{}, errPickURL
		}
		if len(allURL.License) == 0 { // No license data found. Need to search for live info
			s.Debugf("Couldn't find license data for component. Need to search live data")
			allURLs = allURLs[:0]
		} else {
			return allURL, nil // Return the component details
		}
	}
	if len(allURLs) == 0 { // Check pkg.go.dev for the latest data
		s.Debugf("Checking PkgGoDev for live info...")
		allURL, pkgErr := m.getLatestPkgGoDev(ctx, purlName, purlType, purlVersion)
		if pkgErr == nil {
			s.Debugf("Retrieved golang data from pkg.go.dev: %#v", allURL)
			allURLs = append(allURLs, allURL)
		} else {
			s.Infof("Ran into an issue looking up pkg.go.dev for: %v - %v. Ignoring", purlName, purlVersion)
		}
	}
	// Pick the most appropriate version to return
	return PickOneUrl(ctx, allURLs, purlName, purlType, "")
}

// savePkg writes the given package details to the Golang Projects table.
//
//goland:noinspection ALL
func (m *GolangProjects) savePkg(ctx context.Context, allURL AllURL, version Version, license License, comp *pkggodevclient.Package) error {
	s := ctxzap.Extract(ctx).Sugar()
	if len(allURL.PurlName) == 0 {
		s.Error("Please specify a valid Purl to save")
		return errors.New("please specify a valid Purl to save")
	}
	if allURL.MineID <= 0 {
		s.Error("Please specify a valid mine id to save")
		return errors.New("please specify a valid mine id to save")
	}
	if version.ID <= 0 || len(version.VersionName) == 0 {
		s.Error("Please specify a valid version to save")
		return errors.New("please specify a valid version to save")
	}
	if license.ID <= 0 || len(license.LicenseName) == 0 {
		s.Error("Please specify a valid license to save")
		return errors.New("please specify a valid license to save")
	}
	if comp == nil {
		s.Error("Please specify a valid component package to save")
		return errors.New("please specify a valid component package to save")
	}
	s.Debugf("Attempting to save '%#v' - %#v to the golang_projects table...", allURL, version)
	// Search for an existing entry first
	var existingPurl string
	err := m.db.QueryRowxContext(ctx,
		"SELECT purl_name FROM golang_projects"+
			" WHERE purl_name = $1 AND version = $2",
		allURL.PurlName, allURL.Version,
	).Scan(&existingPurl)
	if err != nil && err != sql.ErrNoRows {
		s.Warnf("Error: Problem encountered searching golang_projects table for %v: %v", allURL, err)
	}
	var purlName string
	sqlQueryType := "insert"
	if len(existingPurl) > 0 {
		// update entry
		sqlQueryType = "update"
		s.Debugf("Updating new Golang project: %#v", comp)
		//goland:noinspection ALL
		err = m.db.QueryRowxContext(ctx,
			"UPDATE golang_projects SET component = $1, version = $2, version_id = $3, version_date = $4,"+
				" is_module = $5, is_package = $6, license = $7, license_id = $8, has_valid_go_mod_file = $9,"+
				" has_redistributable_license = $10, has_tagged_version = $11, has_stable_version = $12,"+
				" repository = $13, is_indexed = $14, purl_name = $15, mine_id = $16"+
				" WHERE purl_name = $17 AND version = $18"+
				" RETURNING purl_name",
			allURL.Component, allURL.Version, version.ID, comp.Published,
			comp.IsModule, comp.IsPackage, license.LicenseName, license.ID, comp.HasValidGoModFile,
			comp.HasRedistributableLicense, comp.HasTaggedVersion, comp.HasStableVersion,
			comp.Repository, true, allURL.PurlName, allURL.MineID,
			allURL.PurlName, allURL.Version,
		).Scan(&purlName)
	} else {
		s.Debugf("Inserting new Golang project: %#v", comp)
		// insert new entry
		err = m.db.QueryRowxContext(ctx,
			"INSERT INTO golang_projects (component, version, version_id, version_date, is_module, is_package,"+
				" license, license_id, has_valid_go_mod_file, has_redistributable_license, has_tagged_version,"+
				" has_stable_version, repository, is_indexed, purl_name, mine_id, index_timestamp)"+
				" VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)"+
				" RETURNING purl_name",
			allURL.Component, allURL.Version, version.ID, comp.Published,
			comp.IsModule, comp.IsPackage, license.LicenseName, license.ID, comp.HasValidGoModFile,
			comp.HasRedistributableLicense, comp.HasTaggedVersion, comp.HasStableVersion,
			comp.Repository, true, allURL.PurlName, allURL.MineID, "",
		).Scan(&purlName)
	}
	if err != nil {
		s.Errorf("Error: Failed to %v new component into golang_projects table for %v - %#v: %v", sqlQueryType, allURL, comp, err)
		return fmt.Errorf("failed to %v new component into golang projects: %v", sqlQueryType, err)
	}
	s.Debugf("Completed %v of %v", sqlQueryType, purlName)
	return nil
}

// getLatestPkgGoDev retrieves the latest information about a Golang Package from https://pkg.go.dev
// If requested (via config), it will commit that data to the Golang Projects table.
func (m *GolangProjects) getLatestPkgGoDev(ctx context.Context, purlName, purlType, purlVersion string) (AllURL, error) {
	s := ctxzap.Extract(ctx).Sugar()
	allURL, err := m.queryPkgGoDev(ctx, purlName, purlVersion)
	if err != nil {
		return allURL, err
	}
	cleansedLicense, err := CleanseLicenseName(allURL.License)
	if err != nil {
		return allURL, err
	}
	license, _ := m.lic.GetLicenseByName(ctx, cleansedLicense)
	if len(license.LicenseName) == 0 {
		s.Warnf("No license details in DB for: %v", cleansedLicense)
	} else {
		allURL.License = license.LicenseName
		allURL.LicenseID = license.ID
		allURL.IsSpdx = license.IsSpdx
	}
	version, _ := m.ver.GetVersionByName(ctx, allURL.Version)
	if len(version.VersionName) == 0 {
		s.Warnf("No version details in DB for: %v", allURL.Version)
	}
	mineIDs, _ := m.mine.GetMineIdsByPurlType(ctx, purlType)
	if len(mineIDs) > 0 {
		allURL.MineID = mineIDs[0] // Assign the first mine id
	} else {
		s.Warnf("No mine details in DB for purl type: %v", purlType)
	}
	return allURL, nil
}

// queryPkgGoDev retrieves the latest information about a Golang Package from https://pkg.go.dev
func (m *GolangProjects) queryPkgGoDev(ctx context.Context, purlName, purlVersion string) (AllURL, error) {
	s := ctxzap.Extract(ctx).Sugar()
	if len(purlName) == 0 {
		s.Errorf("Please specify a valid Purl Name to query")
		return AllURL{}, errors.New("please specify a valid Purl Name to query")
	}
	client := pkggodevclient.New()
	pkg := purlName
	if len(purlVersion) > 0 {
		pkg = fmt.Sprintf("%s@%s", purlName, purlVersion)
	}
	s.Debugf("Checking pkg.go.dev for the latest info: %v", pkg)
	comp, err := client.DescribePackage(pkggodevclient.DescribePackageRequest{Package: pkg})
	if err != nil && len(purlVersion) > 0 {
		// We have a version zero search, so look for the latest one
		s.Debugf("Failed to query pkg.go.dev for %v: %v. Trying without version...", pkg, err)
		comp, err = client.DescribePackage(pkggodevclient.DescribePackageRequest{Package: purlName})
	}
	if err != nil {
		s.Warnf("Failed to query pkg.go.dev for %v: %v", pkg, err)
		return AllURL{}, fmt.Errorf("failed to query pkg.go.dev: %v", err)
	}
	var version = comp.Version
	if len(purlVersion) > 0 {
		version = purlVersion // Force the requested version if specified (the returned value can be concatenated)
	}
	allURL := AllURL{
		Component: purlName,
		Version:   version,
		License:   comp.License,
		PurlName:  purlName,
		URL:       fmt.Sprintf("https://%v", comp.Repository),
	}
	return allURL, nil
}

// CheckPurlByNameType checks the golang project table for the count of entries matching a Purl Name and Type.
func (m *GolangProjects) CheckPurlByNameType(ctx context.Context, purlName string, purlType string) (int, error) {
	s := ctxzap.Extract(ctx).Sugar()
	if len(purlName) == 0 {
		s.Error("Please specify a valid Purl Name to query")
		return -1, errors.New("please specify a valid Purl Name to query")
	}
	if len(purlType) == 0 {
		s.Error("Please specify a valid Purl Type to query")
		return -1, errors.New("please specify a valid Purl Type to query")
	}
	var count int
	fmt.Printf("Checking golang projects table for %v, %v...\n", purlName, purlType)
	err := m.db.QueryRowxContext(ctx,
		"SELECT count(*)"+
			" FROM golang_projects u"+
			" INNER JOIN mines m ON u.mine_id = m.id"+
			" WHERE u.purl_name = $1 AND m.purl_type = $2",
		purlName, purlType).Scan(&count)
	if err != nil {
		s.Errorf("Error: Failed to query projects table for %v, %v: %v", purlName, purlType, err)
		return -1, fmt.Errorf("failed to query the projects table: %v", err)
	}
	return count, nil
}
