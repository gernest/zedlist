# Copyright (c) 2015 Geofrey Ernest. All rights reserved.
# Use of this source code is governed by a MIT
# license that can be found in the LICENSE file.

.PHONY: all clean nuke migration-test bindata
POSTGRES_CONN	:=postgres://postgres:postgres@localhost/zedlist_test?sslmode=disable
STATIC_EMBED		:=bindata/static/static.go
TMPl_EMBED		:=bindata/template/templates.go
COMPONENTS		:=./middlewares/... ./modules/... ./routes/...

all: dev

clean:
	go clean
	rm -f  *.out
	
nuke:
	go clean -i


test:migration-test
	@CONFIG_DBCONN=$(POSTGRES_CONN) go test $(COMPONENTS)

bindata:
	@go-bindata -debug=true -pkg=static -o=$(STATIC_EMBED) static/...
	@go-bindata -debug=true -pkg=template -o=$(TMPl_EMBED) -prefix=templates/ templates/...

dev: bindata test
	@go build

cover:
	@CONFIG_POSTGRES_CONN=$(POSTGRES_CONN) bash ./scripts/coverage.sh

watch:
	@sass --watch assets/sass:static/css
	
lint:
	@go vet ./...
	@golint ./... |grep -v modules/template/* |grep -v modules/static/*

migration-test:
	@CONFIG_DBCONN=$(POSTGRES_CONN) go test ./migration

bindata-prod:
	@go-bindata  -pkg=static -o=$(STATIC_EMBED) static/...
	@go-bindata  -pkg=template -o=$(TMPl_EMBED) -prefix=templates/ templates/...