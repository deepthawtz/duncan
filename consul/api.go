package consul

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/betterdoctor/duncan/config"
	"github.com/betterdoctor/kit/notify"
	"github.com/spf13/viper"
)

// TxnItem represents a KV element in a transaction
type TxnItem struct {
	KV *KVPair `json:"KV"`
}

// KVPair represents an individual key/value pair
type KVPair struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
	Verb  string `json:"Verb,omitempty"`
}

// Read returns ENV for given consul KV URL
func Read(url string) (map[string]string, error) {
	url += "?recurse"
	url += fmt.Sprintf("&token=%s", viper.GetString("consul_token"))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var env []KVPair
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("failed to find ENV at: %s", url)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch ENV: %s", resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &env); err != nil {
		return nil, err
	}

	m := envMap(env)
	return m, nil
}

// Write sets ENV vars for a given KV URL and prints what changed
//
// e.g.,
// changing FOO_LEVEL from 9 => 9000
// changing BAR_ENABLED from true => false
func Write(app, deployEnv, url string, kvs []string) (map[string]string, error) {
	changes := map[string][]string{}
	u := EnvURL(app, deployEnv)
	env, err := Read(u)
	if err != nil {
		return nil, err
	}

	var txn []*TxnItem
	for _, kvp := range kvs {
		a := strings.Split(kvp, "=")
		key := a[0]
		val := strings.Join(a[1:], "=")
		for k, v := range env {
			if k == key && v != val {
				changes[k] = []string{v, val}
				fmt.Printf("changing %s from %s => %s\n", k, v, val)
			}
		}
		if _, ok := env[key]; !ok {
			changes[key] = []string{val}
		}
		env[a[0]] = val
		txn = append(txn, &TxnItem{
			KV: &KVPair{
				Key:   fmt.Sprintf("env/%s/%s/%s", app, deployEnv, key),
				Value: base64.StdEncoding.EncodeToString([]byte(val)),
				Verb:  "set",
			},
		})
	}

	client := &http.Client{}
	body, err := json.Marshal(txn)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest("PUT", url, bytes.NewReader(body))
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
		return nil, fmt.Errorf("failed to set Consul KV: %s\n%s", resp.Status, string(b))
	}
	msg := config.Changes("env", changes)
	if msg == "" {
		return env, nil
	}
	if err := notify.Slack(viper.GetString("slack_webhook_url"), fmt.Sprintf("%s %s", app, deployEnv), msg); err != nil {
		return nil, err
	}

	return env, nil
}

// Delete removes key/values from Consul by given keys
func Delete(app, deployEnv, url string, keys []string) error {
	token := viper.GetString("consul_token")
	changes := map[string][]string{}
	for _, k := range keys {
		u := fmt.Sprintf("%s/%s?token=%s", url, k, token)
		client := &http.Client{}
		req, _ := http.NewRequest("DELETE", u, strings.NewReader(""))
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			return fmt.Errorf("failed to set Consul KV: %s\n%s", resp.Status, string(b))
		}
		changes[k] = []string{}
		fmt.Printf("deleted %s\n", k)
	}
	msg := config.Changes("env", changes)
	if msg == "" {
		return nil
	}
	if err := notify.Slack(viper.GetString("slack_webhook_url"), fmt.Sprintf("%s %s", app, deployEnv), msg); err != nil {
		return err
	}
	return nil
}

// CurrentTag returns the last deployed tag for app + env
func CurrentTag(app, env string) (string, error) {
	resp, err := http.Get(CurrentDeploymentTagURL(app, env))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return "", fmt.Errorf("ACL does not allow you to run comand for %s %s", app, env)
	}
	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("No app found %s %s", app, env)
	}

	tag, err := ioutil.ReadAll(resp.Body)
	return string(tag), err
}

func envMap(kvs []KVPair) map[string]string {
	m := make(map[string]string)
	for _, env := range kvs {
		p := strings.Split(env.Key, "/")
		key := p[len(p)-1]
		if key != "" {
			value, _ := base64.StdEncoding.DecodeString(env.Value)
			m[key] = string(value)
		}
	}
	return m
}

// EnvURL returns a Consul KV URL for an app/env
func EnvURL(app, env string) string {
	host := viper.GetString("consul_host")
	return fmt.Sprintf("%s/v1/kv/env/%s/%s", host, app, env)
}

// CurrentDeploymentTagURL returns URL to fetch currently deployed tag
func CurrentDeploymentTagURL(app, env string) string {
	host := viper.GetString("consul_host")
	token := viper.GetString("consul_token")
	return fmt.Sprintf("%s/v1/kv/deploys/%s/%s/current?raw&token=%s", host, app, env, token)
}

// TxnURL returns a Consul transaction (txn) URL
func TxnURL() string {
	host := viper.GetString("consul_host")
	token := viper.GetString("consul_token")
	return fmt.Sprintf("%s/v1/txn?token=%s", host, token)
}
