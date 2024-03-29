VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`
GITHASH=`git rev-parse HEAD`
GIT_BRANCH=`git branch --show-current`

GO=go
GOCOVER=$(GO) tool cover
GOTEST=$(GO) test

.PHONY: fmt dev lint vet test test-cover build-all-in-one build-cobra build-ksqlgrammar all

all:
	make fmt vet lint test

dev:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1

build:  
	touch internal/version.go
	go build -ldflags "all=-X github.com/thmeitz/ksqldb-migrate/internal.version=${VERSION} -X github.com/thmeitz/ksqldb-migrate/internal.build=${BUILD} -X github.com/thmeitz/ksqldb-migrate/internal.hash=${GITHASH}" -a -o ksqldb-migrate .

release:
	./build-release.sh github.com/thmeitz/ksqldb-migrate

test:
	$(GOTEST) -v ./... -short
	gosec -no-fail -fmt=sonarqube -out coverage/secreport.json ./...

test-cover:
	$(GOTEST) ./... -coverprofile=coverage/coverage.out
	$(GOCOVER) -func=coverage/coverage.out 
	$(GOCOVER) -html=coverage/coverage.out
	golangci-lint run ./... --verbose --no-config --out-format checkstyle > coverage/golangci-lint.out

vet:	## run go vet on the source files
	$(GO) vet ./...

doc:	## generate godocs and start a local documentation webserver on port 8085
	GO111MODULE=off godoc -notes=TODO -goroot=. -http=:8085 -index

lint:
	golangci-lint run

clean-compose:	
	docker-compose down && docker-compose up -d

fmt: 
	$(GO) fmt ./...

changelog:
	git-chglog --output CHANGELOG.md