BIN_calendar := "./bin/calendar"
DOCKER_IMG_calendar="calendar:develop"

BIN_calendar_scheduler := "./bin/calendar_scheduler"
DOCKER_IMG_calendar_scheduler="calendar_scheduler:develop"

BIN_calendar_sender := "./bin/calendar_sender"
DOCKER_IMG_calendar_sender="calendar_sender:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

up: run-docker
	@echo "start app"

down: stop-docker
	@echo "stop app"

integration-tests: up
	go test -tags=integration ./tests/integration -v
	docker compose --project-directory ./deployments up -d

build:
	go build -v -o $(BIN_calendar) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_calendar_scheduler) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
	go build -v -o $(BIN_calendar_sender) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

run-docker:
	docker compose --project-directory ./deployments up -d

stop-docker:
	docker compose --project-directory ./deployments down

run: build
	$(BIN_calendar) -config ./configs/calendar_config-dev.yml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_calendar) \
		-f build/Dockerfile_calendar .

run-img: build-img
	docker run $(DOCKER_IMG_calendar)

version: build
	$(BIN_calendar) version
	$(BIN_calendar_scheduler) version
	$(BIN_calendar_sender) version

test:
	go test -race ./internal/... ./pkg/... -v

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.62.0

lint: install-lint-deps
	golangci-lint run ./...

generate: generate
	@echo "Run go:generate"
	@go generate ./...
	protoc api/*.proto --go_out=./internal/server/grpc/pb --go-grpc_out=./internal/server/grpc/pb

.PHONY: build run build-img run-img version test lint
