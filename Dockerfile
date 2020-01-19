ARG GO_VERSION=1.13.6
ARG DEBIAN_VERSION=buster
ARG BASE_IMAGE=gcr.io/distroless/static:nonroot
############################
# STEP 1 build executable binary
############################
# golang debian buster 1.13.6 linux/amd64
# FROM golang:${GO_VERSION}-${DEBIAN_VERSION}
# https://github.com/docker-library/golang/blob/master/1.13/buster/Dockerfile
FROM golang:${GO_VERSION}-${DEBIAN_VERSION} as builder

# Ensure ca-certficates are up to date
RUN update-ca-certificates

# Moving outside of $GOPATH forces modules on without having to set ENVs
WORKDIR /src

ENV GOFLAGS="-mod=readonly"
# ENV GO111MODULE=on

# Moving outside of $GOPATH forces modules on without having to set ENVs
WORKDIR /src

# Add go.mod and go.sum first to maximize caching
COPY ./go.mod ./go.sum ./

RUN go mod download
RUN go mod verify

COPY . .

# Build the static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -ldflags='-w -s -extldflags "-static"' -a \
      -o /go/bin/main ./main.go

############################
# STEP 2 build a small image
############################
# FROM gcr.io/distroless/static:nonroot
# user:group is nobody:nobody, uid:gid = 65534:65534
FROM ${BASE_IMAGE}

# Copy our static executable
COPY --from=builder /go/bin/main /go/bin/main

ARG GIT_COMMIT=unknown
LABEL git-commit=$GIT_COMMIT
ARG GIT_BRANCH=unknown
LABEL git-branch=$GIT_BRANCH
ARG BUILD_TIME=unknown
LABEL build_time=$BUILD_TIME

# Run completely unprivileged.
USER nobody:nobody
# Run the hello binary.
ENTRYPOINT ["/go/bin/main"]
