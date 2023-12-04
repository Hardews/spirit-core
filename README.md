# spirit-core
MindWrite 后端相关

## 关于本项目

本项目为 2023 年计算机设计大赛的作品后端。

## 快速开始

### 前端

在部署前，请确保你的设备具有 npm 包管理器 以及 node 环境（不使用 docker 时）

#### 项目打包并使用 nginx 部署

将代码下载下来后，进入代码主文件夹目录，执行以下命令

```Plain
npm run build
```

此时会生成一个 dist 文件夹，文件夹中含有 index.html。nginx 配置如下

```Nginx
user  nginx;
# 根据自身设备核心数配置
worker_processes  2; 
error_log  /var/log/nginx/error.log warn;
pid        /var/run/nginx.pid;
events {
  worker_connections  1024;
}
http {
  include       /etc/nginx/mime.types;
  default_type  application/octet-stream;
  log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';
  access_log  /var/log/nginx/access.log  main;
  sendfile        on;
  keepalive_timeout  65;
  server {
  # 监听的端口
    listen       8080;
    server_name  localhost;
    location / {
    # 按文件真实路径填写
      root   /spirit-front;
      index  index.html;
      try_files $uri $uri/ /index.html;
    }
  }
}
```

1. 按自身设备情况更改核心数
2. 将 index 的 index.html 改成你的 index.html 对应的地址
3. 启动 nginx
4. 访问 ip + :8080 即可访问前端页面

#### 使用 nginx + docker 部署（推荐使用）

将代码下载下来后，会有一个 Dockerfile 的文件，在此文件的目录下，执行如下命令：

```Nginx
docker build -t mind-write .
```

待上述命令执行完毕后，执行如下命令：

```Nginx
docker run -d -p 8080:80 mind-write
```

即部署成功。可更改 8080 为其他端口。

请注意，若要在服务器上成功访问，需要开放相关端口。

附 nginx 配置及 Dockerfile

nginx：

```Nginx
user  nginx;
# 根据自身设备核心数配置
worker_processes  2; 
error_log  /var/log/nginx/error.log warn;
pid        /var/run/nginx.pid;
events {
  worker_connections  1024;
}
http {
  include       /etc/nginx/mime.types;
  default_type  application/octet-stream;
  log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';
  access_log  /var/log/nginx/access.log  main;
  sendfile        on;
  keepalive_timeout  65;
  server {
  # 监听的端口
    listen       80;
    server_name  localhost;
    location / {
    # 按文件真实路径填写
      root   /spirit-front;
      index  index.html;
      try_files $uri $uri/ /index.html;
    }
  }
}
```

Dockerfile:

```Dockerfile
FROM node:latest
COPY . /spirit-front
WORKDIR /spirit-front
RUN npm install --legacy-peer-deps && npm run build

FROM nginx
RUN mkdir /spirit-front
COPY --from=0 /spirit-front/dist /spirit-front
COPY nginx.conf /etc/nginx/nginx.conf
```

### 后端

请注意，如需要部署后端服务，请确保你的设备有 Go、Mysql、Redis 环境

如果想要自己部署后端服务，可按照以下流程

#### 部署之前

##### 代码修改

你需要修改相关配置的代码

文件路径 `dao/dao.go` 

> 18 行 address = 容器名称 + ":" + 端口号
>
> 19 行 dbName = 创建的数据库名称
>
> 23 行 gmt_website 改为 对应的用户名（如 root）

文件路径 `tool/redi.go`

> 19 行 Addr:     "localhost:6379", // 改为你的 redis 数据库名称 + ":" + 端口号 20 行 Password: "", // 改为你的密码 21 行 DB:       0,  // 改为你想用的数据库

并且，在前端代码中，你需要代码进行如下修改：

文件路径 `src/main.js`

> 13 行 axios.defaults.baseURL = '' // 将路径改为你后端部署的路径

##### 环境变量配置

亦或者设置环境变量(找到上述所说的代码行，根据名字修改即可)

#### 普通部署

将代码下载下来后，在代码主文件夹内执行以下命令

```Nginx
go build -o spirit-core
```

然后可以得到一个可执行文件 spirit-core，此时，可执行命令

```Nginx
nohup ./spirit-core &
```

或使用宝塔 Linux 进行部署

注：spirit-core 是 MindWrite 的后端服务名称，意为精神的核心，也是 MindWrite 的核心。

#### docker 部署

在项目文件夹中创建 Dockerfile 并添加以下内容：

```Dockerfile
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
```

#### 数据库配置

你可以使用本机的数据库，也可以使用 docker 中的数据库，这里只介绍 docker 数据库配置。

##### Mysql

执行以下命令

```Nginx
docker pull mysql:latest
```

然后运行这个 mysql 数据库,根据需要修改密码与名称

```Nginx
docker run -p 3307:3306 --name mysql  -v /Users/hyc/DockerStudy/mysql/log:/var/log/mysql  -v /Users/hyc/DockerStudy/mysql/data:/var/lib/mysql  -v /Users/hyc/DockerStudy/mysql/conf:/etc/mysql  -e MYSQL_ROOT_PASSWORD=root  -d mysql:latest
```

运行成功后，进入该容器并创建对应的数据库。数据库的 DDL 语句放在最后，各表信息请见详细设计。

##### Redis

执行以下命令

```Nginx
docker pull redis:latest
```

然后运行这个 mysql 数据库,根据需要修改密码与名称及相关参数

```Nginx
docker run --restart=always --log-opt max-size=100m --log-opt max-file=2 -p 6379:6379 --name redisname -v /DockerContainerProperties/redis/myredis.conf:/etc/redis/redis.conf -v /DockerContainerProperties/redis/data:/data -d redis redis-server /etc/redis/redis.conf  --appendonly yes  --requirepass 888888
```

