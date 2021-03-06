

# k8s
`import "github.com/deepthawtz/duncan/k8s"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [type KubeAPI](#KubeAPI)
  * [func NewClient() (*KubeAPI, error)](#NewClient)
  * [func (k *KubeAPI) CurrentTag(app, env, repo string) (string, error)](#KubeAPI.CurrentTag)
  * [func (k *KubeAPI) Deploy(app, env, tag, repo string) error](#KubeAPI.Deploy)
  * [func (k *KubeAPI) List(app, env string) error](#KubeAPI.List)


#### <a name="pkg-files">Package files</a>
[client.go](/src/github.com/deepthawtz/duncan/k8s/client.go) [deploy.go](/src/github.com/deepthawtz/duncan/k8s/deploy.go) [list.go](/src/github.com/deepthawtz/duncan/k8s/list.go) 






## <a name="KubeAPI">type</a> [KubeAPI](/src/target/client.go?s=181:254#L13)
``` go
type KubeAPI struct {
    Client    kubernetes.Interface
    Namespace string
}

```
KubeAPI performs all the Kubernetes API operations







### <a name="NewClient">func</a> [NewClient](/src/target/client.go?s=298:332#L19)
``` go
func NewClient() (*KubeAPI, error)
```
NewClient returns a new KubeAPI client





### <a name="KubeAPI.CurrentTag">func</a> (\*KubeAPI) [CurrentTag](/src/target/deploy.go?s=361:428#L17)
``` go
func (k *KubeAPI) CurrentTag(app, env, repo string) (string, error)
```
CurrentTag fetches the currently deployed docker image tag for
given app and env if it exists. First checks Kubernetes Deployment API
and then Stateful Sets API




### <a name="KubeAPI.Deploy">func</a> (\*KubeAPI) [Deploy](/src/target/deploy.go?s=1186:1244#L46)
``` go
func (k *KubeAPI) Deploy(app, env, tag, repo string) error
```
Deploy updates docker image tag for a given k8s deployment




### <a name="KubeAPI.List">func</a> (\*KubeAPI) [List](/src/target/list.go?s=572:617#L27)
``` go
func (k *KubeAPI) List(app, env string) error
```
List displays k8s pods matching given app/env








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
