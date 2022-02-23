FLAGS?=-v
CMD?=hw5
SERVICE?=

services := $(notdir $(shell find ./internal/$(CMD)/ -mindepth 1 -maxdepth 1 -type d))

default: test lint

.SILENT:

build:
	if [ "$(SERVICE)" = "" ]; then \
		go build -o ./bin/$(CMD) $(FLAGS) ./cmd/$(CMD) ;\
	else \
  		go build -o ./bin/$(CMD)/$(SERVICE) $(FLAGS) ./cmd/$(CMD)/$(SERVICE) ;\
	fi

build_all:
	$(foreach dir,$(wildcard cmd/*), go build $(FLAGS) ./$(dir);)

docker:
	if [ "$(SERVICE)" = "" ]; then \
		docker build -f ./internal/$(CMD)/Dockerfile -t arch_course/$(CMD) . ;\
	else \
	  	docker build -f ./internal/$(CMD)/$(SERVICE)/Dockerfile -t arch_course/$(CMD)/$(SERVICE) . ;\
	fi

docker_all:
	for service in $(services) ; do \
  		docker build -f ./internal/$(CMD)/$$service/Dockerfile -t arch_course/$(CMD)/$$service . ;\
  	done

docker_local: docker
	minikube image load arch_course/$(CMD):latest

test:
	go test $(FLAGS) ./...

lint:
	golangci-lint run -v ./...

tidy:
	go mod tidy

.PHONY: build build_all docker docker_all test lint tidy
