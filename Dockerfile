# Stage 1: Build the Go application
FROM golang:alpine AS builder

# 设置 Go 相关环境变量
ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /build

# 复制 go.mod 和 go.sum 文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码到容器
COPY . .

# 编译 Go 应用
RUN go build -o huacache

# Stage 2: Create a minimal image with only the compiled binary
FROM scratch

# 设置时区环境变量
ENV TZ=Asia/Shanghai

# 设置工作目录
WORKDIR /app

# 复制编译后的可执行文件到新镜像中
COPY --from=builder /build/huacache /app/huacache

# 设置容器启动时执行的命令
CMD ["./huacache"]
