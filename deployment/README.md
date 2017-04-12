

# deployment
`import "github.com/betterdoctor/duncan/deployment"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [func CurrentTag(app, env string) (string, error)](#CurrentTag)
* [func GithubDiffLink(app, prev, tag string) string](#GithubDiffLink)
* [func MarathonGroupID(app, env string) string](#MarathonGroupID)
* [func UpdateReleaseTags(app, env, tag string) (string, error)](#UpdateReleaseTags)


#### <a name="pkg-files">Package files</a>
[utils.go](/src/github.com/betterdoctor/duncan/deployment/utils.go) 





## <a name="CurrentTag">func</a> [CurrentTag](/src/target/utils.go?s=1585:1633#L54)
``` go
func CurrentTag(app, env string) (string, error)
```
CurrentTag returns the currently deployed git tag for an app and environment



## <a name="GithubDiffLink">func</a> [GithubDiffLink](/src/target/utils.go?s=2241:2290#L77)
``` go
func GithubDiffLink(app, prev, tag string) string
```
GithubDiffLink returns a GitHub diff link to view deployment changes



## <a name="MarathonGroupID">func</a> [MarathonGroupID](/src/target/utils.go?s=2067:2111#L72)
``` go
func MarathonGroupID(app, env string) string
```
MarathonGroupID returns a Marathon Group id for an app and env



## <a name="UpdateReleaseTags">func</a> [UpdateReleaseTags](/src/target/utils.go?s=487:547#L10)
``` go
func UpdateReleaseTags(app, env, tag string) (string, error)
```
UpdateReleaseTags updates the deployment git tags in Consul KV registry
`tags/{app}/{env}/current` points to the currently deployed tag
`tags/{app}/{env}/previous` points to the previously deployed tag

### This structure allows for rollback if a previous tag exists
Returns previously deployed git tag if one has been deployed








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)