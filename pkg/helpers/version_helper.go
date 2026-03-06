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

package helpers

import (
	"strings"

	"github.com/Masterminds/semver/v3"
)

// SemverTogglePrefix returns the alternate version string by toggling the "v" prefix.
// If the version is valid semver and lacks a "v" prefix, it adds one. If it has a "v" prefix, it removes it.
// If the version is not valid semver, it returns it as-is.
func SemverTogglePrefix(version string) string {
	if len(version) == 0 {
		return version
	}
	if _, err := semver.NewVersion(version); err == nil {
		if version[0] != 'v' {
			return "v" + version
		}
		return strings.TrimLeft(version, "v")
	}
	return version
}
