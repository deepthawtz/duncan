package consul

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

const envJSONTemplate = `
[
  {
    "LockIndex": 0,
    "Key": "env/{{.App}}/{{.Env}}/FOO_ENABLED",
    "Flags": 0,
    "Value": "Kg==",
    "CreateIndex": 2255959,
    "ModifyIndex": 2255959
  },
  {
    "LockIndex": 0,
    "Key": "env/{{.App}}/{{.Env}}/BAR_ENABLED",
    "Flags": 0,
    "Value": "Kg==",
    "CreateIndex": 2255959,
    "ModifyIndex": 2255959
  }
]`

func TestEnvMap(t *testing.T) {
	data := map[string][]byte{
		"YODAWG_LEVEL": []byte("9000"),
		"FOO_ENABLED":  []byte("true"),
	}
	var kvs []KVPair
	for k, v := range data {
		key := fmt.Sprintf("env/yo/stage/%s", k)
		value := base64.StdEncoding.EncodeToString(v)
		kvp := &KVPair{
			Key:   key,
			Value: value,
		}
		kvs = append(kvs, *kvp)
	}

	m := envMap(kvs)

	if len(m) != len(data) {
		t.Errorf("expected %v got %v", len(data), len(m))
	}

	for _, kvp := range kvs {
		p := strings.Split(kvp.Key, "/")
		key := p[len(p)-1]
		value, _ := base64.StdEncoding.DecodeString(kvp.Value)
		if m[key] != string(value) {
			t.Errorf("expected %s but got %s", m[key], value)
		}
	}
}

func TestEnvURL(t *testing.T) {
	cases := []struct {
		app string
		env string
	}{
		{app: "foo", env: "stage"},
		{app: "foo", env: "production"},
	}
	ch := "https://consul.yodawg.com"
	viper.Set("consul_host", ch)

	for _, test := range cases {
		exp := fmt.Sprintf("%s/v1/kv/env/%s/%s", ch, test.app, test.env)
		u := EnvURL(test.app, test.env)
		if exp != u {
			t.Errorf("expected %s but got %s", exp, u)
		}
	}
}

func TestCurrentDeploymentTagURL(t *testing.T) {
	cases := []struct {
		app string
		env string
	}{
		{app: "foo", env: "stage"},
		{app: "foo", env: "production"},
	}
	ch := "https://consul.yodawg.com"
	token := "abc123"
	viper.Set("consul_host", ch)
	viper.Set("consul_token", token)

	for _, test := range cases {
		exp := fmt.Sprintf("%s/v1/kv/deploys/%s/%s/current?raw&token=%s", ch, test.app, test.env, token)
		u := CurrentDeploymentTagURL(test.app, test.env)
		if exp != u {
			t.Errorf("expected %s but got %s", exp, u)
		}
	}
}

// TestApp represents a test app
type TestApp struct {
	App, Env string
	exists   bool
}

func TestRead(t *testing.T) {
	apps := []TestApp{
		{App: "foo", Env: "stage", exists: true},
		{App: "foo", Env: "production", exists: false},
	}

	viper.Set("consul_token", "abc123")
	for _, app := range apps {
		ts := createConsulENVServer(app)
		env, _ := Read(ts.URL)
		if app.exists && len(env) == 0 {
			t.Errorf("expected populated ENV map but got %v", env)
		}
		if !app.exists && len(env) != 0 {
			t.Errorf("expected empty ENV map but got %v", env)
		}
	}
}

func createConsulENVServer(app TestApp) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.exists {
			t := template.Must(template.New("env_json").Parse(envJSONTemplate))
			err := t.Execute(w, app)
			if err != nil {
				panic(err)
			}

			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}
