

# mesos
`import "github.com/betterdoctor/duncan/mesos"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [type ContainerStatus](#ContainerStatus)
* [type Executor](#Executor)
* [type Framework](#Framework)
* [type IPAddress](#IPAddress)
* [type Logs](#Logs)
* [type NetworkInfo](#NetworkInfo)
* [type SlaveTasks](#SlaveTasks)
* [type Status](#Status)
* [type Task](#Task)
  * [func (t *Task) Duration() (float64, error)](#Task.Duration)
  * [func (t *Task) LogDirectory(host string) (string, error)](#Task.LogDirectory)
  * [func (t *Task) SlaveIP() (string, error)](#Task.SlaveIP)
* [type Tasks](#Tasks)
  * [func (t *Tasks) Len() int](#Tasks.Len)
  * [func (t *Tasks) Less(i, j int) bool](#Tasks.Less)
  * [func (t *Tasks) Swap(i, j int)](#Tasks.Swap)
  * [func (t *Tasks) TasksFor(name string) (*Tasks, error)](#Tasks.TasksFor)


#### <a name="pkg-files">Package files</a>
[slave.go](/src/github.com/betterdoctor/duncan/mesos/slave.go) [task.go](/src/github.com/betterdoctor/duncan/mesos/task.go) 


## <a name="pkg-constants">Constants</a>
``` go
const (
    // TaskFramework is the framework name used by duncan run
    TaskFramework = "chronos"

    // TaskRunning represents status for a running task
    TaskRunning = "TASK_RUNNING"

    // TaskStaging represents status for a staged task
    TaskStaging = "TASK_STAGING"

    // TaskFinished represents status for a successfully completed task
    TaskFinished = "TASK_FINISHED"

    // TaskFailed represents status for a failed task
    TaskFailed = "TASK_FAILED"

    // TaskKilled represents status for a prematurely killed task
    TaskKilled = "TASK_KILLED"
)
```




## <a name="ContainerStatus">type</a> [ContainerStatus](/src/target/task.go?s=4058:4141#L160)
``` go
type ContainerStatus struct {
    NetworkInfos []*NetworkInfo `json:"network_infos"`
}
```
ContainerStatus repesents Mesos task container status










## <a name="Executor">type</a> [Executor](/src/target/slave.go?s=587:678#L8)
``` go
type Executor struct {
    ID        string `json:"id"`
    Directory string `json:"directory"`
}
```
Executor represents a completed executor on a Mesos slave










## <a name="Framework">type</a> [Framework](/src/target/slave.go?s=277:524#L1)
``` go
type Framework struct {
    ID                 string      `json:"id"`
    Name               string      `json:"name"`
    Executors          []*Executor `json:"executors,omitempty"`
    CompletedExecutors []*Executor `json:"completed_executors,omitempty"`
}
```
Framework represents a completed framework on a Mesos slave










## <a name="IPAddress">type</a> [IPAddress](/src/target/task.go?s=4324:4380#L170)
``` go
type IPAddress struct {
    IP string `json:"ip_address"`
}
```
IPAddress repesents a Mesos task slave IP










## <a name="Logs">type</a> [Logs](/src/target/slave.go?s=721:768#L14)
``` go
type Logs struct {
    Data string `json:"data"`
}
```
Logs represents logs for a Mesos task










## <a name="NetworkInfo">type</a> [NetworkInfo](/src/target/task.go?s=4202:4277#L165)
``` go
type NetworkInfo struct {
    IPAddresses []*IPAddress `json:"ip_addresses"`
}
```
NetworkInfo repesents Mesos task container network info










## <a name="SlaveTasks">type</a> [SlaveTasks](/src/target/slave.go?s=68:212#L1)
``` go
type SlaveTasks struct {
    Frameworks          []*Framework `json:"frameworks"`
    CompletedFrameworks []*Framework `json:"completed_frameworks"`
}
```
SlaveTasks represents Mesos slave completed tasks










## <a name="Status">type</a> [Status](/src/target/task.go?s=3833:3999#L153)
``` go
type Status struct {
    State     string           `json:"state"`
    Timestamp float64          `json:"timestamp"`
    Container *ContainerStatus `json:"container_status"`
}
```
Status repesents Mesos task status










## <a name="Task">type</a> [Task](/src/target/task.go?s=1720:1940#L66)
``` go
type Task struct {
    ID          string    `json:"id"`
    FrameworkID string    `json:"framework_id"`
    SlaveID     string    `json:"slave_id"`
    State       string    `json:"state"`
    Statuses    []*Status `json:"statuses"`
}
```
Task repesents a Mesos task










### <a name="Task.Duration">func</a> (\*Task) [Duration](/src/target/task.go?s=2359:2401#L91)
``` go
func (t *Task) Duration() (float64, error)
```
Duration returns the duration a task took to complete




### <a name="Task.LogDirectory">func</a> (\*Task) [LogDirectory](/src/target/task.go?s=2873:2929#L109)
``` go
func (t *Task) LogDirectory(host string) (string, error)
```
LogDirectory returns the Mesos sandbox directory for a task




### <a name="Task.SlaveIP">func</a> (\*Task) [SlaveIP](/src/target/task.go?s=2000:2040#L75)
``` go
func (t *Task) SlaveIP() (string, error)
```
SlaveIP returns the IP of slave the task is running on




## <a name="Tasks">type</a> [Tasks](/src/target/task.go?s=702:753#L23)
``` go
type Tasks struct {
    Tasks []*Task `json:"tasks"`
}
```
Tasks represents Mesos tasks
used to deserialize a Mesos tasks API response










### <a name="Tasks.Len">func</a> (\*Tasks) [Len](/src/target/task.go?s=1142:1167#L43)
``` go
func (t *Tasks) Len() int
```
Len helps Tasks implement the sort.Interface




### <a name="Tasks.Less">func</a> (\*Tasks) [Less](/src/target/task.go?s=1243:1278#L48)
``` go
func (t *Tasks) Less(i, j int) bool
```
Less helps Tasks implement the sort.Interface




### <a name="Tasks.Swap">func</a> (\*Tasks) [Swap](/src/target/task.go?s=1604:1634#L61)
``` go
func (t *Tasks) Swap(i, j int)
```
Swap helps Tasks implement the sort.Interface




### <a name="Tasks.TasksFor">func</a> (\*Tasks) [TasksFor](/src/target/task.go?s=818:871#L28)
``` go
func (t *Tasks) TasksFor(name string) (*Tasks, error)
```
TasksFor returns running task for given task app, env, name








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
