package marathon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/betterdoctor/duncan/deployment"
	"github.com/spf13/viper"
)

// Deploy deploys a given marathon app, env and tag
//
// If group has already been deployed, JSON is fetched
// from Marathon API and modified; this prevents any scale
// events from being overwritten by JSON in the betterdoctor/mesos repo.
//
// If the group has not been deployed already, Duncan will
// look for JSON in the Mesos repo and modify and deploy that.
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
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("failed to deploy: %s\n", resp.Status)
			}
			defer resp.Body.Close()
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

	// if we get here, app has not been deployed already
	marathonPath := viper.GetString("marathon_json_path")
	deployment := viper.GetStringMap("apps")[app]
	if deployment == nil {
		return fmt.Errorf("invalid YAML config for %s\n", app)
	}
	for k, v := range deployment.(map[interface{}]interface{}) {
		if k.(string) == "marathon" {
			for _, x := range v.([]interface{}) {
				mjp := marathonJSONPath(marathonPath, x.(string), env)
				body, err := ioutil.ReadFile(mjp)
				if err != nil {
					return fmt.Errorf("Marathon JSON file does not exist: %s\n", mjp)
				}
				mj := marathonJSON(body, app, tag)
				client := &http.Client{}
				req, _ := http.NewRequest("PUT", deploymentURL(), strings.NewReader(mj))
				req.Header.Set("Content-Type", "application/json")
				resp, err := client.Do(req)
				if err != nil {
					return err
				}
				if resp.StatusCode != http.StatusOK {
					return fmt.Errorf("failed to deploy: %s\n", resp.Status)
				}
				defer resp.Body.Close()
				d := &DeploymentResponse{}
				if err := json.NewDecoder(resp.Body).Decode(d); err != nil {
					return err
				}
				fmt.Printf("Deploying %s %s to %s (%s)\n", app, tag, env, d.ID)
				if err := waitForDeployment(d.ID); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
