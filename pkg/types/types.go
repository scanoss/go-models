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

// Package types defines the public API types for the SCANOSS Go models library.
//
// This package contains all request and response types used by the client API.
// It is designed to be imported by both the client and service packages to
// avoid cyclic dependencies while maintaining a clean API surface.
//
// The types in this package represent the public API contract and should be
// considered stable. Breaking changes to these types may require a major
// version bump.
package types

// ComponentRequest represents a request to get component information.
type ComponentRequest struct {
	// Purl is the Package URL identifying the component.
	Purl string `json:"purl"`

	// Requirement specifies version constraints (e.g., ">=1.0.0", "^2.0.0").
	Requirement string `json:"requirement"`
}

// ComponentResponse represents the response containing component information.
type ComponentResponse struct {
	// Purl is the Package URL of the component (without version)
	Purl string `json:"purl"`

	// Version is the component version.
	Version string `json:"version"`
}

// ComponentVersionsResponse represents the response containing all versions for a component.
type ComponentVersionsResponse struct {
	// Purl is the Package URL of the component (without version)
	Purl string `json:"purl"`

	// Versions is the list of all available versions for the component.
	Versions []string `json:"versions"`
}

// SourcePurl represents the raw source-mine data used to build a source PURL
// for a component. It is the public, decoded form of the projects/mines join
// row used by the source PURL lookup.
type SourcePurl struct {
	// SourceMineID is the id of the source mine the component was sourced from.
	SourceMineID int32 `json:"source_mine_id"`
	// SourcePurlName is the PURL name on the source mine.
	SourcePurlName string `json:"source_purl_name"`
	// SourceVendor is the namespace/vendor on the source mine (may be empty).
	SourceVendor string `json:"source_vendor"`
	// MineName is the human-readable name of the source mine.
	MineName string `json:"mine_name"`
	// PurlType is the PURL type of the source mine (e.g., "github").
	PurlType string `json:"purl_type"`
	// RepositoryURL is the canonical repository URL exposed by the source mine.
	RepositoryURL string `json:"repository_url"`
}
