package deploy

import (
	"fmt"

	consul "github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

// UpdateTags updates the deployment git tags in Consul KV registry
// `tags/{app}/{env}/current` points to the currently deployed tag
// `tags/{app}/{env}/previous` points to the previously deployed tag
//
// This structure allows for rollback if a previous tag exists
// Returns previously deployed git tag if one has been deployed
func UpdateTags(app, env, tag string, client *consul.KV) (string, error) {
	prefix := fmt.Sprintf("deploys/%s/%s/", app, env)
	if client == nil {
		client = newConsulClient()
	}
	previous, _, err := client.Get(prefix+"current", nil)
	if err != nil {
		return "", err
	}

	if previous != nil && string(previous.Value) == tag {
		return tag, nil
	}

	curr := &consul.KVPair{
		Key:   prefix + "current",
		Value: []byte(tag),
	}
	_, err = client.Put(curr, nil)
	if err != nil {
		return "", err
	}

	if previous != nil {
		prev := &consul.KVPair{
			Key:   prefix + "previous",
			Value: previous.Value,
		}
		_, err = client.Put(prev, nil)
		if err != nil {
			return "", err
		}
	}

	return string(previous.Value), nil
}

// Diff returns a GitHub link to view the git diff of changes
// being deployed
func Diff(app, prev, tag string) string {
	if prev == tag {
		return "re-deployment, no changes"
	}
	// TODO: handle if YAML is not filled out correctly
	repo := viper.GetStringMapString("repos")[app]
	return fmt.Sprintf("https://github.com/betterdoctor/%s/compare/%s...%s", repo, prev, tag)
}

func newConsulClient() *consul.KV {
	config := consul.Config{
		Address: viper.GetString("consul_host"),
		Scheme:  "https",
	}
	client, err := consul.NewClient(&config)
	if err != nil {
		panic(err)
	}

	return client.KV()
}
