package docker

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
)

func TestVerifyTagExists(t *testing.T) {
	cases := []struct {
		app    string
		tag    string
		exists bool
	}{
		{app: "foo", tag: "1.2.3", exists: true},
		{app: "bar", tag: "4.5.6", exists: false},
	}
	for _, test := range cases {
		ts := createQuayAPIServer(test.tag, test.exists)
		viper.Set("docker_registry_host", ts.URL)
		err := VerifyTagExists(test.app, test.tag)
		if test.exists && err != nil {
			t.Errorf("expected tag '%s' to exist", test.tag)
		}
		if !test.exists && err == nil {
			t.Errorf("did nt expect tag '%s' to exist", test.tag)
		}
	}
}

func createQuayAPIServer(tag string, exists bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		j := `{"tags":[{"name": "%s"}]}`
		io.WriteString(w, fmt.Sprintf(j, tag))
	}))
}
