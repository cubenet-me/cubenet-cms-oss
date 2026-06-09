FROM golang:1.26-alpine AS builder
RUN apk add --no-cache wget && \
    go install github.com/a-h/templ/cmd/templ@latest && \
    wget -qO /usr/local/bin/tailwindcss https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.19/tailwindcss-linux-x64 && \
    chmod +x /usr/local/bin/tailwindcss
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN tailwindcss -i web/static/input.css -o web/static/output.css && \
    templ generate && \
    CGO_ENABLED=0 go build -o /server ./cmd/server

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=builder /server /server
EXPOSE 8080
CMD ["/server"]
