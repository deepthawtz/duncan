package marathon

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// Groups represents a list of Marathon groups
type Groups struct {
	Groups []Group `json:"groups"`
}

// DisplayAppStatus returns the group for the given app
func (gs *Groups) DisplayAppStatus(apps map[string]string, env string) error {
	if env == "" {
		env = "stage|production"
	}
	for a := range apps {
		for _, g := range gs.Groups {
			envs := strings.Split(env, "|")
			for _, e := range envs {
				id := strings.Join([]string{a, e}, "-")
				if g.ID == "/"+id {
					fmt.Println(id)
					var data = make([][]string, 10)
					for _, x := range g.Apps {
						if x.ID == "" {
							continue
						}
						data = append(data, []string{strings.Split(x.ID, "/")[2], x.Release(), strconv.Itoa(x.Instances), fmt.Sprintf("%.2f", x.CPUs), strconv.Itoa(x.Mem)})
					}
					table := tablewriter.NewWriter(os.Stdout)
					table.SetHeader([]string{"ID", "Tag", "Instances", "CPU", "Mem MB"})
					table.AppendBulk(data)
					table.Render()
				}
			}
		}
	}

	return nil
}

// Group represents a Marathon app or group definition
//
// TODO: deploy all applications as groups
type Group struct {
	App        // TODO: can remove when feedback converted to a group
	Apps []App `json:"apps,omitempty"`
}

// App represents a Marathon app
type App struct {
	ID        string    `json:"id"`
	Container Container `json:"container"`
	Instances int       `json:"instances"`
	CPUs      float64   `json:"cpus"`
	Mem       int       `json:"mem"`
}

// Release returns the git tag for a given app
func (a *App) Release() string {
	p := strings.Split(a.Container.Docker.Image, ":")
	if len(p) != 2 {
		return "no tag!!!"
	}
	return p[1]
}

// Container represents a Marathon container
type Container struct {
	Docker Docker `json:"docker"`
}

// Docker represents Marathon Docker metadata
type Docker struct {
	Image string `json:"image"`
}

// Deployment represents a Marathon deployment
type Deployment struct {
	ID string `json:"id"`
}

// DeploymentResponse represents a Marathon deployment response
// when updating/creating a new app/group
type DeploymentResponse struct {
	ID string `json:"deploymentId"`
}
