
# Folder to store temporary files in
BIN_DIR = bin

# Default target name
TARGET := ${BIN_DIR}/rtsp-timestamp

# By default we build for amd64 (x86-64)
ARCH ?= amd64

# By default we build for linux
GOOS ?= linux

# This makefile location to set up the GOPATH
MAKEFILE_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
GOPATH := $(abspath ${MAKEFILE_DIR}/..)

# Targets not associated with files
.PHONY: clean build

clean:
	go clean; \
	sudo rm -rf ${GOPATH}/bin; \
	sudo rm -rf ${GOPATH}/pkg;

build:
	GOARCH=$(ARCH) go get -x -v && \
	GOARCH=$(ARCH) GOOS=$(GOOS) go build -o $(TARGET)