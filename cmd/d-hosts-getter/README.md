# d-hosts-getter

定时地访问此服务的 `POST /curl` 路径，将路由器ip记录在此服务中，然后通过 `GET /ip` 路径获取路由器ip

## build

```bash
GOOS=linux GOARCH=amd64 go build -o dist/amd64/d-hosts-getter
```

## run

```
./d-hosts-getter
```

or docker:

```
docker run -d -p 8007:3000 --restart=always zhizuqiu/d-hosts-getter:latest
```

## use

在路由器的定时任务中设置：

```
* */1 * * * curl -X POST http://your_app_host:app_port/curl
```

1.在浏览器中访问 `GET /ip` 接口，获取ip：

![demo.png](demo.png)

2.或者使用 [d-hosts-setter](https://github.com/zhizuqiu/d-hosts/tree/master/cmd/d-hosts-setter) 定时更新本地的 hosts 文件，实现自定义域名的访问