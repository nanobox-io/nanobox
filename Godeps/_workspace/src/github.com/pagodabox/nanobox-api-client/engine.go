// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package client

import (
	"encoding/json"
	"time"
)

//
type (

	// Engine represents a nanobox published project
	Engine struct {
		ActiveReleaseID   string     `json:"active_release_id"`
		CreatedAt         time.Time  `json:"created_at"`
		CreatorID         string     `json:"creator_id"`
		Downloads         int        `json:"downloads"`
		ID                string     `json:"id"`
		Name              string     `json:"name"`
		Official          bool       `json:"official"`
		RepositoriumKey 	string     `json:"repositorium_key"`
		RepositoriumUser  string     `json:"repositorium_user"`
		State             string     `json:"state"`
		UpdatedAt         time.Time  `json:"updated_at"`
		WarehouseUser     string     `json:"warehouse_user"`
		WarehouseKey      string     `json:"warehouse_key"`
	}

	// EngineCreateOptions represents all available options when creating a engine.
	EngineCreateOptions struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
)

// routes

// CreateEngine creates a new engine, with provided options
func CreateEngine(options *EngineCreateOptions) (*Engine, error) {

	b, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	var engine Engine
	return &engine, post(&engine, "/engines", string(b))
}

// GetEngine returns the specified engine
func GetEngine(userSlug, engineSlug string) (*Engine, error) {
	var engine Engine
	return &engine, get(&engine, "/engines/" + userSlug + "/" + engineSlug)
}
