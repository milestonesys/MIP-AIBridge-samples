#!/bin/bash
set -e -o pipefail

# Locate files in the system
APP_ROOT_DIR=$(realpath `dirname "$0"`)
EXECUTABLE_DIR="${APP_ROOT_DIR}/bin"
EXECUTABLE_FILE="${EXECUTABLE_DIR}/connectivitysample"

# Check the binary was built
if [ ! -f "${EXECUTABLE_FILE}" ]; then
	echo 'Make sure to build the binary first.'
	echo 'Navigate to the build folder and run the `build-binary.sh` script.'
	exit 1
fi

# Read the .env file to obtain the App configuration
source "${APP_ROOT_DIR}/.env"

# Declare needed environment variables
export TAG="1.0.0"
export EXTERNAL_HOSTNAME=${EXTERNAL_HOSTNAME}
export APP_WEBSERVER_PORT=${APP_WEBSERVER_PORT}
export APP_URL_PATH=${APP_URL_PATH}
export TLS_SCHEME=${TLS_SCHEME:-http}

# Run binary providing relevant parameters
${EXECUTABLE_FILE} \
	-aib-webservice-location localhost:4000 \
	-app-registration-file-path "${APP_ROOT_DIR}/config/register.graphql" \
	-app-url-path ${APP_URL_PATH} \
	-app-webserver-port ${APP_WEBSERVER_PORT} \
	-enforce-oauth=true \
	-snapshot-max-height 600 \
	-snapshot-max-width 600 \
	-tls-certificate-file "${APP_ROOT_DIR}/certs/tls-server/server.crt" \
	-tls-key-file "${APP_ROOT_DIR}/certs/tls-server/server.key" \
	-tls-enabled=${TLS_ENABLED:-false}
