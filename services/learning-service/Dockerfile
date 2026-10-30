# 1. ビルド（Go のツール一式が入った大きなイメージ）
FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/server .

# 2. 実行（できあがったプログラムだけを小さなイメージに入れる）
FROM alpine:3.22
COPY --from=builder /bin/server /bin/server
EXPOSE 8080
CMD ["/bin/server"]
