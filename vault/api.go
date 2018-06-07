package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/betterdoctor/duncan/config"
	"github.com/betterdoctor/kit/notify"
	"github.com/spf13/viper"
)

// Secrets represents Vault key/value pairs for a prefix
type Secrets struct {
	KVPairs map[string]string `json:"data"`
}

// Read displays all key/value pairs at the given prefix if no key is given
// If a key is passed will just display the key/value pair for the key
func Read(url string) (*Secrets, error) {
	s, err := readSecrets(url)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Write sets the given key/value pairs at the provided prefix
func Write(url string, kvs []string, s *Secrets) (*Secrets, error) {
	if len(s.KVPairs) == 0 {
		fmt.Println("no secrets yet, bootstrapping new prefix")
		s.KVPairs = make(map[string]string)
		for _, kvp := range kvs {
			p := strings.Split(kvp, "=")
			s.KVPairs[p[0]] = strings.Join(p[1:], "=")
		}
		if err := updateSecrets(url, s); err != nil {
			return nil, err
		}
		return s, nil
	}

	changes := map[string][]string{}
	for _, kvp := range kvs {
		p := strings.Split(kvp, "=")
		key := p[0]
		val := strings.Join(p[1:], "=")
		for k, v := range s.KVPairs {
			if k == key && v != val {
				changes[k] = []string{v, val}
				fmt.Printf("changing %s from %s => %s\n", k, v, val)
			}
		}
		if _, ok := s.KVPairs[key]; !ok {
			changes[key] = []string{val}
		}
		s.KVPairs[key] = val
	}

	if len(changes) == 0 {
		fmt.Println("no secrets were changed")
		return s, nil
	}

	if err := updateSecrets(url, s); err != nil {
		return nil, err
	}

	p := strings.Split(url, "/")
	app, deployEnv := p[len(p)-2], p[len(p)-1]
	msg := config.Changes("secrets", changes)
	if msg == "" {
		return s, nil
	}
	if err := notify.Slack(viper.GetString("slack_webhook_url"), fmt.Sprintf("%s %s", app, deployEnv), msg); err != nil {
		return nil, err
	}

	return s, nil
}

// Delete removes a key/value pair from the prefix
func Delete(url string, keys []string, s *Secrets) (*Secrets, error) {
	changes := map[string][]string{}
	for _, key := range keys {
		for k := range s.KVPairs {
			if k == key {
				fmt.Printf("deleting %s\n", k)
				changes[k] = []string{}
				delete(s.KVPairs, k)
			}
		}
	}
	if err := updateSecrets(url, s); err != nil {
		return nil, err
	}
	p := strings.Split(url, "/")
	app, deployEnv := p[len(p)-2], p[len(p)-1]
	msg := config.Changes("secrets", changes)
	if msg == "" {
		return s, nil
	}
	if err := notify.Slack(viper.GetString("slack_webhook_url"), fmt.Sprintf("%s %s", app, deployEnv), msg); err != nil {
		return nil, err
	}

	return s, nil
}

func readSecrets(url string) (*Secrets, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, strings.NewReader(""))
	req.Header.Set("X-Vault-Token", viper.GetString("vault_token"))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("Read access to Vault secrets %s denied. Either the key does not exist or your token does not have permission to access it", url)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch secrets: %s (%s)", url, resp.Status)
	}
	s := &Secrets{}
	if err := json.NewDecoder(resp.Body).Decode(s); err != nil {
		return nil, err
	}
	return s, nil
}

func updateSecrets(url string, s *Secrets) error {
	j, err := json.Marshal(s.KVPairs)
	if err != nil {
		return err
	}
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewReader(j))
	req.Header.Set("X-Vault-Token", viper.GetString("vault_token"))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Write access to Vault secrets %s denied. Either the key does not exist or your token does not have permission to access it", url)
	}
	return nil
}

// SecretsURL returns the Vault API endpoint to GET and POST secrets
// for a given app and env
func SecretsURL(app, env string) string {
	vaultHost := viper.GetString("vault_host")
	return fmt.Sprintf("%s/v1/%s", vaultHost, prefix(app, env))
}

func prefix(app, env string) string {
	return fmt.Sprintf("secret/%s/%s", app, env)
}
