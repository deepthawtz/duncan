package marathon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/betterdoctor/duncan/deployment"
)

// CurrentTag return the currently deployed release tag
func CurrentTag(app, env, repo string) (string, error) {
	var tag string
	groups, err := listGroups()
	if err != nil {
		return tag, err
	}

	for _, g := range groups.Groups {
		if g.ID == deployment.MarathonGroupID(app, env) {
			for _, a := range g.Apps {
				if a.IsApp(repo) {
					tag := a.ReleaseTag()
					return tag, nil
				}
			}
		}
	}

	return tag, fmt.Errorf("")
}

// Deploy deploys a given marathon app, env and tag
func Deploy(app, env, tag, repo string) error {
	groups, err := listGroups()
	if err != nil {
		return err
	}

	for _, g := range groups.Groups {
		if g.ID == deployment.MarathonGroupID(app, env) {
			for _, a := range g.Apps {
				if a.IsApp(repo) {
					a.UpdateReleaseTag(tag)
				}
			}
			j, err := json.Marshal(g)
			if err != nil {
				return err
			}
			client := &http.Client{}
			req, _ := http.NewRequest("PUT", updateGroupURL(), bytes.NewReader(j))
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				b, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return err
				}
				return fmt.Errorf("failed to deploy: %s\n%s", resp.Status, string(b))
			}
			d := &deploymentResponse{}
			if err := json.NewDecoder(resp.Body).Decode(d); err != nil {
				return err
			}

			return deployment.Watch(d.ID)
		}
	}

	return fmt.Errorf("No Marathon group running for %s-%s", app, env)
}
