ARG BUILDPLATFORM="linux/amd64"

FROM --platform=${BUILDPLATFORM} golang:1.22 AS builder
ARG TARGETARCH
WORKDIR /wspc/src/connectivitysample

# Add source code of connectivitysample sample
COPY . .

# Build artifacts
RUN export GOPATH=/wspc ARCH="${TARGETARCH}" && \
	make -s clean && \
	make -s build

# Build image with ubuntu 22.04 as the base
FROM ubuntu:22.04

# Use bash as the default shell and instruct it to exit immediately if a command exits with a non-zero status
SHELL ["/bin/bash", "-e", "-c"]
WORKDIR /root

# The app version matches the docker image tag.
ENV TAG="1.0.0"

# Update packages and install dependencies
RUN apt update && apt upgrade -y && \
	apt install -y ca-certificates && \
	apt clean && \ 
	rm -rf /var/cache/apt/archives /var/lib/apt/lists/*

# Copy in files that we have built
COPY --from=builder /wspc/src/connectivitysample/bin bin
RUN chmod +x bin/start.sh

# Container runtime arguments
WORKDIR /root/bin
ENTRYPOINT [ "./start.sh" ]
CMD [ ]
