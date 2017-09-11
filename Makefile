# Copyright (c) 2015 Geofrey Ernest. All rights reserved.
# Use of this source code is governed by a MIT
# license that can be found in the LICENSE file.

.PHONY: all clean nuke migration-test bindata
DEFAULT_POSTGRES_CONN	:=postgres://postgres@localhost/zedlist_test?sslmode=disable
STATIC_EMBED		:=bindata/static/static.go
TMPl_EMBED		:=bindata/template/templates.go
COMPONENTS		:=./middlewares/... ./modules/... ./routes/...

ifeq "$(origin CONFIG_DBCONN)" "undefined"
CONFIG_DBCONN=$(DEFAULT_POSTGRES_CONN)
endif
ifeq "$(origin APP_VER)" "undefined"
VERSION=0.1.0
endif
all: lint bindata test build

clean:
	@rm -r bin
	
nuke:
	go clean -i

build:
	@go build -o bin/zedlist ./cmd/zedlist

test:migration-test
	@CONFIG_DBCONN=$(CONFIG_DBCONN) go test $(COMPONENTS)

bindata:
	@go-bindata  -pkg=static -o=$(STATIC_EMBED) static/...
	@go-bindata -pkg=template -o=$(TMPl_EMBED) -prefix=templates/ templates/...

cover:test
	@CONFIG_DBCONN=$(CONFIG_DBCONN) bash ./scripts/coverage.sh

watch:
	@sass --watch assets/sass:static/css
	
lint:
	@-go vet ./...
	@-golint ./... |grep -v bindata/template/* |grep -v bindata/static/*

migration-test:
	@CONFIG_DBCONN=$(CONFIG_DBCONN) go test ./migration

deps:
	@go get -u -v github.com/golang/lint/golint
	@go get -u -v github.com/jteeuwen/go-bindata/...
	@go get github.com/mitchellh/gox

dist:
	@gox -output="dist/{{.Dir}}v$(VERSION)_{{.OS}}_{{.Arch}}" ./cmd/zedlist

server:build
	./bin/zedlist server
	
migrate:
	./bin/zedlist migrate -d