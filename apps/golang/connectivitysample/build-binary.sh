#!/bin/bash
set -eo pipefail

# Build the executable
make -s clean
make -s build
