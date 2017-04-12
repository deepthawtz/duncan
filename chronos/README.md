

# chronos
`import "github.com/betterdoctor/duncan/chronos"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [func Deploy(app, env, tag string) error](#Deploy)
* [func RunCommand(app, env, cmd string, follow bool) error](#RunCommand)
* [type Executor](#Executor)
* [type Framework](#Framework)
* [type SlaveTasks](#SlaveTasks)
* [type TaskVars](#TaskVars)


#### <a name="pkg-files">Package files</a>
[deploy.go](/src/github.com/betterdoctor/duncan/chronos/deploy.go) [task_json.go](/src/github.com/betterdoctor/duncan/chronos/task_json.go) 





## <a name="Deploy">func</a> [Deploy](/src/target/deploy.go?s=1000:1039#L38)
``` go
func Deploy(app, env, tag string) error
```
Deploy deploys Chronos tasks for a given app, env and tag



## <a name="RunCommand">func</a> [RunCommand](/src/target/deploy.go?s=2385:2441#L74)
``` go
func RunCommand(app, env, cmd string, follow bool) error
```
RunCommand spins up a Chronos task to run the given command and exits




## <a name="Executor">type</a> [Executor](/src/target/deploy.go?s=679:770#L25)
``` go
type Executor struct {
    ID        string `json:"id"`
    Directory string `json:"directory"`
}
```
Executor represents a completed executor on a Mesos slave










## <a name="Framework">type</a> [Framework](/src/target/deploy.go?s=467:616#L18)
``` go
type Framework struct {
    ID        string      `json:"id"`
    Name      string      `json:"name"`
    Executors []*Executor `json:"completed_executors"`
}
```
Framework represents a completed framework on a Mesos slave










## <a name="SlaveTasks">type</a> [SlaveTasks](/src/target/deploy.go?s=321:402#L13)
``` go
type SlaveTasks struct {
    Frameworks []*Framework `json:"completed_frameworks"`
}
```
SlaveTasks represents Mesos slave completed tasks










## <a name="TaskVars">type</a> [TaskVars](/src/target/deploy.go?s=818:916#L31)
``` go
type TaskVars struct {
    App, Env, Tag, Command, TaskName, DockerRepoPrefix, DockerConfURL string
}
```
TaskVars represents a one-off Chronos task














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)