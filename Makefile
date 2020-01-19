SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

GO11MODULES=on
APP?=application
REGISTRY?=stevenacoffman
COMMIT_SHA=$(shell git rev-parse --short HEAD)
VERSION=$(shell git rev-parse HEAD)
BUILD=$(shell date +%FT%T%z)
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

.PHONY: build
build: clean ## - Build the application
	@printf "\033[32m\xE2\x9c\x93 Building your code\n\033[0m"
	go build -o ${APP} main.go

.PHONY: run
run: build ## - Runs go run main.go
	@printf "\033[32m\xE2\x9c\x93 Running your code\n\033[0m"
	go run -race main.go

.PHONY: clean
clean: ## - Cleans the binary
	@printf "\033[32m\xE2\x9c\x93 Cleaning your code\n\033[0m"
	rm -rf ${APP}
	go get -u -v golang.org/x/tools/cmd/goimports
	gofmt -l -w -s .
	goimports -l -w .
	go mod tidy

.PHONY: test
test: ## - Runs go test with default values
	@printf "\033[32m\xE2\x9c\x93 Testing your code to find potential problems\n\033[0m"
	go test -v -count=1 -race ./...

.PHONY: cover
cover: test ## - Runs test coverage report
	@printf "\033[32m\xE2\x9c\x93 Running Code Test Coverage Report\n\033[0m"
	go test -count=1 -coverprofile=coverage.out
	go tool cover -html=coverage.out

.PHONY: lint
lint: clean ## - Lint the application code for problems and nits
	@printf "\033[32m\xE2\x9c\x93 Linting your code to find potential problems\n\033[0m"
	golangci-lint run

.PHONY: docker-build
docker-build:	## - Build the smallest secure golang docker image based on distroless static
	@printf "\033[32m\xE2\x9c\x93 Build the smallest and secured golang docker image based on distroless static\n\033[0m"
	docker build -f Dockerfile -t ${REGISTRY}/${APP}:${COMMIT_SHA} .

.PHONY: docker-build-no-cache
docker-build-no-cache:	## - Build the smallest secure golang docker image based on distroless static with no cache
	@printf "\033[32m\xE2\x9c\x93 Build the smallest and secured golang docker image based on scratch\n\033[0m"
	docker build --no-cache -f Dockerfile -t ${REGISTRY}/${APP}:${COMMIT_SHA} .

.PHONY: ls
ls: ## - List size docker images
	@printf "\033[32m\xE2\x9c\x93 Look at the size dude !\n\033[0m"
	@echo image ls ${REGISTRY}/${APP}
	docker image ls ${REGISTRY}/${APP}

.PHONY: docker-run
docker-run:	docker-build ## - Run the smallest and secured golang docker image based on distroless static nonroot
	@printf "\033[32m\xE2\x9c\x93 Run the smallest and secured golang docker image based on scratch\n\033[0m"
	docker run -p 127.0.0.1:8080:8080/tcp ${REGISTRY}/${APP}:${COMMIT_SHA}

.PHONY: docker-push
docker-push: docker-build ## - Pushes the docker image to registry
	docker push ${REGISTRY}/${APP}:${COMMIT_SHA}

.PHONY: help
## help: Prints this help message
help: ## - Show help message
	@printf "\033[32m\xE2\x9c\x93 usage: make [target]\n\n\033[0m"
	grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
