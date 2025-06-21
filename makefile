APP_NAME=k8s-controller-tutorial
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_FLAGS ?= -v -o $(APP_NAME) -ldflags "-X=github.com/JRaver/k8s-controller-tutorial/cmd.appVersion=$(VERSION)"

.PHONY: all build tests run clean

all: build

build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(BUILD_FLAGS) cmd/main.go

test:
	go test ./...

run: 
	go run main.go

docker-build:
	docker build --build-arg VERSION=$(VERSION) -t $(APP_NAME):$(VERSION) .

clean:
	rm -f $(APP_NAME)
