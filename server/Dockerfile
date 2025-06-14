FROM golang:1.24 AS builder

ENV GIN_MODE=release

WORKDIR /src

COPY go.mod go.sum .

RUN go mod download

COPY . .

RUN make build

FROM alpine:latest

WORKDIR /app

ENV PORT="8000"

COPY --from=builder /src/out /app

EXPOSE $PORT
VOLUME /data/conf

CMD ["./server", "-port", "$PORT"]
