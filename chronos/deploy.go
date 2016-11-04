package chronos

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

// Deploy deploys Chronos tasks for a given app, env and tag
func Deploy(app, env, tag string) error {
	chronosPath := viper.GetString("chronos_json_path")
	deployment := viper.GetStringMap("apps")[app]
	if deployment == nil {
		return fmt.Errorf("invalid YAML config for %s\n", app)
	}
	for k, v := range deployment.(map[string]interface{}) {
		if k == "chronos" {
			for _, x := range v.([]interface{}) {
				cj := path.Join(chronosPath, strings.Replace(x.(string), "{{env}}", env, -1))
				body, err := ioutil.ReadFile(cj)
				if err != nil {
					return fmt.Errorf("Chronos JSON does not exist %s: %s\n", cj, err)
				}
				re := regexp.MustCompile(fmt.Sprintf("(quay.io/betterdoctor/%s):.*(\",?)", app))
				chronosJSON := re.ReplaceAllString(string(body), fmt.Sprintf("$1:%s$2", tag))

				url := fmt.Sprintf("%s/service/chronos/scheduler/iso8601", viper.GetString("chronos_host"))
				client := &http.Client{}
				req, _ := http.NewRequest("POST", url, strings.NewReader(chronosJSON))
				req.Header.Set("Content-Type", "application/json")
				resp, err := client.Do(req)
				if err != nil {
					return err
				}
				if resp.StatusCode != http.StatusNoContent {
					return fmt.Errorf("failed to deploy chronos tasks for %s: %s", app, resp.Status)
				}
				fmt.Printf("deployed Chronos task %s\n", path.Base(cj))
			}
		}
	}
	return nil
}
