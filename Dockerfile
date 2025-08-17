FROM golang:1.23-alpine AS build

WORKDIR /app

RUN apk add --no-cache make nodejs npm curl
RUN npm install -g pnpm

COPY go.mod go.sum ./
RUN go mod download

COPY Makefile ./
COPY ./web ./web

RUN make build-web

COPY . .

RUN make build-only-aio

# ---- FINAL -----
FROM alpine:latest

WORKDIR /app
RUN apk add --no-cache docker docker-cli docker-cli-compose

COPY --from=build /app/build/aio /app/aio

EXPOSE 4000

CMD ["/app/aio"]