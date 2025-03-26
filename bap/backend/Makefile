# Makefile for building bpp
# Reference Guide - https://www.gnu.org/software/make/manual/make.html

#
# Internal variables or constants.
# NOTE - These will be executed when any make target is invoked.
#

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)
SWAG ?= $(LOCALBIN)/swag

.PHONY: swag
swag: $(SWAG) ## Download swag locally if necessary.
$(SWAG): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/swaggo/swag/cmd/swag@v1.16.1

.PHONY: server-api-doc
server-api-doc: swag
	$(SWAG) init --parseDependency --parseInternal -g main.go
