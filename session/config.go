package session

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/gorilla/securecookie"
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

// loadCookieConfig reads a json file and returns the configuration values. If the json file does not exist, it will be
// created.
func loadCookieConfig(file string) (ck cookieKeys, err error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		// file does not exist
		if err := genCookieConfigJsonFile(file); err != nil {
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

// genCookieConfigJsonFile generates a json file containing cookie authorisation and encryption keys.
func genCookieConfigJsonFile(file string) error {
	c := cookieKeys{
		keys{
			securecookie.GenerateRandomKey(32),
			securecookie.GenerateRandomKey(32),
		},
	}
	configJson, err := json.Marshal(c)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(file, configJson, 0644); err != nil {
		return err
	}
	return nil
}
