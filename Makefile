SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later)
endif
.RECIPEPREFIX = >

# go builds are fast enough that we can just build on demand instead of trying to do any fancy
# change detection
build: clean prcr
.PHONY: build

prcr:
> go build ./cmd/prcr

clean:
> rm -f ./prcr
.PHONY: clean

test:
> go test ./...
.PHONY: test
