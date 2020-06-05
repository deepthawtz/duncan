package = github.com/deepthawtz/duncan

.PHONY: install release test

GOOS=$(shell uname -s | awk '{print tolower($$0)}')
OUTPUT_PATH="/usr/local/bin/duncan"

build:
	goreleaser --rm-dist --skip-validate --skip-publish

release:
	goreleaser --rm-dist

install:
	cp dist/duncan_$(GOOS)_amd64/duncan $(OUTPUT_PATH)
	chmod +x $(OUTPUT_PATH)

test:
	go test -v -cover `go list ./... | grep -v vendor` | sed ''/PASS/s//`printf "\033[32mPASS\033[0m"`/'' | sed ''/FAIL/s//`printf "\033[31mFAIL\033[0m"`/''

docs:
	sh ./generate_docs.sh
