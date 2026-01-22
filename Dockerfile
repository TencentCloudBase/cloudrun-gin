# 构建阶段
FROM golang:1.22.3-alpine as builder

# 安装依赖包，选用国内镜像源以提高下载速度
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tencent.com/g' /etc/apk/repositories \
    && apk add --update --no-cache git

# 指定构建过程中的工作目录
WORKDIR /app

# 将 go.mod 和 go.sum 文件拷贝到工作目录
COPY go.mod go.sum ./

# 设置 Go 模块代理
ENV GOPROXY=https://goproxy.cn,direct

# 下载依赖
RUN go mod download

# 将当前目录下所有文件都拷贝到工作目录下
COPY . .

# 执行代码编译命令。操作系统参数为linux，编译后的二进制产物命名为main
RUN GOOS=linux go build -o main .

# 运行阶段
FROM alpine:3.18

# 安装依赖包，选用国内镜像源以提高下载速度
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tencent.com/g' /etc/apk/repositories \
    && apk add --update --no-cache ca-certificates tzdata \
    && rm -f /var/cache/apk/*

# 容器默认时区为UTC，如需使用上海时间请启用以下时区设置命令
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo Asia/Shanghai > /etc/timezone

# 指定运行时的工作目录
WORKDIR /app

# 将构建产物 main 拷贝到运行时的工作目录中
COPY --from=builder /app/main /app/

# 暴露端口
EXPOSE 8080

# 设置环境变量
ENV PORT=8080
ENV GIN_MODE=release

# 执行启动命令
CMD ["./main"]
