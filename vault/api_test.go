package vault

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
)

func TestPrefix(t *testing.T) {
	p := prefix("foo", "stage")
	if p != "secret/foo/stage" {
		t.Errorf("expected secret/foo/stage but got %s", p)
	}
}

func TestSecretsURL(t *testing.T) {
	vh := viper.GetString("vault_host")
	u := SecretsURL("foo", "stage")
	exp := fmt.Sprintf("%s/v1/secret/foo/stage", vh)
	if u != exp {
		t.Errorf("expected %s but got %s", exp, u)
	}
}

func TestRead(t *testing.T) {
	ts := getSecretsServer(true)
	defer ts.Close()
	s, err := Read(ts.URL)
	if err != nil {
		t.Errorf("expected success but failed: %s", err)
	}
	if len(s.KVPairs) != 2 {
		t.Error("expected secrets to exits")
	}

	ts = getSecretsServer(false)
	defer ts.Close()
	s, err = Read(ts.URL)
	if err != nil {
		t.Errorf("expected nil but got %s", err)
	}

	ts = failServer()
	defer ts.Close()
	s, err = Read(ts.URL)
	if err == nil {
		t.Error("expected error but got nil")
	}
}

func TestWrite(t *testing.T) {
	cases := []struct {
		kvs   []string
		getOK bool
		setOK bool
		exp   int
	}{
		{kvs: []string{"SECRET_ONE=xxxxxxxxxxxx", "YABBA=doo"}, getOK: true, setOK: true, exp: 3},
		{kvs: []string{"FOO=bar", "YABBA=doo"}, getOK: true, setOK: true, exp: 4},
		{kvs: []string{"FOO=bar", "YABBA=doo"}, getOK: true, setOK: false, exp: 2},
		{kvs: []string{"VAR_WITH_EQUAL_SIGNS_IN_IT=bnv23=@#$"}, getOK: true, setOK: true, exp: 3},
		{kvs: []string{"YO=123"}, getOK: false, setOK: true, exp: 1},
	}

	for _, test := range cases {
		gss := getSecretsServer(test.getOK)
		defer gss.Close()
		sss := setSecretsServer(test.setOK)
		defer sss.Close()
		s, err := Read(gss.URL)
		if err != nil {
			t.Errorf("expected nil but got %s", err)
		}
		before := &Secrets{
			KVPairs: map[string]string{},
		}
		for k, v := range s.KVPairs {
			before.KVPairs[k] = v
		}
		se, err := Write(sss.URL, test.kvs, s)
		if test.setOK && err != nil {
			t.Errorf("expected success but failed: %s for %v", err, test.kvs)
		}
		if !test.setOK && err == nil {
			t.Errorf("expected error got nil for: %v", test.kvs)
		}
		if test.setOK && len(se.KVPairs) != test.exp {
			fmt.Println(test)
			t.Errorf("expected %d secrets but got %d for: %v", test.exp, len(se.KVPairs), se.KVPairs)
		}
	}
}

func TestDelete(t *testing.T) {
	gss := getSecretsServer(true)
	defer gss.Close()
	sss := setSecretsServer(true)
	defer sss.Close()
	s, err := Read(gss.URL)
	if err != nil {
		t.Errorf("expected nil but got %s", err)
	}
	se, err := Delete(sss.URL, []string{"SECRET_ONE", "SECRET_TWO"}, s)
	if err != nil {
		t.Errorf("expected success but failed: %s", err)
	}
	if err == nil && len(se.KVPairs) != 0 {
		t.Errorf("expected 0 secrets but got %v", se.KVPairs)
	}
}

func setSecretsServer(ok bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ok {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	}))
}

func failServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
}

func getSecretsServer(exist bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if exist {
			json := `
{
  "request_id": "8f80-84a694b421ec5-8d2a-ee88-4b51466",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 2764800,
  "data": {
    "SECRET_ONE": "ooooooooooooo",
    "SECRET_TWO": "my-precious"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}`

			fmt.Fprintln(w, json)
			return
		}
		fmt.Fprintln(w, `{"data": {}}`)
	}))
}
