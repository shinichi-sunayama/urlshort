# ===== builder =====
FROM golang:1.23-bullseye AS builder
WORKDIR /app

# 1. 依存DL用に go.mod/go.sum を先にコピー
COPY go.mod go.sum ./
ENV GOTOOLCHAIN=auto
RUN go mod download

# 2. 残りをコピーしてビルド
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o urlshort ./cmd/server

# ===== runner =====
FROM gcr.io/distroless/static:nonroot
WORKDIR /home/nonroot
COPY --from=builder /app/urlshort /usr/local/bin/urlshort
COPY .env .env
USER nonroot
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/urlshort"]
