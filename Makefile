package = github.com/betterdoctor/duncan

.PHONY: install release test

release: deps
	goreleaser --rm-dist

deps:
	glide install

install:
	cp dist/duncan_*_darwin_amd64/duncan /usr/local/bin/duncan
	chmod +x /usr/local/bin/duncan

test: deps
	go test -v -cover `go list ./... | grep -v vendor` | sed ''/PASS/s//`printf "\033[32mPASS\033[0m"`/'' | sed ''/FAIL/s//`printf "\033[31mFAIL\033[0m"`/''
