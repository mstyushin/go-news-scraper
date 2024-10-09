SHELL = /usr/bin/bash

APP_NAME := go-news-scraper

$(eval TAGVERSION := $(shell git describe --tags))
$(eval HASHCOMMIT := $(shell git log --pretty=tformat:"%h" -n1 ))
$(eval BRANCHNAME := $(shell git branch --show-current))
ifeq ($(TAGVERSION),undefined)
    # default tag is undefined
    VERSION := $(BRANCHNAME)
else ifeq ($(TAGVERSION),)
    # is empty tag 
    VERSION := $(BRANCHNAME)
else
    VERSION := $(TAGVERSION)
endif
$(eval VERSIONDATE := $(shell git show -s --format=%cI $($VERSION)))

#PG_STARTED=$(shell echo $$((`docker ps --filter "name=db-scraper" --quiet 2> /dev/null | wc -l` + `ps aux|grep -m 1 [p]ostgres:| wc -l`+0)))
PG_STARTED=$(shell echo $$((`docker ps --filter "name=db-scraper" --quiet 2> /dev/null | wc -l` +0)))
pg-run:
ifeq ($(PG_STARTED),0)
	docker run --name db-scraper -d -e POSTGRES_HOST_AUTH_METHOD=trust -p 5432:5432 postgres:15.4
	sleep 2
	psql -d 'postgres://postgres@localhost:5432/postgres?sslmode=disable' -c "CREATE DATABASE news;"
	psql -d 'postgres://postgres@localhost:5432/news?sslmode=disable' -f scripts/schema.sql
endif

PG_TEST_STARTED=$(shell echo $$((`docker ps --filter "name=db-scraper-test" --quiet 2> /dev/null | wc -l` + `ps aux|grep -m 1 [p]ostgres:| wc -l`+0)))
pg-run-test:
ifeq ($(PG_TEST_STARTED),0)
	docker run --name db-scraper-test -d -e POSTGRES_HOST_AUTH_METHOD=trust -p 5432:5432 postgres:15.4
	sleep 2
	psql -d 'postgres://postgres@localhost:5432/postgres?sslmode=disable' -c "CREATE DATABASE news;"
	psql -d 'postgres://postgres@localhost:5432/news?sslmode=disable' -f scripts/schema.sql
	psql -d 'postgres://postgres@localhost:5432/news?sslmode=disable' -f scripts/test_fixtures.sql
endif

build:
	@go mod tidy && go build -ldflags="-X 'go-news-scraper/pkg/config.Version=$(VERSION)' -X 'go-news-scraper/pkg/config.Hash=$(HASHCOMMIT)' -X 'go-news-scraper/pkg/config.VersionDate=$(VERSIONDATE)'" -o bin/$(APP_NAME) go-news-scraper/cmd/server
	@chmod +x bin/$(APP_NAME)

run: build pg-run
	@mkdir -p bin log
	@bin/$(APP_NAME) > log/$(APP_NAME).log 2>&1 & echo "$$!" > /tmp/$(APP_NAME).pid

test: pg-run-test
	@go mod tidy && go test -v ./...

stop:
	-pkill -f $(APP_NAME)

clean: stop
	@rm -f bin/*
	@rm -f log/*
	@rm -f /tmp/$(APP_NAME).pid
	docker stop db-scraper || docker stop db-scraper-test || true
	docker container prune -f && docker volume prune -f
