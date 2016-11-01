<img src="https://s3.amazonaws.com/betterdoctor-images/1/SAS.svg" width="100"> duncan: Docker deployment tool
=============================================================================================================

Duncan is a Docker deployment tool which aims to be like [Tim Duncan](https://en.wikipedia.org/wiki/Tim_Duncan):
consistent, reliable, and un-flashy.

```
Usage:
  duncan [command]

Available Commands:
  deploy      deploy an application
  list        List applications
  scale       Scale an app process
  version     Print the version of duncan
```

### Getting Started

Download a [binary release](https://github.com/betterdoctor/duncan/releases)

**OR**

Requires [Golang](https://golang.org/)

```bash
# clone to $GOPATH
cd $GOPATH/src/github.com/betterdoctor
git clone git@github.com:betterdoctor/duncan.git && cd duncan
make release
make install

# confirm installation
duncan version
```
