package marathon

import (
	"testing"
)

func TestInstanceType(t *testing.T) {
	cases := []struct {
		id  string
		exp string
	}{
		{id: "/bar-production/web", exp: "web"},
		{id: "/coolshit", exp: ""},
	}

	for _, test := range cases {
		a := &App{
			ID: test.id,
		}
		it := a.InstanceType()
		if it != test.exp {
			t.Errorf("expected '%s' but got '%s'", test.exp, it)
		}
	}
}

func TestReleaseTag(t *testing.T) {
	cases := []struct {
		image string
		tag   string
	}{
		{image: "quay.io/foo/bar:v1.2.3", tag: "v1.2.3"},
		{image: "quay.io/foo/bar", tag: "no tag!!!"},
		{image: "quay.io/foo/bar:release-1.2.3", tag: "release-1.2.3"},
	}
	for _, test := range cases {
		a := &App{
			Container: &Container{
				Docker: &Docker{Image: test.image},
			},
		}

		tag := a.ReleaseTag()
		if tag != test.tag {
			t.Errorf("expected %s but got %s", test.tag, tag)
		}
	}
}

func TestAppUpdateReleaseTag(t *testing.T) {
	cases := []struct {
		image string
		tag   string
	}{
		{image: "quay.io/betterdoctor/foo:v1.2.3", tag: "v1.3.0"},
		{image: "quay.io/betterdoctor/foo-service:release-1.2.3", tag: "release-1.3.0"},
	}

	for _, test := range cases {
		a := &App{
			Container: &Container{
				Docker: &Docker{Image: test.image},
			},
		}
		a.UpdateReleaseTag(test.tag)

		tag := a.ReleaseTag()
		if tag != test.tag {
			t.Errorf("expected %s but got %s", test.tag, tag)
		}
	}
}

func TestAppIsApp(t *testing.T) {
	cases := []struct {
		image string
		app   string
		exp   bool
	}{
		{image: "redis:3.2", app: "foo", exp: false},
		{image: "quay.io/yo/foo:3.2", app: "foo", exp: true},
		{image: "foo:3.2", app: "foo", exp: false},
	}

	for _, test := range cases {
		a := &App{
			Container: &Container{
				Docker: &Docker{Image: test.image},
			},
		}
		b := a.IsApp(test.app)
		if test.exp != b {
			t.Errorf("expected %v but got %v for %v", test.exp, b, test.image)
		}
	}
}
