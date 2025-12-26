FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy
COPY . .
RUN go mod download
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/file-validator ./cmd/app

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /bin/file-validator /app/file-validator
COPY config.yml /app/config.yml
ENV CONFIG_PATH=/app/config.yml
EXPOSE 8065
CMD ["/app/file-validator"]
