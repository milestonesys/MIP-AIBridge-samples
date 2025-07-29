
# Folders to store temporary files in
BIN_DIR = bin
GEN_DIR = generated
DATA_DIR = data
SYSTEM_IP ?= 127.0.0.1

# Default target name
TARGET := ${BIN_DIR}/ds-playback

# By default we build for amd64 (x86-64)
ARCH ?= amd64

# By default we build for linux
GOOS ?= linux

# Protobuffer compiler output
GRPC_ARTIFACTS =
GRPC_ARTIFACTS += ${GEN_DIR}/protos/ds/v1/ds.pb.go
GRPC_ARTIFACTS += ${GEN_DIR}/protos/ds/v1/ds_grpc.pb.go

# This makefile location to set up the GOPATH
MAKEFILE_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
GOPATH := $(abspath ${MAKEFILE_DIR}/..)

# Targets not associated with files
.PHONY: setup clean requirements get_proto generate_grpc_files clean_grpc_files build

setup: get_proto requirements generate_grpc_files

clean: clean_grpc_files
	go clean; \
	sudo rm -rf ${GOPATH}/bin; \
	sudo rm -rf ${GOPATH}/pkg; \
	rm -rf ${BIN_DIR}; \
	rm -rf ${DATA_DIR}

requirements:
	# 1- Install protoc reading from https://grpc.io/docs/protoc-installation/
	sudo apt install -y protobuf-compiler
	# 2- Install needed plugin https://grpc.io/docs/languages/go/quickstart/
	# The binaries will be installed at ${GOPATH}/bin
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest; \
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

get_proto:
	rm -rf ${BIN_DIR}/protos/ds/v1/ && \
	mkdir -p ${BIN_DIR}/protos/ds/v1/ && \
	wget http://${SYSTEM_IP}:4000/resources/ds.proto -O ${BIN_DIR}/protos/ds/v1/ds.proto

generate_grpc_files: ${GRPC_ARTIFACTS}

clean_grpc_files:
	rm -f ${GRPC_ARTIFACTS} && \
	rm -rf ${GEN_DIR}

${GEN_DIR}/protos/ds/v1/ds.pb.go ${GEN_DIR}/protos/ds/v1/ds_grpc.pb.go:
	PATH="${PATH}:${GOPATH}/bin" protoc -I${GOPATH}/src/ \
		--go_out=. --go_opt=M${BIN_DIR}/protos/ds/v1/ds.proto=${GEN_DIR}/protos/ds/v1 \
		--go-grpc_out=. --go-grpc_opt=M${BIN_DIR}/protos/ds/v1/ds.proto=${GEN_DIR}/protos/ds/v1 \
		--plugin=${GOPATH}/bin/protoc-gen-go ${BIN_DIR}/protos/ds/v1/ds.proto

build:
	GOARCH=$(ARCH) go get -x -v && \
	GOARCH=$(ARCH) GOOS=$(GOOS) go build -o $(TARGET)
