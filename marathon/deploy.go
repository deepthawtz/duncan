package marathon

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

// Deploy deploys a given marathon app, env and tag
func Deploy(app, env, tag string) error {
	marathonPath := viper.GetString("marathon_json_path")
	deployment := viper.GetStringMap("apps")[app]
	if deployment == nil {
		return fmt.Errorf("invalid YAML config for %s\n", app)
	}
	for k, v := range deployment.(map[interface{}]interface{}) {
		if k.(string) == "marathon" {
			for _, x := range v.([]interface{}) {
				mj := marathonJSONPath(marathonPath, x.(string), env)
				body, err := ioutil.ReadFile(mj)
				if err != nil {
					return fmt.Errorf("Marathon JSON file does not exist: %s\n", mj)
				}
				marathonJSON := marathonJSON(string(body), app, tag)
				client := &http.Client{}
				req, _ := http.NewRequest("PUT", deploymentURL(), strings.NewReader(marathonJSON))
				req.Header.Set("Content-Type", "application/json")
				resp, err := client.Do(req)
				if err != nil {
					return err
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
