FROM golang:1.24-alpine AS builder

ENV GIN_MODE=release

WORKDIR /src

RUN apk update && apk upgrade && apk add --no-cache make

COPY go.mod go.sum .

RUN go mod download

COPY . .

RUN make build

FROM alpine:latest

WORKDIR /app

ENV PORT=8000
ENV JWK_URL=""

COPY --from=builder /src/out /app

EXPOSE $PORT
VOLUME /data/conf

CMD ./server -port ${PORT} -jwk-url ${JWK_URL}
