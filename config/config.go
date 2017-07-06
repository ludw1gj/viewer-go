// Package config contains configuration structure required, and Load function to load configuration file.
package config

import (
	"encoding/json"
	"os"
)

// Config is the required configuration information for database connection and session initialisation.
type Config struct {
	Database struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
	} `json:"database"`
	CookieStoreAuthKey string `json:"cookie_store_auth_key"`
}

// Load reads a json file and returns the configuration values in config struct.
func Load(filePath string) (c Config, err error) {
	configFile, err := os.Open(filePath)
	if err != nil {
		return c, err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&c)
	return c, nil
}
