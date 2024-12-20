ARG BUILD_FROM
FROM $BUILD_FROM as builder

# Argument to get the target platform from Docker
ARG TARGETPLATFORM="linux/arm64"

# Set default environment variables
ENV GOOS=linux \
    GOARCH=amd64 \
    GOARM=""

# copy the code
COPY . /go/src/app/

WORKDIR /go/src/app/frontend/slideshow
RUN yarn install && yarn build

# care for the server
WORKDIR /go/src/app
RUN go mod download && go mod verify
RUN case "$TARGETPLATFORM" in \
      "linux/amd64") \
        export GOOS=linux GOARCH=amd64 ;; \
      "linux/arm/v6") \
        export GOOS=linux GOARCH=arm GOARM=6 ;; \
      "linux/arm/v7") \
        export GOOS=linux GOARCH=arm GOARM=7 ;; \
      "linux/386") \
        export GOOS=linux GOARCH=386 ;; \
      "linux/arm64") \
        export GOOS=linux GOARCH=arm64 ;; \
      *) \
        echo "Unsupported platform: $TARGETPLATFORM" && exit 1 ;; \
    esac && \
    echo "Building for TARGETPLATFORM=$TARGETPLATFORM" && \
    echo "GOOS=$GOOS" && \
    echo "GOARCH=$GOARCH" && \
    VERSION=$(jq -r '.version' config.json) && \
    GITHASH=$(git rev-parse --short HEAD) && \
    if [ -n "$GOARM" ]; then echo "GOARM=$GOARM"; fi && \
    go build -a -tags netgo -ldflags "-w -s -X main.addon_version=$VERSION -X main.addon_githash=$GITHASH" -o app .

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /go/src/app/app /bin/app
COPY --from=builder /go/src/app/frontend/slideshow/dist /dist

# Use an rivileged user because options.json is owned by root
USER root:root

ENTRYPOINT [ "/bin/app" ]