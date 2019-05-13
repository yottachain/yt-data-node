# YTDataNode

```
ytf矿机存储节点
```

使用
```bash
# 编译
make build
# 使用示例
./out/linux-amd64-0.0.1/ytfs-node daemon 5JB7HxBrsDYtMrjJUsJ5WLdJRK3KJUrPHD6eBphVYPrXoxcqLtd
```

restful接口

serverHost = “127.0.0.1:9002”
+ 获取节点id “/api/v0/node/id”
+ 获取节点收益 “/api/v0/node/income”
+ 获取存储空间使用状况 “/api/v0/ytfs/state”