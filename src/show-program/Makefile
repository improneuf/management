# Makefile for Go project

# Variables
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
BINARY_NAME=app

# Source files
SOURCES=main.go price.go showtype.go showlanguage.go show.go drive.go photo.go

# Build the project
build:
	$(GOBUILD) -o $(BINARY_NAME) $(SOURCES)

# Run the project
run:
	mkdir -p output/screenshots
	cp -R icons output/screenshots/icons
	$(GORUN) $(SOURCES)

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# PHONY targets
.PHONY: build run clean
