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

package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/scanoss/go-models/pkg/models"
	"github.com/scanoss/go-models/pkg/types"
	purlutils "github.com/scanoss/go-purl-helper/pkg"
	"sort"
)

// ComponentService orchestrates component lookup logic using extracted business logic.
type ComponentService struct {
	models *models.Models
}

// NewComponentService creates a new ComponentService instance.
// Uses the Models wrapper to access all necessary data access methods.
func NewComponentService(models *models.Models) *ComponentService {
	return &ComponentService{
		models: models,
	}
}

// GetComponent retrieves component information based on PURL and requirements.
func (cs *ComponentService) GetComponent(ctx context.Context, req types.ComponentRequest) (types.ComponentResponse, error) {
	// TODO: Simplify component selection logic.
	// The code was inspired from scanoss.com/dependencies and heavily refactored

	if len(req.Purl) == 0 {
		return types.ComponentResponse{}, errors.New("please specify a valid purl to query")
	}

	purl, err := purlutils.PurlFromString(req.Purl)
	if err != nil {
		return types.ComponentResponse{}, fmt.Errorf("failed to parse purl: %w", err)
	}

	purlName, err := purlutils.PurlNameFromString(req.Purl) //Make sure we just have the bare minimum for a Purl Name
	if err != nil {
		return types.ComponentResponse{}, fmt.Errorf("failed to extract purl name: %w", err)
	}

	purlReq := req.Requirement

	// Extract an exact version from requirement if no version in PURL
	if len(purl.Version) == 0 && len(purlReq) > 0 {
		ver := purlutils.GetVersionFromReq(purlReq)
		if len(ver) > 0 {
			purl.Version = ver
			purlReq = ""
		}
	}

	var allUrls []models.AllURL
	if len(purl.Version) > 0 {
		allUrls, err = cs.models.AllUrls.GetURLsByPurlNameTypeVersion(ctx, purlName, purl.Type, purl.Version)
	} else {
		allUrls, err = cs.models.AllUrls.GetURLsByPurlNameType(ctx, purlName, purl.Type)
	}

	if err != nil {
		return types.ComponentResponse{}, err
	}

	allUrl, err := cs.pickOneUrl(ctx, allUrls, purlName, purl.Type, purlReq)
	if err != nil {
		return types.ComponentResponse{}, err
	}

	if len(allUrl.Version) == 0 {
		return types.ComponentResponse{}, fmt.Errorf("cannot find version for purl %s", req.Purl)
	}

	return types.ComponentResponse{
		Purl:    req.Purl,
		Version: allUrl.Version,
	}, nil
}

// pickOneUrl takes the potential matching component/versions and selects the most appropriate one.
func (cs *ComponentService) pickOneUrl(ctx context.Context, allUrls []models.AllURL, purlName, purlType, purlReq string) (models.AllURL, error) {
	s := ctxzap.Extract(ctx).Sugar()

	if len(allUrls) == 0 {
		s.Infof("No component match (in urls) found for %v, %v", purlName, purlType)
		return models.AllURL{}, nil
	}

	var c *semver.Constraints
	var urlMap = make(map[*semver.Version]models.AllURL)
	if len(purlReq) > 0 {
		s.Debugf("Building version constraint for %v: %v", purlName, purlReq)
		var err error
		c, err = semver.NewConstraint(purlReq)
		if err != nil {
			s.Warnf("Encountered an issue parsing version constraint string '%v' (%v,%v): %v", purlReq, purlName, purlType, err)
		}
	}
	s.Debugf("Checking versions...")
	for _, url := range allUrls {
		if len(url.SemVer) > 0 || len(url.Version) > 0 {
			v, err := semver.NewVersion(url.Version)
			if err != nil && len(url.SemVer) > 0 {
				s.Debugf("Failed to parse SemVer: '%v'. Trying Version instead: %v (%v)", url.Version, url.SemVer, err)
				v, err = semver.NewVersion(url.SemVer) // Semver failed, try the normal version
			}
			if err != nil {
				s.Warnf("Encountered an issue parsing version string '%v' (%v) for %v: %v. Using v0.0.0", url.Version, url.SemVer, url, err)
				v, err = semver.NewVersion("v0.0.0") // Semver failed, just use a standard version zero (for now)
			}
			if err == nil {
				if c == nil || c.Check(v) {
					_, ok := urlMap[v]
					if !ok {
						urlMap[v] = url // fits inside the constraint and hasn't already been stored
					}
				}
			}
		} else {
			s.Infof("Skipping match as it doesn't have a version: %#v", url)
		}
	}
	if len(urlMap) == 0 { // TODO should we return the latest version anyway?
		s.Warnf("No component match found for %v, %v after filter %v", purlName, purlType, purlReq)
		return models.AllURL{}, nil
	}
	var versions = make([]*semver.Version, len(urlMap))
	var vi = 0
	for version := range urlMap { // Save the list of versions so they can be sorted
		versions[vi] = version
		vi++
	}
	sort.Sort(semver.Collection(versions))
	version := versions[len(versions)-1] // Get the latest (acceptable) URL version
	s.Debugf("Sorted versions: %v. Highest: %v", versions, version)

	url, ok := urlMap[version] // Retrieve the latest accepted URL version
	if !ok {
		s.Errorf("Problem retrieving URL data for %v (%v, %v)", version, purlName, purlType)
		return models.AllURL{}, fmt.Errorf("failed to retrieve specific URL version: %v", version)
	}
	url.URL, _ = purlutils.ProjectUrl(purlName, purlType)

	s.Debugf("Selected version: %#v", url)
	return url, nil // Return the best component match
}
