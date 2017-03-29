

# marathon
`import "github.com/betterdoctor/duncan/marathon"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [func Deploy(app, env, tag string) error](#Deploy)
* [func List(app, env string) error](#List)
* [type App](#App)
  * [func (a *App) IsApp(app string) bool](#App.IsApp)
  * [func (a *App) ReleaseTag() string](#App.ReleaseTag)
  * [func (a *App) UpdateReleaseTag(tag string)](#App.UpdateReleaseTag)
* [type Container](#Container)
* [type Deployment](#Deployment)
* [type DeploymentResponse](#DeploymentResponse)
* [type Docker](#Docker)
* [type Group](#Group)
* [type Groups](#Groups)
  * [func (gs *Groups) DisplayAppStatus(apps map[string]string, env string) error](#Groups.DisplayAppStatus)
* [type ScaleEvent](#ScaleEvent)
  * [func Scale(app, env string, procs []string) (ScaleEvent, error)](#Scale)


#### <a name="pkg-files">Package files</a>
[deploy.go](/src/github.com/betterdoctor/duncan/marathon/deploy.go) [list.go](/src/github.com/betterdoctor/duncan/marathon/list.go) [scale.go](/src/github.com/betterdoctor/duncan/marathon/scale.go) [structs.go](/src/github.com/betterdoctor/duncan/marathon/structs.go) [utils.go](/src/github.com/betterdoctor/duncan/marathon/utils.go) 





## <a name="Deploy">func</a> [Deploy](/src/target/deploy.go?s=538:577#L13)
``` go
func Deploy(app, env, tag string) error
```
Deploy deploys a given marathon app, env and tag

If group has already been deployed, JSON is fetched
from Marathon API and modified; this prevents any scale
events from being overwritten by JSON in the betterdoctor/mesos repo.

If the group has not been deployed already, Duncan will
look for JSON in the Mesos repo and modify and deploy that.



## <a name="List">func</a> [List](/src/target/list.go?s=164:196#L3)
``` go
func List(app, env string) error
```
List shows the list of applications duncan knows about




## <a name="App">type</a> [App](/src/target/structs.go?s=1756:2507#L61)
``` go
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
```
App represents a Marathon app










### <a name="App.IsApp">func</a> (\*App) [IsApp](/src/target/structs.go?s=3002:3038#L92)
``` go
func (a *App) IsApp(app string) bool
```
IsApp returns true if Docker image matches app name




### <a name="App.ReleaseTag">func</a> (\*App) [ReleaseTag](/src/target/structs.go?s=2559:2592#L77)
``` go
func (a *App) ReleaseTag() string
```
ReleaseTag returns the git tag for a given app




### <a name="App.UpdateReleaseTag">func</a> (\*App) [UpdateReleaseTag](/src/target/structs.go?s=2773:2815#L86)
``` go
func (a *App) UpdateReleaseTag(tag string)
```
UpdateReleaseTag updates the release tag of an app's Docker image




## <a name="Container">type</a> [Container](/src/target/structs.go?s=3238:3295#L98)
``` go
type Container struct {
    Docker *Docker `json:"docker"`
}
```
Container represents a Marathon container










## <a name="Deployment">type</a> [Deployment](/src/target/structs.go?s=3741:3790#L112)
``` go
type Deployment struct {
    ID string `json:"id"`
}
```
Deployment represents a Marathon deployment










## <a name="DeploymentResponse">type</a> [DeploymentResponse](/src/target/structs.go?s=3898:3965#L118)
``` go
type DeploymentResponse struct {
    ID string `json:"deploymentId"`
}
```
DeploymentResponse represents a Marathon deployment response
when updating/creating a new app/group










## <a name="Docker">type</a> [Docker](/src/target/structs.go?s=3343:3692#L103)
``` go
type Docker struct {
    Type           string                   `json:"type"`
    Image          string                   `json:"image"`
    ForcePullImage bool                     `json:"forcePullImage,omitempty"`
    Network        string                   `json:"network,omitempty"`
    PortMappings   []map[string]interface{} `json:"portMappings,omitempty"`
}
```
Docker represents Marathon Docker metadata










## <a name="Group">type</a> [Group](/src/target/structs.go?s=1638:1721#L55)
``` go
type Group struct {
    ID   string `json:"id"`
    Apps []*App `json:"apps,omitempty"`
}
```
Group represents a Marathon app or group definition










## <a name="Groups">type</a> [Groups](/src/target/structs.go?s=265:319#L8)
``` go
type Groups struct {
    Groups []Group `json:"groups"`
}
```
Groups represents a list of Marathon groups










### <a name="Groups.DisplayAppStatus">func</a> (\*Groups) [DisplayAppStatus](/src/target/structs.go?s=377:453#L13)
``` go
func (gs *Groups) DisplayAppStatus(apps map[string]string, env string) error
```
DisplayAppStatus returns the group for the given app




## <a name="ScaleEvent">type</a> [ScaleEvent](/src/target/scale.go?s=277:318#L8)
``` go
type ScaleEvent map[string]map[string]int
```
ScaleEvent represents a Marathon Group scale event

One or more containers within a Group may be scaled
up or down at once







### <a name="Scale">func</a> [Scale](/src/target/scale.go?s=425:488#L12)
``` go
func Scale(app, env string, procs []string) (ScaleEvent, error)
```
Scale increases or decreases number of running instances of
an application within a Marathon Group









- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
