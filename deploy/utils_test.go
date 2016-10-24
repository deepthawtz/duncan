package deploy

import (
	"testing"

	consul "github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/testutil"
)

type configCallback func(c *consul.Config)

func makeClient(t *testing.T) (*consul.Client, *testutil.TestServer) {
	return makeClientWithConfig(t, nil, nil)
}

func makeClientWithConfig(t *testing.T, cb1 configCallback, cb2 testutil.ServerConfigCallback) (*consul.Client, *testutil.TestServer) {
	// Make client config
	conf := consul.DefaultConfig()
	if cb1 != nil {
		cb1(conf)
	}

	// Create server
	server := testutil.NewTestServerConfig(t, cb2)
	conf.Address = server.HTTPAddr

	// Create client
	client, err := consul.NewClient(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	return client, server
}

func TestUpdateTags(t *testing.T) {
	c, s := makeClient(t)
	defer s.Stop()
	kv := c.KV()

	cases := []struct {
		prev string
		curr string
		out  string
	}{
		{prev: "", curr: "v1.2.3", out: ""},
		{prev: "v1.2.3", curr: "v1.2.3", out: "v1.2.3"},
		{prev: "v1.2.3", curr: "v1.2.4", out: "v1.2.3"},
	}

	for _, test := range cases {
		_, _ = kv.Put(&consul.KVPair{Key: "deploys/yo/stage/previous", Value: []byte(test.prev)}, nil)
		prev, err := UpdateTags("yo", "stage", test.curr, kv)
		if err != nil {
			t.Error(err)
		}
		if prev != test.out {
			t.Errorf("expected %s but got %s", test.out, prev)
		}
	}
}

// func TestDiff(t *testing.T) {
// 	d := Diff("yo", "v1.2.3", "v1.2.3")
// 	if d != "re-deployment, no changes" {
// 		t.Error("expected no diff link")
// 	}
// 	y := []byte(`
// repos:
//   yo: yo
// `)

// 	viper.ReadConfig(bytes.NewBuffer(y))
// 	d = Diff("yo", "v1.2.3", "v1.2.4")
// 	if d != "https://github.com/betterdoctor/app/compare/v1.2.3...v1.2.4" {
// 		t.Error("expected diff link")
// 	}
// }
