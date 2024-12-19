ARG BUILD_FROM
# hadolint ignore=DL3006

FROM golang:1.23.4-bookworm as builder

# Argument to get the target platform from Docker
ARG TARGETPLATFORM="linux/arm64"

# care for the frontend
RUN apt update && apt install -y nodejs npm
RUN npm install -g yarn

