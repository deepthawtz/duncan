package deployment

import (
	"fmt"
	"strings"
	"testing"
)

func TestMarathonGroupID(t *testing.T) {
	cases := []struct {
		app string
		env string
	}{
		{app: "foo", env: "stage"},
		{app: "foo", env: "production"},
	}

	for _, test := range cases {
		id := MarathonGroupID(test.app, test.env)
		exp := "/" + strings.Join([]string{test.app, test.env}, "-")
		if id != exp {
			t.Errorf("expected %s but got %s", exp, id)
		}
	}
}

func TestGithubDiffLink(t *testing.T) {
	cases := []struct {
		app  string
		prev string
		tag  string
	}{
		{app: "foo", prev: "v1.2.3", tag: "v1.2.4"},
		{app: "foo", prev: "v1.2.3", tag: "v1.2.3"},
		{app: "foo", prev: "v1.2.4", tag: "v1.2.3"},
	}

	for _, test := range cases {
		d := GithubDiffLink(test.app, test.prev, test.tag)
		if test.prev == test.tag {
			if d != "no changes" {
				t.Errorf("expected no changes got %s", d)
			}
		} else {
			dl := fmt.Sprintf("https://github.com/betterdoctor/%s/compare/%s...%s", test.app, test.prev, test.tag)
			if d != dl {
				t.Errorf("expected %s but got %s", dl, d)
			}
		}
	}
}
