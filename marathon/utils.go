package marathon

import (
	"fmt"
	"regexp"

	"github.com/spf13/viper"
)

// deploymentURL returns a Marathon API endpoint to deploy/scale
func deploymentURL() string {
	return fmt.Sprintf("%s/service/marathon/v2/groups/", viper.GetString("marathon_host"))
}

func marathonJSON(body []byte, app, tag string) string {
	re := regexp.MustCompile(fmt.Sprintf("(%s/%s):.*(\",?)", viper.GetString("docker_repo_prefix"), app))
	return re.ReplaceAllString(string(body), fmt.Sprintf("$1:%s$2", tag))
}
