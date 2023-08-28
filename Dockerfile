# This Docker allows running a Caerus server instance.
#
# How to build the image:
# > docker build --tag desmoslabs/caerus .
#
# How to run the image:
# > docker run desmoslabs/caerus

FROM golang:1.20-alpine
ARG arch=x86_64

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3 ca-certificates build-base
RUN set -eux; apk add --no-cache $PACKAGES;

# Set working directory for the build
WORKDIR /code

# Add sources files
COPY . /code/

# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.3.0/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep b1610f9c8ad8bdebf5b8f819f71d238466f83521c74a2deb799078932e862722

ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.3.0/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep b4aad4480f9b4c46635b4943beedbb72c929eab1d1b9467fe3b43e6dbf617e32

# Copy the library you want to the final location that will be found by the linker flag `-lwasmvm_muslc`
RUN cp /lib/libwasmvm_muslc.${arch}.a /usr/local/lib/libwasmvm_muslc.a

# Set the environment variables
ENV GIN_MODE=release

# Build the executable
RUN BUILD_TAGS=muslc GOOS=linux GOARCH=amd64 LINK_STATICALLY=true make build

# Move the executable inside the bin folder to make it runnable without specifying the full path
RUN cp /code/build/caerus /usr/bin/caerus

# Set the entrypoint, so that the user can set the config using the CMD
ENTRYPOINT ["caerus"]