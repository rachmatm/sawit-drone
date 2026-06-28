# Dockerfile definition for Backend application service.

# From which image we want to build. This is basically our environment.
FROM golang:1.25-alpine

WORKDIR /app

RUN apk add --no-cache git

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

EXPOSE 1323

CMD ["air", "-c", ".air.toml"]