package = github.com/betterdoctor/duncan
TAG := $(shell git tag | sort -r | head -n 1)

.PHONY: install release test

install:
	go get -t -v ./...

release:
	mkdir -p release
	perl -p -i -e 's/{{VERSION}}/$(TAG)/g' cmd/version.go
	GOOS=darwin GOARCH=amd64 go build -o release/duncan-darwin-amd64 $(package)
	GOOS=linux GOARCH=amd64 go build -o release/duncan-linux-amd64 $(package)
	perl -p -i -e 's/$(TAG)/{{VERSION}}/g' cmd/version.go

test:
	go test -v -cover `go list ./... | grep -v vendor` | sed ''/PASS/s//`printf "\033[32mPASS\033[0m"`/'' | sed ''/FAIL/s//`printf "\033[31mFAIL\033[0m"`/''
