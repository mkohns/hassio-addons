ARG BUILD_FROM
ARG COMMIT_HASH
FROM $BUILD_FROM as builder

# Argument to get the target platform from Docker
ARG TARGETPLATFORM="linux/arm64"

# Set default environment variables
ENV GOOS=linux \
    GOARCH=amd64 \
    GOARM=""

RUN echo "Show all envrironment variables" && \
    env

# copy the code
COPY . /go/src/app/

# first care for the server 
WORKDIR /go/src/app
RUN go mod download && \
    go mod verify && \
    go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
RUN ls -al /go/bin && \
    if [ -d "/go/bin/linux_amd64" ]; then export PATH="$PATH:/go/bin/linux_amd64"; else echo "/go/bin/linux_amd64 does not exist"; fi && \
    oapi-codegen -version && \
    go generate
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
    GITHASH=$COMMIT_HASH && \
    echo "VERSION=$VERSION" && \
    echo "GITHASH=$GITHASH" && \
    if [ -n "$GOARM" ]; then echo "GOARM=$GOARM"; fi && \
    go build -a -tags netgo -ldflags "-w -s -X main.addon_version=$VERSION -X main.addon_githash=$GITHASH" -o app .

# second care for the frontend
WORKDIR /go/src/app/frontend/slideshow
RUN yarn install && yarn build


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