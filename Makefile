PWD := $(shell pwd)
SYSTEM := $(shell uname -s)-amd64
VERSION := v$(shell cat version)
PKGNAME := mknote-$(VERSION)-linux-amd64
OPTS := -ldflags "-X main.version=$(VERSION)"

default: clean bin tar

clean:
	@rm -rf _output

bin:
	@echo "building $(PKGNAME)"
	@mkdir -p _output
	@cp -r build _output/$(PKGNAME)
	@GOOS=linux go build $(OPTS) -o _output/$(PKGNAME)/bin/mknote ./cmd/mknote
	@echo "successful binary to _output/$(PKGNAME)/bin/mknote"

tar:
	@cd _output && tar -zcf $(PKGNAME).tar.gz $(PKGNAME)
	@echo "successful tarball to _output/$(PKGNAME).tar.gz"
