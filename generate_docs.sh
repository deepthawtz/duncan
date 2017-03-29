#!/bin/bash

# list all the packges, trim out the vendor directory and any main packages,
# then strip off the package name
PACKAGES="$(go list -f "{{.Name}}:{{.ImportPath}}" ./... | grep -v -E "main:|vendor/|examples" | cut -d ":" -f 2)"

# loop over all packages generating all their documentation
for pkg in $PACKAGES; do
  echo "godoc2md $pkg > $GOPATH/src/$pkg/README.md"
  godoc2md $pkg > $GOPATH/src/$pkg/README.md
done
