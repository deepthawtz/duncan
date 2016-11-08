package consul

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

// KVPair represents an individual key/value pair
type KVPair struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

// Read returns ENV for given consul KV URL
func Read(url string) (map[string]string, error) {
	url += "?recurse"
	token := viper.GetString("consul_token")
	if token != "" {
		url += fmt.Sprintf("&token=%s", token)
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var env []KVPair
	if resp.StatusCode == http.StatusNotFound {
		m := envMap(env)
		return m, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch ENV: %s", resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &env)
	if err != nil {
		return nil, err
	}

	m := envMap(env)
	return m, nil
}

// Write sets ENV vars for a given KV URL
func Write(url string, kvs []string) (map[string]string, error) {
	env, err := Read(url)
	if err != nil {
		return nil, err
	}

	for _, kvp := range kvs {
		a := strings.Split(kvp, "=")
		for k, v := range env {
			if k == a[0] && v != a[1] {
				fmt.Printf("changing %s from %s => %s\n", k, v, a[1])
			}
		}
	}

	for _, kvp := range kvs {
		a := strings.Split(kvp, "=")
		env[a[0]] = a[1]

		client := &http.Client{}
		req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/%s", url, a[0]), strings.NewReader(fmt.Sprintf("%s", a[1])))
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
	}

	return env, nil
}

// Delete removes key/values from Consul by given keys
func Delete(url string, keys []string) error {
	for _, k := range keys {
		url += fmt.Sprintf("/%s", k)
		client := &http.Client{}
		req, _ := http.NewRequest("DELETE", url, strings.NewReader(""))
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
		fmt.Printf("deleted %s\n", k)
	}
	return nil
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
	ch := viper.GetString("consul_host")
	return fmt.Sprintf("https://%s/v1/kv/env/%s/%s", ch, app, env)
}

// CurrentDeploymentTagURL returns URL to fetch currently deployed tag
func CurrentDeploymentTagURL(app, env string) string {
	ch := viper.GetString("consul_host")
	return fmt.Sprintf("https://%s/v1/kv/deploys/%s/%s/current?raw", ch, app, env)
}

// DeploymentTagURL returns URL to PUT release tags to (current/previous)
func DeploymentTagURL(app, env string) string {
	ch := viper.GetString("consul_host")
	return fmt.Sprintf("https://%s/v1/kv/deploys/%s/%s", ch, app, env)
}
