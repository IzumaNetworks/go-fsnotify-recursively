BUILD_FOLDER=dist
BIN_FOLDER := $$(go env GOPATH)/bin
REPO := $$(go mod why | tail -n 1)
SEMVER := $$(git tag --sort=-version:refname | head -n 1)
BINARY_NAME=gorph

.PHONY: test

tidy:
	go mod tidy

clean:
	go clean

publish:
	GOPROXY=proxy.golang.org go list -m ${REPO}@${SEMVER}
 
build:
	go build -o ${BUILD_FOLDER}/${BINARY_NAME} ./cmd/gorph/ 

install: build
	cp -fv ${BUILD_FOLDER}/${BINARY_NAME} ${BIN_FOLDER}/${BINARY_NAME}

test:
	go test .

benchmark:
	go test -bench=. -count 5 -run=^#
