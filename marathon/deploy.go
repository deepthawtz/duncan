package marathon

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var marathonPath string

func init() {
	marathonPath = viper.GetString("marathon_json_path")
}

// Deploy deploys a given marathon app, env and tag
func Deploy(app, env, tag string) error {
	deployment := viper.GetStringMap("apps")[app]
	if deployment == nil {
		return fmt.Errorf("invalid YAML config for %s\n", app)
	}
	for k, v := range deployment.(map[interface{}]interface{}) {
		if k.(string) == "marathon" {
			for _, x := range v.([]interface{}) {
				mj := marathonJSONPath(x.(string), env)
				body, err := ioutil.ReadFile(mj)
				if err != nil {
					return fmt.Errorf("Marathon JSON file does not exist %s: %s\n", mj, err)
				}
				marathonJSON := marathonJSON(string(body), app, tag)
				url, err := deploymentURL(marathonJSON)
				if err != nil {
					return err
				}
				client := &http.Client{}
				req, _ := http.NewRequest("PUT", url, strings.NewReader(marathonJSON))
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

func waitForDeployment(id string) error {
	fmt.Println("Waiting for deploy....")
	go func() {
		for {
			for _, r := range `-\|/` {
				fmt.Printf("\r%c", r)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	url := fmt.Sprintf("%s/service/marathon/v2/deployments", viper.GetString("marathon_host"))

	for {
		time.Sleep(5 * time.Second)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		b, _ := ioutil.ReadAll(resp.Body)
		var d []Deployment
		if err := json.Unmarshal(b, &d); err != nil {
			return err
		}
		if len(d) == 0 {
			fmt.Printf("\rDONE\n")
			break
		}
	}

	return nil
}

// deploymentURL returns a Marathon API deployment URL to deploy
// a Marathon App or Marathon Group depending on the JSON
func deploymentURL(mj string) (string, error) {
	dj := &DeploymentJSON{}
	if err := json.Unmarshal([]byte(mj), &dj); err != nil {
		return "", err
	}
	if len(dj.Apps) == 0 {
		return fmt.Sprintf("%s/service/marathon/v2/apps/%s", viper.GetString("marathon_host"), dj.ID), nil
	}
	return fmt.Sprintf("%s/service/marathon/v2/groups/", viper.GetString("marathon_host")), nil
}

func marathonJSONPath(f, env string) string {
	return path.Join(marathonPath, strings.Replace(f, "{{env}}", env, -1))
}

func marathonJSON(body, app, tag string) string {
	re := regexp.MustCompile(fmt.Sprintf("(quay.io/betterdoctor/%s):.*(\",?)", app))
	return re.ReplaceAllString(string(body), fmt.Sprintf("$1:%s$2", tag))
}
