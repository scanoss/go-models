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
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/scanoss/go-models/pkg/models"
	"github.com/scanoss/go-models/pkg/types"
	purlutils "github.com/scanoss/go-purl-helper/pkg"
)

var errGolangNotResolved = errors.New("golang component not fully resolved")

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

func (cs *ComponentService) CheckPurl(ctx context.Context, p string) (int, error) {
	if len(p) == 0 {
		return -1, errors.New("please specify a valid purl to query")
	}

	purl, err := purlutils.PurlFromString(p)
	if err != nil {
		return -1, fmt.Errorf("failed to parse purl: %w", err)
	}

	purlName, err := purlutils.PurlNameFromString(p) // Make sure we just have the bare minimum for a Purl Name
	if err != nil {
		return -1, fmt.Errorf("failed to extract purl name: %w", err)
	}

	if purl.Type == "golang" {
		return cs.models.GolangProjects.CheckPurlByNameType(ctx, purlName, purl.Type)
	}

	return cs.models.Projects.CheckPurlByNameType(ctx, purlName, purl.Type)
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

	purlName, err := purlutils.PurlNameFromString(req.Purl) // Make sure we just have the bare minimum for a Purl Name
	if err != nil {
		return types.ComponentResponse{}, fmt.Errorf("failed to extract purl name: %w", err)
	}

	purlReq := req.Requirement
	if len(purlReq) > 0 && len(purl.Version) > 0 {
		return types.ComponentResponse{}, errors.New("cannot specify both a version and a requirement")
	}

	// Extract an exact version from requirement if no version in PURL
	if len(purl.Version) == 0 && len(purlReq) > 0 {
		ver := purlutils.GetVersionFromReq(purlReq)
		if len(ver) > 0 {
			purl.Version = ver
			purlReq = ""
		}
	}
	if purl.Type == "golang" {
		resolved, golangErr := cs.resolveGolangComponent(ctx, req.Purl, purlReq)
		if golangErr != nil && !errors.Is(golangErr, errGolangNotResolved) {
			return types.ComponentResponse{}, golangErr
		}
		if resolved != nil {
			return types.ComponentResponse{
				Purl:    req.Purl,
				Version: resolved.Version,
			}, nil
		}
		// No component/license found — if it's a GitHub component, try GitHub lookup
		if strings.HasPrefix(req.Purl, "pkg:golang/github.com/") {
			purl.Type, purlName, purl.Version, err = convertGolangToGithubPurl(req.Purl)
			if err != nil {
				return types.ComponentResponse{}, err
			}
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

	allUrl, err := models.PickOneUrl(ctx, allUrls, purlName, purl.Type, purlReq)
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

// resolveGolangComponent looks up a golang component in the golang_projects table.
// Returns a non-nil AllURL if the component was fully resolved (has both component and license data),
// nil if the caller should fall through to all_urls lookup.
func (cs *ComponentService) resolveGolangComponent(ctx context.Context, purlString, purlReq string) (*models.AllURL, error) {
	s := ctxzap.Extract(ctx).Sugar()
	allURL, err := cs.models.GolangProjects.GetGoLangURLByPurlString(ctx, purlString, purlReq)
	if err != nil {
		return nil, err
	}
	if len(allURL.Component) == 0 {
		s.Debugf("Didn't find component in golang projects table for %v. Checking all urls...", purlString)
		return nil, errGolangNotResolved
	}
	if len(allURL.License) == 0 {
		s.Debugf("Didn't find license in golang projects table for %v. Checking all urls...", purlString)
		return nil, errGolangNotResolved
	}
	return &allURL, nil
}

// convertGolangToGithubPurl converts a golang GitHub purl string to a GitHub purl type,
// returning the parsed purl components needed for all_urls lookup.
func convertGolangToGithubPurl(purlString string) (string, string, string, error) {
	ghPurlString := purlutils.ConvertGoPurlStringToGithub(purlString)
	purl, err := purlutils.PurlFromString(ghPurlString)
	if err != nil {
		return "", "", "", err
	}
	purlName, err := purlutils.PurlNameFromString(ghPurlString)
	if err != nil {
		return "", "", "", err
	}
	return purl.Type, purlName, purl.Version, nil
}
