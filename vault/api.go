package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

// Secrets represents Vault key/value pairs for a prefix
type Secrets struct {
	KVPairs map[string]string `json:"data"`
}

// Read displays all key/value pairs at the given prefix if no key is given
// If a key is passed will just display the key/value pair for the key
func Read(app, env string) error {
	s, err := readSecrets(app, env)
	if err != nil {
		return err
	}
	if len(s.KVPairs) == 0 {
		return fmt.Errorf("no secrets for --app %s --env %s", app, env)
	}
	for k, v := range s.KVPairs {
		fmt.Printf("%s=%s\n", k, v)
	}
	return nil
}

// Write sets the given key/value pairs at the provided prefix
func Write(app, env string, kvs []string) error {
	s, err := readSecrets(app, env)
	if err != nil {
		return err
	}
	for i, kvp := range kvs {
		for k, v := range s.KVPairs {
			p := strings.Split(kvp, "=")
			if k == p[0] {
				fmt.Printf("changing %s from %s => %s\n", k, v, p[1])
				s.KVPairs[k] = p[1]
				kvs = append(kvs[:i], kvs[i+1:]...)
			}
		}
	}

	for _, kvp := range kvs {
		p := strings.Split(kvp, "=")
		s.KVPairs[p[0]] = p[1]
	}

	if err := updateSecrets(app, env, s); err != nil {
		return err
	}

	for k, v := range s.KVPairs {
		fmt.Printf("%s=%s\n", k, v)
	}

	return nil
}

// Delete removes a key/value pair from the prefix
func Delete(app, env, key string) error {
	s, err := readSecrets(app, env)
	if err != nil {
		return err
	}
	b := len(s.KVPairs)
	for k, _ := range s.KVPairs {
		if k == key {
			delete(s.KVPairs, k)
		}
	}
	if len(s.KVPairs) == b {
		fmt.Printf("key %s does not exist\n", key)
		return nil
	}
	if err := updateSecrets(app, env, s); err != nil {
		return err
	}
	fmt.Printf("deleted key %s\n", key)

	return nil
}

func readSecrets(app, env string) (*Secrets, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", vaultSecretURL(app, env), strings.NewReader(""))
	req.Header.Set("X-Vault-Token", viper.GetString("vault_token"))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to reach vault: %s\n%s", resp.Status, string(b))
	}
	s := &Secrets{}
	if err := json.NewDecoder(resp.Body).Decode(s); err != nil {
		return nil, err
	}
	return s, nil
}

func updateSecrets(app, env string, s *Secrets) error {
	j, err := json.Marshal(s.KVPairs)
	if err != nil {
		return err
	}
	client := &http.Client{}
	req, _ := http.NewRequest("POST", vaultSecretURL(app, env), bytes.NewReader(j))
	req.Header.Set("X-Vault-Token", viper.GetString("vault_token"))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to set vault secrets: %s\n%s", resp.Status, string(b))
	}
	return nil
}

func vaultSecretURL(app, env string) string {
	vaultHost := viper.GetString("vault_host")
	return fmt.Sprintf("https://%s/v1/%s", vaultHost, prefix(app, env))
}

func prefix(app, env string) string {
	return fmt.Sprintf("secret/%s/%s", app, env)
}
