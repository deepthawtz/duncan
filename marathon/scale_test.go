package marathon

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"
)

func TestScaledMarathonJSON(t *testing.T) {
	app, env, tag := "foo", "production", "v1.2.3"
	cases := []struct {
		scale []string
		valid bool
	}{
		{scale: []string{"web=2"}, valid: true},
		{scale: []string{"web=2", "worker=4"}, valid: true},
		{scale: []string{"worker=4", "web=2"}, valid: true},
		{scale: []string{"web=-1"}, valid: false},
		{scale: []string{"web=fuuuuuu"}, valid: false},
		{scale: []string{"fuuuuuu=23"}, valid: false},
	}

	for _, test := range cases {
		g := &Group{
			Apps: []*App{
				&App{
					ID:        "/foo-production/web",
					Instances: 1,
					Container: &Container{
						Docker: &Docker{Image: strings.Join([]string{app, tag}, ":")},
					},
				},
				&App{
					ID:        "/foo-production/worker",
					Instances: 1,
					Container: &Container{
						Docker: &Docker{Image: strings.Join([]string{app, tag}, ":")},
					},
				},
			},
		}
		mj, _ := scaledMarathonJSON(g, app, env, tag, test.scale)
		if test.valid && len(mj) == 0 {
			t.Errorf("expected %s argument to be invalid", test.scale)
		}

		if len(mj) > 0 {
			ng := &Group{}
			if err := json.Unmarshal(mj, &ng); err != nil {
				t.Error(err)
			}

			for _, a := range ng.Apps {
				for _, proc := range test.scale {
					s := strings.Split(proc, "=")
					proc := s[0]
					count, _ := strconv.Atoi(s[1])
					if a.ID == "/"+strings.Join([]string{app, env}, "-")+"/"+proc {
						if count != a.Instances {
							t.Errorf("expected %v to scale instances to %v", test.scale, count)
						}
					}
				}
			}
		}
	}
}
