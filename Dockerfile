# 构建阶段
FROM golang:1.25-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装依赖工具
RUN apk add --no-cache git make

# 复制go mod文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cloud-disk ./cmd/server/main.go

# 运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN adduser -D -g '' appuser

WORKDIR /app

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/cloud-disk .

# 创建必要的目录
RUN mkdir -p /app/storage/uploads && \
    mkdir -p /app/storage/uploads/temp && \
    chown -R appuser:appuser /app

# 复制配置文件
COPY --from=builder /app/config.prod.yaml ./config.yaml

# 切换到非root用户
USER appuser

# 或者为了兼容性，同时复制多个配置
COPY config.dev.yaml /app/config.dev.yaml
COPY config.prod.yaml /app/config.prod.yaml
COPY config.yaml /app/config.yaml

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动命令
CMD ["./cloud-disk"]