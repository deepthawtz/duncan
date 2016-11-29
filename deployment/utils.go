package deployment

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/betterdoctor/duncan/consul"
	"github.com/spf13/viper"
)

// UpdateReleaseTags updates the deployment git tags in Consul KV registry
// `tags/{app}/{env}/current` points to the currently deployed tag
// `tags/{app}/{env}/previous` points to the previously deployed tag
//
// This structure allows for rollback if a previous tag exists
//
// Returns previously deployed git tag if one has been deployed
func UpdateReleaseTags(app, env, tag string) (string, error) {
	url := consul.CurrentDeploymentTagURL(app, env)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	prev := string(b)

	if prev == tag {
		return tag, nil
	}
	m := map[string]string{
		"current":  tag,
		"previous": prev,
	}

	token := viper.GetString("consul_token")
	for k, v := range m {
		url = strings.Join([]string{consul.DeploymentTagURL(app, env), k}, "/")
		url += fmt.Sprintf("?token=%s", token)
		client := &http.Client{}
		req, _ := http.NewRequest("PUT", url, strings.NewReader(v))
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return "", err
			}
			return "", fmt.Errorf("failed to update release tags in Consul: %s\n%s", resp.Status, string(b))
		}
	}
	return prev, nil
}

// CurrentTag returns the currently deployed git tag for an app and environment
func CurrentTag(app, env string) (string, error) {
	url := consul.CurrentDeploymentTagURL(app, env)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	return "", fmt.Errorf("could not fetch current release tag: %s", resp.Status)
}

// MarathonGroupID returns a Marathon Group id for an app and env
func MarathonGroupID(app, env string) string {
	return "/" + strings.Join([]string{app, env}, "-")
}

// GithubDiffLink returns a GitHub diff link to view deployment changes
func GithubDiffLink(app, prev, tag string) string {
	if prev == tag || prev == "" {
		return "no changes"
	}
	// TODO: handle if YAML is not filled out correctly
	repo := viper.GetStringMapString("repos")[app]
	return fmt.Sprintf("https://github.com/betterdoctor/%s/compare/%s...%s", repo, prev, tag)
}
