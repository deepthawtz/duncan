

# marathon
`import "github.com/betterdoctor/duncan/marathon"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [func AssertAppExistsInGroup(app, env, typ string) error](#AssertAppExistsInGroup)
* [func Deploy(app, env, tag string) error](#Deploy)
* [func List(app, env string) error](#List)
* [func Scale(group *Group, rules map[string]int) (string, error)](#Scale)
* [type App](#App)
  * [func (a *App) InstanceType() string](#App.InstanceType)
  * [func (a *App) IsApp(app string) bool](#App.IsApp)
  * [func (a *App) ReleaseTag() string](#App.ReleaseTag)
  * [func (a *App) UpdateReleaseTag(tag string)](#App.UpdateReleaseTag)
* [type Container](#Container)
* [type Docker](#Docker)
* [type Group](#Group)
  * [func GroupDefinition(app, env string) (*Group, error)](#GroupDefinition)
* [type Groups](#Groups)
  * [func (gs *Groups) DisplayAppStatus(apps []string, env string) error](#Groups.DisplayAppStatus)


#### <a name="pkg-files">Package files</a>
[app.go](/src/github.com/betterdoctor/duncan/marathon/app.go) [deploy.go](/src/github.com/betterdoctor/duncan/marathon/deploy.go) [groups.go](/src/github.com/betterdoctor/duncan/marathon/groups.go) [scale.go](/src/github.com/betterdoctor/duncan/marathon/scale.go) 





## <a name="AssertAppExistsInGroup">func</a> [AssertAppExistsInGroup](/src/target/groups.go?s=2571:2626#L91)
``` go
func AssertAppExistsInGroup(app, env, typ string) error
```
AssertAppExistsInGroup checks if an app type exists in a Marathon group
and returns an error if it does not



## <a name="Deploy">func</a> [Deploy](/src/target/deploy.go?s=186:225#L4)
``` go
func Deploy(app, env, tag string) error
```
Deploy deploys a given marathon app, env and tag



## <a name="List">func</a> [List](/src/target/groups.go?s=2932:2964#L108)
``` go
func List(app, env string) error
```
List shows the list of applications duncan knows about



## <a name="Scale">func</a> [Scale](/src/target/scale.go?s=180:242#L2)
``` go
func Scale(group *Group, rules map[string]int) (string, error)
```
Scale increases or decreases number of running instances of
an application within a Marathon Group




## <a name="App">type</a> [App](/src/target/app.go?s=126:877#L3)
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










### <a name="App.InstanceType">func</a> (\*App) [InstanceType](/src/target/app.go?s=1005:1040#L20)
``` go
func (a *App) InstanceType() string
```
InstanceType returns the Marathon app name ("instance type" in our terminology)
e.g., web, worker, worker-special-sauce




### <a name="App.IsApp">func</a> (\*App) [IsApp](/src/target/app.go?s=1633:1669#L45)
``` go
func (a *App) IsApp(app string) bool
```
IsApp returns true if Docker image matches app name




### <a name="App.ReleaseTag">func</a> (\*App) [ReleaseTag](/src/target/app.go?s=1190:1223#L30)
``` go
func (a *App) ReleaseTag() string
```
ReleaseTag returns the git tag for a given app




### <a name="App.UpdateReleaseTag">func</a> (\*App) [UpdateReleaseTag](/src/target/app.go?s=1404:1446#L39)
``` go
func (a *App) UpdateReleaseTag(tag string)
```
UpdateReleaseTag updates the release tag of an app's Docker image




## <a name="Container">type</a> [Container](/src/target/app.go?s=1869:1926#L51)
``` go
type Container struct {
    Docker *Docker `json:"docker"`
}
```
Container represents a Marathon container










## <a name="Docker">type</a> [Docker](/src/target/app.go?s=1974:2323#L56)
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










## <a name="Group">type</a> [Group](/src/target/groups.go?s=1841:1924#L65)
``` go
type Group struct {
    ID   string `json:"id"`
    Apps []*App `json:"apps,omitempty"`
}
```
Group represents a Marathon app or group definition







### <a name="GroupDefinition">func</a> [GroupDefinition](/src/target/groups.go?s=1992:2045#L71)
``` go
func GroupDefinition(app, env string) (*Group, error)
```
GroupDefinition returns the deployed Marathon group definition





## <a name="Groups">type</a> [Groups](/src/target/groups.go?s=289:343#L9)
``` go
type Groups struct {
    Groups []Group `json:"groups"`
}
```
Groups represents a list of Marathon groups










### <a name="Groups.DisplayAppStatus">func</a> (\*Groups) [DisplayAppStatus](/src/target/groups.go?s=401:468#L14)
``` go
func (gs *Groups) DisplayAppStatus(apps []string, env string) error
```
DisplayAppStatus returns the group for the given app








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
