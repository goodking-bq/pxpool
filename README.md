## pxpool 说明

这是一个练手项目，本地ip代理池，

### 爬虫支持

- [快代理](https://www.kuaidaili.com)
- ... 

### ip 端口 扫描

- 扫描IP段
- 自定义端口

### api

- /random/ 随机获取一个代理
- /counter/ 拥有的代理数

### 数据库

使用的bbolt本地KV数据库，方便快速

## 使用
```shell
D:\go-program>go run .\src\pxpool -h
pxpool 代理池全功能支持 网站爬虫 端口扫描 web api

Usage:
  pxpool [flags]
  pxpool [command]

Available Commands:
  crawl       启动网络爬虫
  help        Help about any command
  scan        启动网络爬虫
  web         启动web api等

Flags:
  -c, --config string        配置文件 (默认是 config.yaml)
  -h, --help                 help for pxpool
      --log string           日志文件路径 (默认是 os.stdout)
      --loglevel string      日志级别 (默认是 info)
  -t, --storageType string   数据库存储类型 (默认是 bolt) (default "bolt")

Use "pxpool [command] --help" for more information about a command.
```

## 更新计划

- [ ] 添加更多的爬虫
- [ ] 分布式扫描
- [ ] web支持