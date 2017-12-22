package marathon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/betterdoctor/duncan/deployment"
)

// Deploy deploys a given marathon app, env and tag
func Deploy(app, env, tag string) (string, error) {
	var prev string
	groups, err := listGroups()
	if err != nil {
		return prev, err
	}

	for _, g := range groups.Groups {
		if g.ID == deployment.MarathonGroupID(app, env) {
			for _, a := range g.Apps {
				if a.IsApp(app) {
					prev = a.ReleaseTag()
					a.UpdateReleaseTag(tag)
				}
			}
			j, err := json.Marshal(g)
			if err != nil {
				return prev, err
			}
			client := &http.Client{}
			req, _ := http.NewRequest("PUT", updateGroupURL(), bytes.NewReader(j))
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				return prev, err
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				b, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return prev, err
				}
				return prev, fmt.Errorf("failed to deploy: %s\n%s\n", resp.Status, string(b))
			}
			d := &deploymentResponse{}
			if err := json.NewDecoder(resp.Body).Decode(d); err != nil {
				return prev, err
			}
			if err := deployment.Watch(d.ID); err != nil {
				return prev, err
			}

			return prev, nil
		}
	}

	return prev, nil
}
