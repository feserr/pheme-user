# Builder
FROM golang:1.19-buster AS builder

ARG VERSION=dev

WORKDIR /go/src/app
COPY . .
COPY .env .
RUN go get -d -v ./... \
  && go install -v ./... \
  && go build -o main -ldflags=-X=main.version=${VERSION} main.go

# Runner
FROM debian:buster-slim

COPY --from=builder /go/src/app/main /go/bin/main
ENV PATH="/go/bin:${PATH}"
ARG SERVER_PORT=8001
EXPOSE ${SERVER_PORT}

CMD ["main"]