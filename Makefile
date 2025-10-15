# Makefile for go-cfprefs

BASEDIR ?= $(PWD)
SRCDIR ?= $(BASEDIR)

APPNAME ?= spec
APPVER ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: all
all: build test


.PHONY: init
init:
	mkdir -p "$(SRCDIR)/tmp"
	cd $(SRCDIR) && go mod download
	cd $(SRCDIR) && go mod tidy


.PHONY: build
build:


.PHONY: unit-test
unit-test: init
	cd $(SRCDIR) && go test -v -coverprofile=tmp/coverage.out ./...


.PHONY: test
test: unit-test
	cd $(SRCDIR) && go tool cover -func=tmp/coverage.out


.PHONY: static-checks
static-checks: unit-test


.PHONY: preflight
preflight: static-checks


.PHONY: clean
clean:
	cd $(SRCDIR) && go clean


.PHONY: clobber
clobber: clean
	rm -Rf "$(SRCDIR)/dist"
	cd $(SRCDIR) && go clean -modcache
