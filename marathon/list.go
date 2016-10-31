package marathon

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

// List shows the list of applications duncan knows about
func List(app, env string) error {
	groups, err := listGroups()
	if err != nil {
		return err
	}
	apps := viper.GetStringMapString("apps")
	if app != "" {
		apps = map[string]string{
			app: apps[app],
		}
	}
	if err := groups.DisplayAppStatus(apps, env); err != nil {
		return err
	}
	return nil
}

// Groups returns a Marathon groups API response
func listGroups() (*Groups, error) {
	url := fmt.Sprintf("%s/service/marathon/v2/groups", viper.GetString("marathon_host"))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	groups := &Groups{}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(b, &groups); err != nil {
		fmt.Println(b)
		return nil, err
	}

	return groups, nil
}
