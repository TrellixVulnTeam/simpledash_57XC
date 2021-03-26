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
