// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2025 SCANOSS.COM
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

package scanoss

import (
	"context"
	"github.com/jmoiron/sqlx"

	"github.com/scanoss/go-grpc-helper/pkg/grpc/database"
	"github.com/scanoss/go-models/pkg/models"
	"github.com/scanoss/go-models/pkg/services"
	"go.uber.org/zap"
)

// Client provides a unified interface to SCANOSS data models and services.
type Client struct {
	ctx context.Context
	s   *zap.SugaredLogger
	q   *database.DBQueryContext
	db  *sqlx.DB //TODO: remove db *sqlx.DB once QueryRowxContext is implemented on database.DBQueryContext and used across pkg/models

	Models    *models.Models
	Component *services.ComponentService
}

// New creates a SCANOSS Model Client.
func New(ctx context.Context, s *zap.SugaredLogger, q *database.DBQueryContext, db *sqlx.DB) *Client {
	m := models.NewModels(ctx, s, q, db)

	//Initialize services
	component := services.NewComponentService(ctx, s, m)

	return &Client{
		ctx:       ctx,
		s:         s,
		Models:    m,
		Component: component,
	}
}
