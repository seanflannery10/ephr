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
	sqlc compile
	sqlc generate
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit: tidy
	golangci-lint run --fix
	go test -race -vet=off ./...
	go mod verify

## upgrade: upgrade modfile versions
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
# DEV
# ==================================================================================== #

## up: start the application stack
.PHONY: up
up:
	tilt up

## down: stop the application stack
.PHONY: down
down:
	tilt down --delete-namespaces

# ==================================================================================== #
# DB
# ==================================================================================== #

## migrations: apply db migrations
.PHONY: migrations
migrations:
	dbmate -d "./sql/migrations" --url ${DB_URL} up

## psql: connect to the database using psql
.PHONY: psql
psql:
	psql ${DB_URL}
