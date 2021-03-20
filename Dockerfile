FROM golang:1.16.2-alpine3.13 as builder


WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
RUN go build \
    -o /go/bin/main \
    -ldflags '-s -w'

ENV PORT=${PORT}

FROM scratch as runner

COPY --from=builder /go/bin/main /app/main

ENTRYPOINT ["/app/main"]
