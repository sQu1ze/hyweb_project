# Step 1: Build stage (編譯階段)
FROM golang:1.26.2-alpine AS builder

# 設定工作目錄
WORKDIR /app

# 安裝編譯時需要的工具
RUN apk add --no-cache git

# 複製依賴描述文件並下載套件 (利用 Docker 快取機制加速)
COPY go.mod go.sum ./
RUN go version
RUN go env GOTOOLCHAIN
RUN go mod download

# 複製專案所有原始碼
COPY . .

# 編譯成靜態執行檔，名稱定為 main
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Step 2: Run stage (執行階段)
FROM alpine:latest

# 安裝 ca-certificates (處理 HTTPS 請求需要)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 從 builder 階段複製編譯好的執行檔
COPY --from=builder /app/main .
# 複製 Swagger 文件目錄與 .env 檔
COPY --from=builder /app/cmd/server/docs ./docs
COPY --from=builder /app/.env .

# 暴露 Gin 預設的 8080 埠號
EXPOSE 8080

# 啟動應用程式
CMD ["./main"]