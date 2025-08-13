FROM golang:1.24.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY configs ./configs
COPY migrations ./migrations
COPY docs ./docs

RUN go build -o myservice cmd/service/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/myservice .

CMD ["./myservice"]
