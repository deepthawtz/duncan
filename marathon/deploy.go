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
func Deploy(app, env, tag string) error {
	groups, err := listGroups()
	if err != nil {
		return err
	}

	// check JSON to see if group has already been deployed
	for _, g := range groups.Groups {
		if g.ID == deployment.MarathonGroupID(app, env) {
			for _, a := range g.Apps {
				if a.IsApp(app) {
					a.UpdateReleaseTag(tag)
				}
			}
			j, err := json.Marshal(g)
			if err != nil {
				return err
			}
			client := &http.Client{}
			req, _ := http.NewRequest("PUT", deploymentURL(), bytes.NewReader(j))
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
				return fmt.Errorf("failed to deploy: %s\n%s\n", resp.Status, string(b))
			}
			d := &DeploymentResponse{}
			if err := json.NewDecoder(resp.Body).Decode(d); err != nil {
				return err
			}
			if err := waitForDeployment(d.ID); err != nil {
				return err
			}

			return nil
		}
	}

	return nil
}
