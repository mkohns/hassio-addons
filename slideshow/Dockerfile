ARG BUILD_FROM
# hadolint ignore=DL3006

FROM golang:1.23.4-bookworm as builder

# Argument to get the target platform from Docker
ARG TARGETPLATFORM

# Set default environment variables
ENV GOOS=linux \
    GOARCH=amd64 \
    GOARM=""

# Set GOOS, GOARCH, and GOARM based on TARGETPLATFORM
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
    if [ -n "$GOARM" ]; then echo "GOARM=$GOARM"; fi


# Create appuser.
ENV USER=appuser
ENV UID=10001 

# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

# copy the code
COPY . /go/src/app/

# care for the frontend
RUN apt update && apt install -y nodejs npm
RUN npm install -g yarn
WORKDIR /go/src/app/frontend/slideshow
RUN ls -alh
RUN yarn install
RUN yarn build

# care for the server
WORKDIR /go/src/app
RUN go mod download
RUN go mod verify
RUN go build -a -tags netgo -ldflags "-w -s" -o app .

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /go/src/app/app /bin/app
COPY --from=builder /go/src/app/frontend/slideshow/dist /dist

# Use an unprivileged user.
USER appuser:appuser

ENTRYPOINT [ "/bin/app" ]