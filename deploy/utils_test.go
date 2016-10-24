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
