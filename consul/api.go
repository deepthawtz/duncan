package consul

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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
	url += "&recurse"
	url += fmt.Sprintf("?token=%s", viper.GetString("consul_token"))
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
func Write(app, deployEnv, url string, kvs []string) (map[string]string, error) {
	u := EnvURL(app, deployEnv)
	env, err := Read(u)
	if err != nil {
		return nil, err
	}

	for _, kvp := range kvs {
		a := strings.Split(kvp, "=")
		for k, v := range env {
			val := strings.Join(a[1:], "")
			if k == a[0] && v != val {
				fmt.Printf("changing %s from %s => %s\n", k, v, val)
			}
		}
	}

	var txn []*TxnItem
	for _, kvp := range kvs {
		a := strings.Split(kvp, "=")
		val := strings.Join(a[1:], "")
		env[a[0]] = val
		txn = append(txn, &TxnItem{
			KV: &KVPair{
				Key:   fmt.Sprintf("env/%s/%s/%s", app, deployEnv, a[0]),
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
	fmt.Println(string(body))
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

	return env, nil
}

// Delete removes key/values from Consul by given keys
func Delete(url string, keys []string) error {
	token := viper.GetString("consul_token")
	for _, k := range keys {
		url += fmt.Sprintf("/%s", k)
		url += fmt.Sprintf("?token=%s", token)
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
	token := viper.GetString("consul_token")
	return fmt.Sprintf("https://%s/v1/kv/deploys/%s/%s/current?raw&token=%s", ch, app, env, token)
}

// DeploymentTagURL returns URL to PUT release tags to (current/previous)
func DeploymentTagURL(app, env string) string {
	ch := viper.GetString("consul_host")
	token := viper.GetString("consul_token")
	return fmt.Sprintf("https://%s/v1/kv/deploys/%s/%s?token=%s", ch, app, env, token)
}
