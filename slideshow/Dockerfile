ARG BUILD_FROM
# hadolint ignore=DL3006

FROM golang:1.23.4-bookworm as builder

# Set build arguments with default value
ARG ARCH=amd64

# Set default architecture
ENV GOARCH=amd64

# Translate matrix.arch to goarch
RUN case ${ARCH} in \
        amd64) export GOARCH=amd64 ;; \
        x86_64) export GOARCH=amd64 ;; \
        armv7) export GOARCH=arm ;; \
        armv6) export GOARCH=arm ;; \
        arm64) export GOARCH=arm64 ;; \
        aarch64) export GOARCH=arm64 ;; \
        386) export GOARCH=386 ;; \
        i386) export GOARCH=386 ;; \
        riscv64) export GOARCH=riscv64 ;; \
        ppc64le) export GOARCH=ppc64le ;; \
        s390x) export GOARCH=s390x ;; \
        mips) export GOARCH=mips ;; \
        mipsle) export GOARCH=mipsle ;; \
        mips64) export GOARCH=mips64 ;; \
        mips64le) export GOARCH=mips64le ;; \
        *) echo "Unsupported architecture ${MATRIX_ARCH}" >&2; exit 1 ;; \
    esac && \
    echo "Building for GOARCH=${GOARCH}"

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
RUN GOOS=linux GOARCH=$GOARCH go build -a -tags netgo -ldflags "-w -s" -o app .

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