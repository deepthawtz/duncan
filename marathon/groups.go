package marathon

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/betterdoctor/duncan/deployment"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
)

// Groups represents a list of Marathon groups
type Groups struct {
	Groups []Group `json:"groups"`
}

// DisplayAppStatus returns the group for the given app
func (gs *Groups) DisplayAppStatus(apps []string, env string) error {
	if env == "" {
		env = "stage|production"
	}
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	white := color.New(color.FgWhite, color.Bold).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	for _, a := range apps {
		for _, g := range gs.Groups {
			envs := strings.Split(env, "|")
			for _, e := range envs {
				id := deployment.MarathonGroupID(a, e)
				if g.ID == id {
					fmt.Println(green(id))
					var data = make([][]string, 10)
					for _, x := range g.Apps {
						if x.ID == "" {
							continue
						}
						var (
							host string
							ok   bool
						)
						host, ok = x.Labels["HAPROXY_0_VHOST"]
						if ok {
							host = fmt.Sprintf("https://%s", host)
						}
						data = append(data, []string{
							cyan(x.InstanceType()),
							white(x.ReleaseTag()),
							yellow(strconv.Itoa(x.Instances)),
							white(fmt.Sprintf("%.2f", x.CPUs)),
							white(strconv.Itoa(x.Mem)),
							white(x.Version.Format("2006-01-02T15:04:05")),
							cyan(host),
						})
					}
					table := tablewriter.NewWriter(os.Stdout)
					table.SetHeader([]string{"ID", "Tag", "Instances", "CPU", "Mem MB", "Deployed At", "Host"})
					table.AppendBulk(data)
					table.Render()
				}
			}
		}
	}

	return nil
}

// Group represents a Marathon app or group definition
type Group struct {
	ID   string `json:"id"`
	Apps []*App `json:"apps,omitempty"`
}

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

// AssertAppExistsInGroup checks if an app type exists in a Marathon group
// and returns an error if it does not
func AssertAppExistsInGroup(app, env, typ string) error {
	id := fmt.Sprintf("/%s-%s/%s", app, env, typ)
	g, err := GroupDefinition(app, env)
	if err != nil {
		return err
	}

	for _, a := range g.Apps {
		if a.ID == id {
			return nil
		}
	}

	return fmt.Errorf("could not find %s in Marathon", id)
}

// List shows the list of applications duncan knows about
func List(app, env string) error {
	groups, err := listGroups()
	if err != nil {
		return err
	}
	apps := viper.GetStringSlice("apps")
	if app != "" {
		apps = []string{app}
	}
	return groups.DisplayAppStatus(apps, env)
}

// listGroups returns all deployed Marathon groups
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
		return nil, fmt.Errorf("%s is not responding with valid JSON, perhaps a bad MARATHON_HOST", url)
	}

	return groups, nil
}

// updateGroupURL returns a Marathon API endpoint to update a group (deploy, scale)
func updateGroupURL() string {
	return fmt.Sprintf("%s/service/marathon/v2/groups/", viper.GetString("marathon_host"))
}

// deploymentResponse deserializes a successful Marathon deployment response
// This ID is used to check when the deployments is finished
type deploymentResponse struct {
	ID string `json:"deploymentId"`
}
