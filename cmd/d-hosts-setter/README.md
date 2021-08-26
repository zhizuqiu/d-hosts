# d-hosts-setter

```
> d-hosts-setter.exe -h
定时更新本地的 hosts 文件，实现自定义域名的访问

用法：
注册到 windows 服务: sc.exe create HostsSetter binPath={path-to-this-project}/d-hosts-setter.exe -a={http://d-hosts-getter-host:port} -n={hostname}

-a 指定d-hosts-getter的地址
-n 指定要更新的域名/主机
-i 指定更新间隔，单位：秒

Usage:
  d-hosts-setter [flags]

Flags:
  -a, --address string    d-hosts-getter 的地址
  -h, --help              help for d-hosts-setter
  -n, --hostname string   hostname
  -i, --interval int      更新间隔，单位：秒 (default 3600)

```

## build

windows:

```bash
go build -o dist/windows/d-hosts-setter.exe
```

## use

首先需要在一个具有公网ip的主机上运行 [d-hosts-getter](https://github.com/zhizuqiu/d-hosts/tree/master/cmd/d-hosts-getter)

1.单独运行

```
d-hosts-setter.exe -a={http://d-hosts-getter-host:port} -n={hostname}
```

2.或注册到 windows 服务

```
sc.exe create HostsSetter binPath={path-to-this-project}/d-hosts-setter.exe -a={http://d-hosts-getter-host:port} -n={hostname}
```

![create](images/sc_create.png)

查看已安装的服务，并启动：

![list](images/sc_list.png)

默认每个小时更新一下 hosts 文件：

![list](images/hosts.png)

停止、卸载服务：

![delete](images/sc_delete.png)