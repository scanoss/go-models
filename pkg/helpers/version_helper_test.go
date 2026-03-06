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

import "testing"

func TestSemverTogglePrefix(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "v1.0.0",
			expected: "1.0.0",
		},
		{
			input:    "1.0.0",
			expected: "v1.0.0",
		},
		{
			input:    "@tryghost/koenig-lexical@1.1.9",
			expected: "@tryghost/koenig-lexical@1.1.9",
		},
		{
			input:    "release-1758393970",
			expected: "release-1758393970",
		},
		{
			input:    "v697",
			expected: "697",
		},
		{
			input:    "v3.3.0.3",
			expected: "v3.3.0.3",
		},
		{
			input:    "v3.3.0.3-beta.1",
			expected: "v3.3.0.3-beta.1",
		},
		{
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		got := SemverTogglePrefix(tt.input)
		if got != tt.expected {
			t.Errorf("SemverTogglePrefix() = %v, want %v", got, tt.expected)
		}
	}
}
