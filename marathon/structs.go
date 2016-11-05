package marathon

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/betterdoctor/duncan/deployment"
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
				id := deployment.MarathonGroupID(a, e)
				if g.ID == id {
					fmt.Println(id)
					var data = make([][]string, 10)
					for _, x := range g.Apps {
						if x.ID == "" {
							continue
						}
						data = append(data, []string{
							strings.Split(x.ID, "/")[2],
							x.ReleaseTag(),
							strconv.Itoa(x.Instances),
							fmt.Sprintf("%.2f", x.CPUs),
							strconv.Itoa(x.Mem),
							x.Version.Format("2006-01-02T15:04:05"),
						})
					}
					table := tablewriter.NewWriter(os.Stdout)
					table.SetHeader([]string{"ID", "Tag", "Instances", "CPU", "Mem MB", "Deployed At"})
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

// App represents a Marathon app
type App struct {
	ID           string                   `json:"id"`
	Instances    int                      `json:"instances"`
	CPUs         float64                  `json:"cpus"`
	Mem          int                      `json:"mem"`
	Cmd          string                   `json:"cmd,omitempty"`
	URIs         []string                 `json:"uris,omitempty"`
	Dependencies []string                 `json:"dependencies,omitempty"`
	Container    *Container               `json:"container"`
	Env          map[string]string        `json:"env,omitempty"`
	Labels       map[string]string        `json:"labels,omitempty"`
	HealthChecks []map[string]interface{} `json:"healthChecks,omitempty"`
	Version      time.Time                `json:"version,omitempty"`
}

// ReleaseTag returns the git tag for a given app
func (a *App) ReleaseTag() string {
	p := strings.Split(a.Container.Docker.Image, ":")
	if len(p) != 2 {
		return "no tag!!!"
	}
	return p[1]
}

// UpdateReleaseTag updates the release tag of an app's Docker image
func (a *App) UpdateReleaseTag(tag string) {
	image := strings.Split(a.Container.Docker.Image, ":")[0]
	a.Container.Docker.Image = strings.Join([]string{image, tag}, ":")
}

// IsApp returns true if Docker image matches app name
func (a *App) IsApp(app string) bool {
	re := regexp.MustCompile(fmt.Sprintf("(quay.io/betterdoctor/)?%s:?(.*)?", app))
	return re.MatchString(a.Container.Docker.Image)
}

// Container represents a Marathon container
type Container struct {
	Docker *Docker `json:"docker"`
}

// Docker represents Marathon Docker metadata
type Docker struct {
	Type           string                   `json:"type"`
	Image          string                   `json:"image"`
	ForcePullImage bool                     `json:"forcePullImage,omitempty"`
	Network        string                   `json:"network,omitempty"`
	PortMappings   []map[string]interface{} `json:"portMappings,omitempty"`
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
