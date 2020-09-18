VERSION := $(shell cat .version )
GO = go

PROGNAME = binsplit
PROGNAME_VERSION = $(PROGNAME)-$(VERSION)
TARGZ_FILENAME = $(PROGNAME)-$(VERSION).tar.gz
TARGZ_CONTENTS = binsplit README.md Makefile .version

PREFIX = /tmp
PWD = $(shell pwd)

export PROGROOT=$(PWD)/$(PROGNAME_VERSION)

.PHONY: all version build clean install test


linux: $(PROGNAME)
	env GOOS=linux $(GO) build -ldflags="-X 'main.BuildVersion=$(VERSION)'" -v .

osx:
	env GOOS=darwin $(GO) build -ldflags="-X 'main.BuildVersion=$(VERSION)'" -v .

.PHONY: clean
clean:
	rm -vf "$(PROGNAME)"
