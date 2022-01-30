FLAGS?=-v
CMD?=hw1

default: test lint

.SILENT:

build:
	go build -o ./bin/$(CMD) $(FLAGS) ./cmd/$(CMD)

build-all:
	$(foreach dir,$(wildcard cmd/*), go build $(FLAGS) ./$(dir);)

docker:
	docker build -f ./internal/$(CMD)/Dockerfile -t arch_course/$(CMD) .

docker_local: docker
	minikube image load arch_course/$(CMD):latest

test:
	go test $(FLAGS) ./...

lint:
	golangci-lint run -v ./...

tidy:
	go mod tidy

.PHONY: build build-all test lint tidy
