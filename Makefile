DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

all: deps
	go-bindata -o gotcha/assets.go assets/...
	go install ./...

test: test-deps
	go list ./... | xargs -n1 go test

release: release-deps
	gox ./...

deps:
	go get github.com/jteeuwen/go-bindata/...

test-deps:
	go get github.com/stretchr/testify

release-deps:
	go get github.com/mitchellh/gox

.PNONY: all test deps
