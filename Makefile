include .env

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v
	sqlc generate

## audit: run quality control checks
.PHONY: audit
audit: tidy
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go test -race -vet=off ./...
	go mod verify

## audit: upgrade modfile versions
.PHONY: upgrade
upgrade:
	go get -u ./...

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build: build the application
.PHONY: build
build: tidy
	go mod verify
	ko build ./cmd/ephr

# ==================================================================================== #
# DB
# ==================================================================================== #

## up: start docker compose and apply db migrations
.PHONY: up
up:
	docker compose up -d; sleep 1
	dbmate up

## down: build the cmd/api application
.PHONY: down
down:
	docker compose down

## psql: connect to the database using psql
.PHONY: psql
psql:
	psql ${DATABASE_URL}
