FROM golang:1.21-alpine3.18 AS builder

LABEL maintainer="cyclimse"

WORKDIR /code

# Using git to add the version to the binary
# See: https://tip.golang.org/doc/go1.18#go-version
RUN apk add --no-cache \
    git 

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

# TODO: fix sqlc unknown field Env in struct literal of type wasm.Runner 
# RUN --mount=type=cache,target=/go/pkg/mod \
#     go generate ./...

RUN --mount=type=cache,target=/go/pkg/mod \
    go build -o bin/blahaj ./cmd/blahaj

FROM alpine:3.18 as blahaj

WORKDIR /app

COPY --from=builder /code/bin/ .

FROM blahaj as api

EXPOSE 8080

ENTRYPOINT ["/app/blahaj", "api"]

FROM blahaj as crawler

EXPOSE 8081

ENTRYPOINT ["/app/blahaj", "crawl"]
