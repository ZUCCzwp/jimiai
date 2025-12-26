# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的工具
RUN apk add --no-cache git make

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN go build -o /app/bin/app ./cmd/app/main.go
RUN go build -o /app/bin/admin ./cmd/admin/main.go

# 运行阶段
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/bin/app /app/bin/app
COPY --from=builder /app/bin/admin /app/bin/admin

# 复制配置文件
COPY --from=builder /app/config /app/config

# 创建运行时目录
RUN mkdir -p /app/runtime/log /app/runtime/upload

# 暴露端口（app 和 admin 的端口）
EXPOSE 27355 27356

# 默认命令（可以在 docker-compose 中覆盖）
CMD ["/app/bin/app"]

