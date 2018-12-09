PWD := $(shell pwd)
SYSTEM := $(shell uname -s)-amd64
VERSION := v$(shell cat version)
OPTS := -ldflags "-X main.version=$(VERSION)"
NAME := mknote

default: all

all: clean bin tar

clean: 
	@rm -rf _output/*

bin:
	@echo "build $(NAME) $(VERSION)"
	@cp -r build _output/$(NAME)-$(VERSION)-linux-amd64 
	@GO111MODULE=off GOOS=linux go build $(OPTS) -o _output/$(NAME)-$(VERSION)-linux-amd64/bin/$(NAME) ./cmd/$(NAME)
	@echo "successful binary to _output/$(NAME)-$(VERSION)-linux-amd64/bin/$(NAME)"

tar:
	@cd _output && tar -zcf $(NAME)-$(VERSION)-linux-amd64.tar.gz $(NAME)-$(VERSION)-linux-amd64
	@echo "successful tarball to _output/$(NAME)-$(VERSION)-linux-amd64.tar.gz"
	
