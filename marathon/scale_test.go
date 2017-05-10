package marathon

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
)

var group = &Group{
	Apps: []*App{
		&App{
			ID:        "/foo-production/web",
			Instances: 1,
			Container: &Container{
				Docker: &Docker{Image: "foo:1.2.3"},
			},
		},
		&App{
			ID:        "/foo-production/worker",
			Instances: 1,
			Container: &Container{
				Docker: &Docker{Image: "foo:1.2.3"},
			},
		},
	},
}

func TestScaledMarathonJSON(t *testing.T) {
	cases := []struct {
		scale map[string]int
		valid bool
	}{
		{scale: map[string]int{"web": 2}, valid: true},
		{scale: map[string]int{"web": 2, "worker": 4}, valid: true},
		{scale: map[string]int{"worker": 4, "web": 2}, valid: true},
		{scale: map[string]int{"web": -1}, valid: false},
		{scale: map[string]int{"fuuuuuu": 23}, valid: false},
	}

	for _, test := range cases {
		mj, err := scaledMarathonJSON(group, test.scale)
		if test.valid && len(mj) == 0 {
			t.Errorf("expected %v argument to be invalid, got: %v", test.scale, err)
		}

		if len(mj) > 0 {
			ng := &Group{}
			if err := json.Unmarshal(mj, &ng); err != nil {
				t.Error(err)
			}

			for _, a := range ng.Apps {
				for proc, count := range test.scale {
					if a.ID == "/foo-production/"+proc {
						if count != a.Instances {
							t.Errorf("expected %v to scale instances to %v", test.scale, count)
						}
					}
				}
			}
		}
	}
}

func TestScale(t *testing.T) {
	ts := marathonScaleServer()
	viper.Set("marathon_host", ts.URL)
	rules := map[string]int{
		"web":    2,
		"worker": 2,
	}
	_, err := Scale(group, rules)
	if err != nil {
		t.Errorf("expected error to be nil got %s", err)
	}
}

func marathonScaleServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"deploymentId": "c0e7434c-df47-4d23-99f1-78bd78662231"}`)
	}))
}
