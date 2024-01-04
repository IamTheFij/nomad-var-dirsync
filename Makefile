APP_NAME = nomad-var-dirsync
DOCKER_TAG ?= $(APP_NAME)-USER
VERSION ?= $(shell git describe --tags --dirty)
GOFILES = *.go
# Multi-arch targets are generated from this
TARGET_ALIAS = $(APP_NAME)-linux-amd64 $(APP_NAME)-linux-arm64 $(APP_NAME)-darwin-amd64 $(APP_NAME)-darwin-arm64
TARGETS = $(addprefix dist/,$(TARGET_ALIAS))
#
# Default make target will run tests
.DEFAULT_GOAL = test

# Build all static Minitor binaries
.PHONY: all
all: $(TARGETS)

# Build all static Linux Minitor binaries
.PHONY: all-linux
all-linux: $(filter dist/$(APP_NAME)-linux-%,$(TARGETS))

# Build nomad-var-dirsync for the current machine
$(APP_NAME): $(GOFILES)
	@echo Version: $(VERSION)
	go build -ldflags '-X "main.version=$(VERSION)"' -o $(APP_NAME)

.PHONY: build
build: $(APP_NAME)

# Run all tests
.PHONY: test
test:
	go test -coverprofile=coverage.out
	go tool cover -func=coverage.out
	@go tool cover -func=coverage.out | awk -v target=80.0% \
		'/^total:/ { print "Total coverage: " $3 " Minimum coverage: " target; if ($3+0.0 >= target+0.0) print "ok"; else { print "fail"; exit 1; } }'

# Installs pre-commit hooks
.PHONY: install-hooks
install-hooks:
	pre-commit install --install-hooks

# Runs pre-commit checks on files
.PHONY: check
check:
	pre-commit run --all-files

.PHONY: clean
clean:
	rm -f ./$(APP_NAME)
	rm -f ./coverage.out
	rm -fr ./dist

## Multi-arch targets
$(TARGETS): $(GOFILES)
	mkdir -p ./dist
	GOOS=$(word 2, $(subst -, ,$(@))) GOARCH=$(word 3, $(subst -, ,$(@))) CGO_ENABLED=0 \
		 go build -ldflags '-X "main.version=$(VERSION)"' -a -installsuffix nocgo \
		 -o $@

.PHONY: $(TARGET_ALIAS)
$(TARGET_ALIAS):
	$(MAKE) $(addprefix dist/,$@)

# Docker targets
.PHONY: docker-build
docker-build:
	docker build -f ./Dockerfile.multi-stage -t $(DOCKER_TAG) .

.PHONY: docker-run
docker-run: docker-build
	docker run --rm -v $(shell pwd)/config.yml:/root/config.yml $(DOCKER_TAG)

# Arch specific docker build targets
.PHONY: docker-build-arm
docker-build-amd64: dist/$(APP_NAME)-linux-amd64
	docker build --platform linux/amd64 . -t DOCKER_TAG-linux-amd64

.PHONY: docker-build-arm64
docker-build-arm64: dist/$(APP_NAME)-linux-arm64
	docker build --build-arg REPO=arm64v8 --build-arg ARCH=arm64 . -t DOCKER_TAG-linux-arm64

# Cross run on host architechture
.PHONY: docker-run-amd64
docker-run-amd64: docker-build-amd64
	docker run --rm -v /var/run/docker.sock:/var/run/docker.sock --name $(DOCKER_TAG)-run DOCKER_TAG-linux-amd64

.PHONY: docker-run-arm
docker-run-arm: docker-build-arm
	docker run --rm -v /var/run/docker.sock:/var/run/docker.sock --name $(DOCKER_TAG)-run DOCKER_TAG-linux-arm

.PHONY: docker-run-arm64
docker-run-arm64: docker-build-arm64
	docker run --rm -v /var/run/docker.sock:/var/run/docker.sock --name $(DOCKER_TAG)-run DOCKER_TAG-linux-arm64
