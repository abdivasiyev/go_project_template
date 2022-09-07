include .env
CURRENT_DIR=$(shell pwd)

.PHONY: run

run: migrate-up swag
	go run cmd/logistics/*

.PHONY: build
build:
	go build -o backend cmd/logistics/*

.PHONY: test
test:
	go test -v ./...

.PHONY: cover-test
cover-test:
	go test -race -covermode=atomic -coverprofile=coverage.out -v ./...

.PHONY: migrate-up
migrate-up:
	migrate -path=${CURRENT_DIR}/migrations -database=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable up

.PHONY: migrate-down
migrate-down:
	migrate -path=${CURRENT_DIR}/migrations -database=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable down

.PHONY: create-new-migration
create-new-migration: # make create-new-migration name=file_name
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: swag
swag:
	swag init -g v1.go --dir internal/handler/v1,pkg/response,internal/models --output docs/v1 --ot json

.PHONY: deploy
deploy:
	git checkout master && git pull origin master && git merge --no-ff dev && git push origin master && git checkout dev

.DEFAULT_GOAL:=run