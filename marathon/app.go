package marathon

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/viper"
)

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

// InstanceType returns the Marathon app name ("instance type" in our terminology)
// e.g., web, worker, worker-special-sauce
func (a *App) InstanceType() string {
	x := strings.Split(a.ID, "/")
	if len(x) != 3 {
		return ""
	}
	id := x[len(x)-1]
	return id
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
	re := regexp.MustCompile(fmt.Sprintf("(%s/%s):(.*)?", viper.GetString("docker_repo_prefix"), app))
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
