FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod go.sum *.go ./
RUN go mod download && \
    go build -o /rot-notifier

FROM alpine AS release

WORKDIR /

COPY --from=build /rot-notifier /rot-notifier

ENTRYPOINT ["/rot-notifier"]