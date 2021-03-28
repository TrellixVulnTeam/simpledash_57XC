/*
	Copyright 2021 Arsen Musayelyan

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package main

// Conf struct stores root of TOML config
type Conf struct {
	Title         string
	Session       SessionConf
	AllowProxy    []string
	LoginRequired bool
	Theme         string
	Users         map[string]User
}

// SessionConf stores session configuration
type SessionConf struct {
	Name string
}

// User stores user configuration from TOML
type User struct {
	PasswordHash string
	ShowPublic   bool
	Cards        []Card `toml:"card"`
}

// Card stores card configuration from TOML
type Card struct {
	Type        string                 `toml:"type"`
	Title       string                 `toml:"title"`
	Description string                 `toml:"desc,omitempty"`
	Icon        string                 `toml:"icon,omitempty"`
	URL         string                 `toml:"url,omitempty"`
	Data        map[string]interface{} `toml:"data,omitempty"`
}
