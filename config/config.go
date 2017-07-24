// Package config contains configuration structure required, and Load function to load configuration file.
package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/gorilla/securecookie"
)

// Keys contain key types for session CookieStore.
type Keys struct {
	AuthorisationKey []byte `json:"authorisation_key"`
	EncryptionKey    []byte `json:"encryption_key"`
}

// Config contains the required configuration information.
type Config struct {
	Cookie Keys `json:"cookie"`
}

// Load reads a json file and returns the configuration values in config struct. If the json file does not exist, it
// will be created.
func Load(file string) (c Config, err error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		// file does not exist
		err = genCookieConfigJsonFile(file)
		if err != nil {
			return c, err
		}
	}

	configFile, err := os.Open(file)
	if err != nil {
		return c, err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&c)
	if err != nil {
		return c, err
	}
	return c, nil
}

// genCookieConfigJsonFile generates a json file containing cookie authorisation and encryption keys.
func genCookieConfigJsonFile(file string) error {
	c := Config{
		Keys{
			securecookie.GenerateRandomKey(32),
			securecookie.GenerateRandomKey(32),
		},
	}
	configJson, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, configJson, 0644)
	if err != nil {
		return err
	}
	return nil
}
