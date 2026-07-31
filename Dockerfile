FROM golang:1.26-alpine AS build

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o /cli-auth ./cmd

FROM alpine:latest

RUN apk add --no-cache libgcc

WORKDIR /app

RUN mkdir -p /app/data

COPY --from=build /cli-auth /app/cli-auth

ENTRYPOINT ["/app/cli-auth"]
