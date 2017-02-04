<img src="https://s3.amazonaws.com/betterdoctor-images/1/SAS.svg" width="100"> duncan: Docker deployment tool
=============================================================================================================

Duncan is a Docker deployment tool which aims to be like [Tim Duncan](https://en.wikipedia.org/wiki/Tim_Duncan):
consistent, reliable, and un-flashy.

```
Usage:
  duncan [command]

Available Commands:
  deploy      Deploy an application
  env         Manage Consul key/values (ENV vars) for an app
  list        List applications
  logs        Streams logs of your service
  run         run a one-off process inside a remote container
  scale       Scale an app process
  secrets     Manage Vault secrets (ENV vars) for an app
  version     Print the version of duncan
```

### Getting Started

#### Installation

Download a [binary release](https://github.com/betterdoctor/duncan/releases)

**OR** build from source

Requires [Golang](https://golang.org/)

```bash
# clone to $GOPATH
cd $GOPATH/src/github.com/betterdoctor
git clone git@github.com:betterdoctor/duncan.git && cd duncan
make
make install

# confirm installation
duncan version
```

#### Configuration

```bash
cp example_duncan.yml $HOME/.duncan.yml
# populate YAML w/ valid values (ask ops team for help if stuck)
```
