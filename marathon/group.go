package marathon

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/betterdoctor/duncan/deployment"
	"github.com/spf13/viper"
)

// GroupDefinition returns the deployed Marathon group definition
func GroupDefinition(app, env string) (*Group, error) {
	url := fmt.Sprintf("%s/service/marathon/v2/groups/", viper.GetString("marathon_host"))
	url += fmt.Sprintf("%s", deployment.MarathonGroupID(app, env))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	group := &Group{}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(b, &group); err != nil {
		fmt.Println(b)
		return nil, err
	}

	return group, nil
}
