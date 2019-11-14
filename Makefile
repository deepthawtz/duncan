package = github.com/betterdoctor/duncan

.PHONY: install release test

build:
	goreleaser --rm-dist --skip-validate --skip-publish

release:
	goreleaser --rm-dist

install:
	cp dist/duncan_darwin_amd64/duncan /usr/local/bin/duncan
	chmod +x /usr/local/bin/duncan

test:
	go test -v -cover `go list ./... | grep -v vendor` | sed ''/PASS/s//`printf "\033[32mPASS\033[0m"`/'' | sed ''/FAIL/s//`printf "\033[31mFAIL\033[0m"`/''

docs:
	sh ./generate_docs.sh
