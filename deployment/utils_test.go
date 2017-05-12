package deployment

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spf13/viper"
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

func TestBeginDeploy(t *testing.T) {
	cases := []struct {
		app     string
		env     string
		allowed bool
	}{
		{app: "yo", env: "stage", allowed: true},
		{app: "badboy", env: "stage", allowed: false},
	}

	for _, test := range cases {
		ts := createConsulDeployLockServer(test.allowed)
		viper.Set("consul_host", ts.URL)
		err := BeginDeploy(test.app, test.env)
		if test.allowed && err != nil {
			t.Errorf("expected error to be nil got: %v", err)
		}
		if !test.allowed && err == nil {
			t.Errorf("expected error but got nil")
		}
	}
}

func TestFinishDeploy(t *testing.T) {
	cases := []struct {
		app     string
		env     string
		allowed bool
	}{
		{app: "yo", env: "stage", allowed: true},
		{app: "badboy", env: "stage", allowed: false},
	}

	for _, test := range cases {
		ts := createConsulDeployLockServer(test.allowed)
		viper.Set("consul_host", ts.URL)
		err := FinishDeploy(test.app, test.env)
		if test.allowed && err != nil {
			t.Errorf("expected error to be nil got: %v", err)
		}
		if !test.allowed && err == nil {
			t.Errorf("expected error but got nil")
		}
	}
}

func TestUpdateReleaseTags(t *testing.T) {
	cases := []struct {
		app  string
		env  string
		curr string
		prev string
		ok   bool
	}{
		{app: "yo", env: "stage", curr: "1.2.3", prev: "1.2.2", ok: true},
		{app: "badboy", env: "stage", curr: "4.5.6", prev: "4.5.7", ok: false},
	}

	for _, test := range cases {
		ts := createConsulDeployLockServer(test.ok)
		viper.Set("consul_host", ts.URL)
		err := UpdateReleaseTags(test.app, test.env, test.curr, test.prev)
		if test.ok && err != nil {
			t.Errorf("expected error to be nil but got: %v", err)
		}
		if !test.ok && err == nil {
			t.Errorf("expected error but got nil")
		}
	}
}

func TestCurrentTag(t *testing.T) {
	cases := []struct {
		app string
		env string
		tag string
		ok  bool
	}{
		{app: "yo", env: "stage", tag: "1.2.3", ok: true},
		{app: "badboy", env: "stage", tag: "4.5.6", ok: false},
	}

	for _, test := range cases {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !test.ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, test.tag)
		}))
		viper.Set("consul_host", ts.URL)
		tag, err := CurrentTag(test.app, test.env)
		if test.ok && err != nil {
			t.Errorf("expected error to be nil but got: %v", err)
		}
		if err == nil && test.tag != tag {
			t.Errorf("expected tag '%s' but got '%s'", test.tag, tag)
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
	test := cases[0]
	d := GithubDiffLink(test.app, test.prev, test.tag)
	exp := "no github_org set: cannot generate diff link"
	if d != exp {
		t.Errorf("expected '%s' got '%s'", exp, d)
	}

	org := "bar"
	viper.Set("github_org", org)
	for _, test := range cases {
		d := GithubDiffLink(test.app, test.prev, test.tag)
		if test.prev == test.tag {
			if d != "no changes" {
				t.Errorf("expected no changes got %s", d)
			}
		} else {
			dl := fmt.Sprintf("https://github.com/%s/%s/compare/%s...%s", org, test.app, test.prev, test.tag)
			if d != dl {
				t.Errorf("expected %s but got %s", dl, d)
			}
		}
	}
}

func createConsulDeployLockServer(ok bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !ok {
			w.WriteHeader(http.StatusConflict)
			io.WriteString(w, `{"Results":null,"Errors":[{"OpIndex":0,"What":"Permission denied"}]}`)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
}
