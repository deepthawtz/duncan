package marathon

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

func TestDeploymentURL(t *testing.T) {
	// marathon group JSON, truncated for brevity
	group := `{"id": "yo-dawg-group", "apps": [{"id": "web"}]}`

	base := fmt.Sprintf("%s/service/marathon/v2/", viper.GetString("marathon_host"))
	cases := []struct {
		in  string
		out string
	}{
		{in: group, out: base + "groups/"},
	}

	for _, test := range cases {
		url := deploymentURL()

		if url != test.out {
			t.Errorf("expected %s but got %s", test.out, url)
		}
	}
}

func TestMarathonJSONPath(t *testing.T) {
	cases := []struct {
		env string
		in  string
		out string
	}{
		{env: "stage", in: "/blah/yodawg-{{env}}.json", out: "/blah/yodawg-stage.json"},
		{env: "production", in: "/blah/yodawg-{{env}}.json", out: "/blah/yodawg-production.json"},
	}

	for _, test := range cases {
		mjp := marathonJSONPath("", test.in, test.env)
		if mjp != test.out {
			t.Errorf("expected %s but got %s", test.out, mjp)
		}
	}
}

func TestMarathonJSON(t *testing.T) {
	group := `{"id": "yo-dawg-group", "apps": [{"id": "web", "container": {"docker": {"image": "quay.io/betterdoctor/yodawg:v1.2.3"}}}, {"id": "worker", "container": {"docker": {"image": "quay.io/yo/yodawg:v1.2.3"}}}]}`
	app := `{"id": "yo-dawg", "container": { "docker": {"image": "quay.io/betterdoctor/yodawg:v1.2.3"}}}`
	cases := []struct {
		body string
		app  string
		tag  string
		out  string
	}{
		{body: group, app: "yodawg", tag: "release-3.2.1", out: "quay.io/betterdoctor/yodawg:release-3.2.1"},
		{body: app, app: "yodawg", tag: "release-3.2.1", out: "quay.io/betterdoctor/yodawg:release-3.2.1"},
	}

	for _, test := range cases {
		m := marathonJSON([]byte(test.body), test.app, test.tag)
		dj := &Group{}
		if err := json.Unmarshal([]byte(m), &dj); err != nil {
			t.Error(err)
		}

		for _, a := range dj.Apps {
			image := a.Container.Docker.Image
			if image != test.out {
				t.Errorf("expected %s but got %s", test.out, image)
			}
		}
	}
}
