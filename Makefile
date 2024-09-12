# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=sshselect

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

install: build
	mkdir -p $(HOME)/bin
	cp $(BINARY_NAME) $(HOME)/bin/

.PHONY: all build clean install