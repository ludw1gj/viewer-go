package session

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// keys contain key types for session CookieStore.
type keys struct {
	AuthorisationKey []byte `json:"authorisation_key"`
	EncryptionKey    []byte `json:"encryption_key"`
}

// cookieKeys contains the required configuration information.
type cookieKeys struct {
	Cookie keys `json:"cookie"`
}

// generateCookieStore initialises and returns a CookieStore.
func generateCookieStore(configJSONFile string) (*sessions.CookieStore, error) {
	// generateCookieConfigJSONFile generates a json file containing cookie authorisation and encryption keys.
	generateCookieConfigJSONFile := func(file string) error {
		c := cookieKeys{
			keys{
				securecookie.GenerateRandomKey(32),
				securecookie.GenerateRandomKey(32),
			},
		}
		configJSON, err := json.Marshal(c)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(file, configJSON, 0644); err != nil {
			return err
		}
		return nil
	}

	// loadCookieConfig reads a json file and returns the configuration values. If the json file does not exist, it will
	// be created.
	loadCookieConfig := func(file string) (ck cookieKeys, err error) {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			// file does not exist
			if err := generateCookieConfigJSONFile(file); err != nil {
				return ck, err
			}
		}

		configFile, err := os.Open(file)
		if err != nil {
			return ck, err
		}
		defer configFile.Close()

		jsonParser := json.NewDecoder(configFile)
		if err := jsonParser.Decode(&ck); err != nil {
			return ck, err
		}
		return ck, nil
	}

	ck, err := loadCookieConfig(configJSONFile)
	if err != nil {
		return nil, errors.New("Failed to initialise a CookieStore: " + err.Error())
	}
	store := sessions.NewCookieStore(ck.Cookie.AuthorisationKey, ck.Cookie.EncryptionKey)
	return store, nil
}
