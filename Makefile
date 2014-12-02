all: deps
	go-bindata -o gotcha/assets.go -ignore=.gitignore assets/...
	go-bindata -o assets/demo_app/assets.go -ignore=.gitignore -prefix assets/demo_app/ assets/demo_app/assets/...
	go-bindata -o assets/new_app/assets.go -ignore=.gitignore -prefix assets/new_app/ assets/new_app/assets/...
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
