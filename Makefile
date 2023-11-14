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

build: gvt 
	go build  -C cmd/gophermart .

run:
	go run ./cmd/gophermart

rundb:
	docker compose up -d

	
racetest:
	go test -v -race ./...

autotest: build
	bin/gophermarttest \
		-test.v -test.run=^TestGophermart$$ \
		-gophermart-binary-path=cmd/gophermart/gophermart \
		-gophermart-host=localhost \
		-gophermart-port=8080 \
		-gophermart-database-uri="postgresql://postgres:postgres@postgres/praktikum?sslmode=disable" \
		-accrual-binary-path=cmd/accrual/accrual_linux_amd64 \
		-accrual-host=localhost \
		-accrual-port=8081 \
		-accrual-database-uri="postgresql://postgres:postgres@postgres/praktikum?sslmode=disable"
