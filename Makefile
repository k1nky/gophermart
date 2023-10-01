SHELL:=/bin/bash
STATICCHECK=$(shell which staticcheck)

.DEFAULT_GOAL := build

test:
	go test -cover ./...

vet:
	go vet ./...
	$(STATICCHECK) ./...

generate:
	go generate ./...

gvt: generate vet test

cover:
	go test -cover ./... -coverprofile cover.out
	go tool cover -html cover.out -o cover.html

build: gvt buildagent buildserver

buildserver:
	go build  -C cmd/server .

buildagent:
	go build -C cmd/agent .

runserver:
	go run ./cmd/server

runagent:
	go run ./cmd/agent

rundb:
	docker compose up -d

racetest:
	go test -v -race ./...

autotest: autotest1

autotest1: buildserver
	metricstest -test.v -test.run=^TestIteration1$$ -binary-path=cmd/server/server
