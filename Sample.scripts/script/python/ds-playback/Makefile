# Folders to store temporary files in
BIN_DIR = bin
GEN_DIR = generated
DATA_DIR = data
VENV_DIR ?= playback_venv
SYSTEM_IP ?= 127.0.0.1

# Protobuffer compiler output
GRPC_ARTIFACTS =
GRPC_ARTIFACTS += ${GEN_DIR}/protos/ds/v1/ds_pb2_grpc.py
GRPC_ARTIFACTS += ${GEN_DIR}/protos/ds/v1/ds_pb2.py
GRPC_ARTIFACTS += ${GEN_DIR}/protos/ds/v1/ds_pb2.pyi

# This makefile location
MAKEFILE_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))


.PHONY: setup clean create_environment clean_environment get_proto generate_grpc_files clean_grpc_files

setup: create_environment get_proto generate_grpc_files

clean: clean_environment clean_grpc_files
	@rm -rf ${BIN_DIR}; \
	 rm -rf ${DATA_DIR}

# Setup env
create_environment:
	@if [ ! -d ${VENV_DIR} ]; then \
		python3 -m venv ${VENV_DIR}; \
	fi; \
	. ${VENV_DIR}/bin/activate && pip3 install -r requirements.txt

clean_environment:
	@rm -rf ${VENV_DIR}

# Proto and grpc
get_proto:
	@rm -rf ${BIN_DIR}/protos/ds/v1/ && \
	 mkdir -p ${BIN_DIR}/protos/ds/v1/ && \
	 wget http://${SYSTEM_IP}:4000/resources/ds.proto -O ${BIN_DIR}/protos/ds/v1/ds.proto

generate_grpc_files: ${GRPC_ARTIFACTS}

clean_grpc_files:
	@rm -f ${GRPC_ARTIFACTS} && \
	 rm -rf ${GEN_DIR}

${GEN_DIR}/protos/ds/v1/ds_pb2_grpc.py ${GEN_DIR}/protos/ds/v1/ds_pb2.py ${GEN_DIR}/protos/ds/v1/ds_pb2.pyi:
	@mkdir -p ${GEN_DIR}/protos/ds/v1 && \
	. ${VENV_DIR}/bin/activate && \
	python3 -m grpc_tools.protoc \
		--python_out=${GEN_DIR}/protos/ds/v1 \
		--grpc_python_out=${GEN_DIR}/protos/ds/v1 \
		--pyi_out=${GEN_DIR}/protos/ds/v1 \
		-I ${BIN_DIR}/protos/ds/v1/ ${BIN_DIR}/protos/ds/v1/ds.proto && \
	touch ${GEN_DIR}/protos/ds/v1/__init__.py
