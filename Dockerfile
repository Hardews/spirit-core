FROM golang:latest  AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

ENV GOPROXY=https://goproxy.cn,direct

# 移动到工作目录 build
WORKDIR /build

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
COPY go.mod .
COPY go.sum .
RUN go mod download

# 将代码复制到容器中
COPY . .

# 将代码编译为可执行文件到blog
RUN go build -o spirit-core .

# 创建一个小镜像
#FROM scratch
FROM busybox

COPY --from=builder /build/spirit-core /

EXPOSE 8089
EXPOSE 3306
EXPOSE 6379
ENTRYPOINT ["/spirit-core"]