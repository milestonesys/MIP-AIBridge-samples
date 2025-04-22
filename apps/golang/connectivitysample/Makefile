# Folders to store temporary files in
OBJ_DIR = obj
BIN_DIR = bin

# Default target name
TARGET := $(BIN_DIR)/connectivitysample

# By default we build for amd64 (x86-64)
ARCH ?= amd64

# By default we build for linux
GOOS ?= linux

build: dependencies
	@GOARCH=$(ARCH) go get -x -v && \
	 GOARCH=$(ARCH) GOOS=$(GOOS) go build -o $(TARGET)

dependencies:
	@mkdir -p $(BIN_DIR) && \
	 cp start.sh $(BIN_DIR)/start.sh
	 cp opensource.html $(BIN_DIR)/opensource.html

clean:
	@go clean && \
	 rm -fr ${BIN_DIR} && \
	 rm -fr ${OBJ_DIR}
