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

// Client provides a unified interface to SCANOSS data models and operations.
type Client struct {
	ctx    context.Context
	logger *zap.SugaredLogger
	conn   *sqlx.Conn
	q      *database.DBQueryContext

	// Data access layer
	Models *models.DB

	// Business logic layer
	Component *services.ComponentService
}

func New(ctx context.Context, logger *zap.SugaredLogger, conn *sqlx.Conn, q *database.DBQueryContext) *Client {
	// Initialize data access layer
	models := models.NewDB(ctx, logger, conn, q)

	// Initialize business logic layer
	component := services.NewComponentService(ctx, logger, models)

	return &Client{
		ctx:       ctx,
		logger:    logger,
		conn:      conn,
		Models:    models,
		Component: component,
	}
}

// Close closes the client and releases any resources.
// This should be called when the client is no longer needed.
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
