# syntax=docker/dockerfile:1
FROM golang:1.20 AS builder
WORKDIR /go/src/

COPY go.* .
COPY cmd ./cmd/

RUN CGO_ENABLED=0 GOOS=linux go build -a -o ./bin/sshdocker ./cmd/sshdocker

FROM alpine:3.18.2
RUN apk --no-cache add \
  ca-certificates \
  docker

WORKDIR /
COPY --from=builder /go/src/bin/sshdocker .

ENTRYPOINT ["./sshdocker"]
