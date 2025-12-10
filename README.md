# 快速部署 Gin 应用

本篇文章为您介绍应用控制台的部署方案, 您可以通过以下操作完成部署。

## 模版部署 Gin 应用

1、登录 [腾讯云托管控制台](https://tcb.cloud.tencent.com/dev#/platform-run/service/create?type=image)

2、点击通过模版部署，选择 ```Gin 模版```

3、输入自定义服务名称，点击部署

4、等待部署完成之后，点击左上角箭头，返回到服务详情页

5、点击概述，获取默认域名并访问，会显示云托管默认首页

## 自定义部署 Gin 应用

### 创建一个 Gin 应用

创建一个新的 Gin 应用，首先需要确保机器上安装 Go 服务程序。

创建一个目录```cloudrun-gin```, 然后```cd```进入该目录。

执行如下命令在```cloudrun-gin```目录:

```sh
go mod init cloudrun-gin
go get -u github.com/gin-gonic/gin
```

在`cloudrun-gin`目录下创建 main.go 文件,内容如下:

```go
func main() {
  router := gin.Default()

  router.GET("/json", func(c *gin.Context) {
    data := map[string]interface{}{
      "id": 0,
      "name": "zhangsan",
    }

    c.JSON(http.StatusOK, data)
  })

  // 监听并在 0.0.0.0:8080 上启动服务
  router.Run(":8080")
}
```

执行 `go run main.go` 启动服务，访问`http://localhost:8080/json`可查看访问结果

### 源码

[cloudrun-gin](https://github.com/TencentCloudBase/tcbr-templates/tree/main/cloudrun-gin)

### 部署到云托管

1、在cloudrun-gin目录下创建一个名称为Dockerfile的新文件,内容如下:

```
FROM golang:1.22.3-alpine as builder

# 指定构建过程中的工作目录
WORKDIR /app

# 将当前目录（dockerfile所在目录）下所有文件都拷贝到工作目录下（.dockerignore中文件除外）
COPY . /app/

# 执行代码编译命令。操作系统参数为linux，编译后的二进制产物命名为main，并存放在当前目录下。
RUN GOOS=linux go build -o main .

# 选用运行时所用基础镜像（GO语言选择原则：尽量体积小、包含基础linux内容的基础镜像）
FROM alpine:latest

# 容器默认时区为UTC，我们使用以下时区设置命令启用上海时区
# RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo Asia/Shanghai > /etc/timezone

# 使用 HTTPS 协议访问容器云调用证书安装
RUN apk add ca-certificates

# 指定运行时的工作目录
WORKDIR /app

# 将构建产物/app/main拷贝到运行时的工作目录中
COPY --from=builder /app/main /app/

# 执行启动命令
# 写多行独立的CMD命令是错误写法！只有最后一行CMD命令会被执行，之前的都会被忽略，导致业务报错。
# 请参考[Docker官方文档之CMD命令](https://docs.docker.com/engine/reference/builder/#cmd)
CMD ["/app/main"]
```

2、进入 [腾讯云托管](https://tcb.cloud.tencent.com/dev#/platform-run/service/create?type=package)。

3、选择 ```通过本地代码``` 部署。

4、填写配置信息:

  * 代码包类型: 选择文件夹
  * 代码包: 点击选择 cloudrun-gin 目录，并上传目录文件
  * 服务名称: 填写服务名称
  * 部署类型: 选择容器服务型
  * 端口: 默认填写 8080
  * 目标目录: 默认为空
  * Dockerfile 名称: Dockerfile
  * 环境变量: 如果有按需要填写
  * 公网访问: 默认打开
  * 内网访问: 默认关闭

5、配置填写完成之后，点击部署等待部署完成。

6、部署完成之后，跳转到服务概述页面，点击默认域名进行公网访问及测试。
