package = github.com/betterdoctor/duncan
TAG := $(shell git tag --sort=v:refname | tail -n 1)

.PHONY: install release test

release: deps
	mkdir -p release
	perl -p -i -e 's/{{VERSION}}/$(TAG)/g' cmd/version.go
	GOOS=darwin GOARCH=amd64 go build -o release/duncan-darwin-amd64 $(package)
	GOOS=linux GOARCH=amd64 go build -o release/duncan-linux-amd64 $(package)
	perl -p -i -e 's/$(TAG)/{{VERSION}}/g' cmd/version.go

deps:
	glide install

install:
	cp release/duncan-darwin-amd64 /usr/local/bin/duncan
	chmod +x /usr/local/bin/duncan

test: deps
	go test -v -cover `go list ./... | grep -v vendor` | sed ''/PASS/s//`printf "\033[32mPASS\033[0m"`/'' | sed ''/FAIL/s//`printf "\033[31mFAIL\033[0m"`/''
