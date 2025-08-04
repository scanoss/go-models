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

// Handle all interaction with the mines table

package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/scanoss/go-grpc-helper/pkg/grpc/database"
)

type MineModel struct {
	q *database.DBQueryContext
}

type Mine struct {
	ID       int32  `db:"id"`
	Name     string `db:"mine_name"`
	PurlType string `db:"purl_type"`
}

// NewMineModel creates a new instance of the 'Mine' Model.
func NewMineModel(q *database.DBQueryContext) *MineModel {
	return &MineModel{q: q}
}

// GetMineIdsByPurlType retrieves a list of the Purl Type IDs associated with the given Purl Type (string).
func (m *MineModel) GetMineIdsByPurlType(ctx context.Context, purlType string) ([]int32, error) {
	s := ctxzap.Extract(ctx).Sugar()
	if len(purlType) == 0 {
		s.Error("Please specify a Purl Type to query")
		return nil, errors.New("please specify a Purl Type to query")
	}
	var mines []Mine
	err := m.q.SelectContext(ctx, &mines,
		"SELECT id,mine_name,purl_type FROM mines WHERE purl_type = $1", purlType,
	)
	if err != nil {
		s.Errorf("Error: Failed to query mines table for %v: %v", purlType, err)
		return nil, fmt.Errorf("failed to query the mines table: %v", err)
	}
	if len(mines) > 0 {
		var mineIds []int32
		for _, mine := range mines {
			mineIds = append(mineIds, mine.ID)
		}
		return mineIds, nil
	}
	s.Error("No entries found in the mines table.")
	return nil, errors.New("no entry in mines table")
}
