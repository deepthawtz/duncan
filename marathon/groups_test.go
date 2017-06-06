package marathon

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

type testGroup struct {
	App  string
	Env  string
	Tag  string
	Apps []*testApp
}

type testApp struct {
	InstanceType string
}

func TestGroupDefinition(t *testing.T) {
	cases := []struct {
		app           string
		env           string
		tag           string
		instanceTypes []string
		ok            bool
	}{
		{app: "foo", env: "production", tag: "1.2.3", instanceTypes: []string{"web"}, ok: true},
		{app: "bar", env: "stage", tag: "4.5.6", instanceTypes: []string{"web", "worker"}, ok: true},
		{app: "bar", env: "stage", tag: "4.5.6", instanceTypes: []string{"web", "worker"}, ok: false},
	}
	for _, test := range cases {
		tg := &testGroup{
			App: test.app,
			Env: test.env,
			Tag: test.tag,
		}
		for _, it := range test.instanceTypes {
			tg.Apps = append(tg.Apps, &testApp{InstanceType: it})
		}
		ts := createMarathonGroupsServer(tg, "group.json.tmpl", test.ok)
		viper.Set("marathon_host", ts.URL)
		g, err := GroupDefinition(test.app, test.env)

		if test.ok && err != nil {
			t.Errorf("expected error to be nil but got %v", err)
		}
		if !test.ok && err == nil {
			t.Errorf("expected error but got nil")
		}
		if test.ok && len(tg.Apps) != len(g.Apps) {
			t.Errorf("expected %d apps but got %d", len(tg.Apps), len(g.Apps))
		}
	}
}

func TestList(t *testing.T) {
	cases := []struct {
		app           string
		env           string
		tag           string
		instanceTypes []string
		ok            bool
	}{
		{app: "foo", env: "production", tag: "1.2.3", instanceTypes: []string{"web"}, ok: true},
		{app: "bar", env: "stage", tag: "4.5.6", instanceTypes: []string{"web", "worker"}, ok: true},
		{app: "bar", env: "stage", tag: "4.5.6", instanceTypes: []string{"web", "worker"}, ok: false},
	}
	for _, test := range cases {
		tg := &testGroup{
			App: test.app,
			Env: test.env,
			Tag: test.tag,
		}
		for _, it := range test.instanceTypes {
			tg.Apps = append(tg.Apps, &testApp{InstanceType: it})
		}
		ts := createMarathonGroupsServer(tg, "groups.json.tmpl", test.ok)
		viper.Set("marathon_host", ts.URL)
		err := List(test.app, test.env)
		if test.ok && err != nil {
			t.Errorf("expected error to be nil but got %v", err)
		}
		if !test.ok && err == nil {
			t.Errorf("expected error but got nil")
		}
	}
}

func createMarathonGroupsServer(tg *testGroup, tmpl string, ok bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		f := filepath.Join("testdata", tmpl)
		b, err := ioutil.ReadFile(f)
		if err != nil {
			panic(err)
		}
		j := new(bytes.Buffer)
		template := template.Must(template.New("task_json").Parse(string(b)))
		err = template.Execute(j, tg)
		if err != nil {
			panic(err)
		}
		io.WriteString(w, j.String())
	}))
}
