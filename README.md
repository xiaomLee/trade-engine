# 交易撮合引擎
只负责提供撮合服务，不涉及具体业务逻辑。撮合成功对外(mq)推送价格，返回成交记录。

# 对外提供的服务 （http, grpc）
- 接收委托
- 查询当前盘口
- 查询委托是否存在委托队列中

### 项目结构
- 最外层通过http, grpc, mq接收请求
- 所有请求进入hub, 同步等待被处理, 核心处理引擎

```
    撮合业务所限, 当前引擎采用的是内存撮合模型, 所以hub是个单点服务.
    hub支持可选
    NewMemoryHub()用于单机版, 使用机器本地内存, 重启时从数据库恢复队列.
    NewRaftHub() 是使用raft协议实现的状态机, 可用于分布式部署, 实现高可用. 
    分布式模式下, 每个节点都可响应请求, follwoer节点会自动将请求转发至leader节点, 并等待返回. 
```

- 委托队列数据结构实现在/entrust/queue/queue.go

```
    队列结构采用分段队列形式, 每个队列初始化时支持排序方式、桶大小设置. 详情参考_test.go
    
    go test --bench=.Queue_Add 用于测试写性能

    BenchmarkQueue_AddItem50-4        571498             70391 ns/op            5322 B/op          4 allocs/op
    BenchmarkQueue_AddItem100-4       444999             26144 ns/op            3565 B/op          4 allocs/op
    BenchmarkQueue_AddItem200-4       272756             10284 ns/op            5465 B/op          3 allocs/op
    BenchmarkQueue_AddItem500-4       142905             11000 ns/op           13048 B/op          3 allocs/op
    BenchmarkQueue_AddItem1000-4       95238             17672 ns/op           25039 B/op          3 allocs/op

    go test --bench=.Queue_Get 用于测试读性能

    BenchmarkQueue_Get50_2000-4               255153              4738 ns/op              32 B/op          1 allocs/op
    BenchmarkQueue_Get100_5000-4              399614              2925 ns/op              32 B/op          1 allocs/op
    BenchmarkQueue_Get200_5000-4              546223              2210 ns/op              32 B/op          1 allocs/op
    BenchmarkQueue_Get500_5000-4              413806              2733 ns/op              32 B/op          1 allocs/op
    BenchmarkQueue_Get50_100000-4               9998            125525 ns/op              32 B/op          1 allocs/op
    BenchmarkQueue_Get100_100000-4             19016             62579 ns/op              32 B/op          1 allocs/op
    BenchmarkQueue_Get200_100000-4             32082             35067 ns/op              32 B/op          1 allocs/op
    BenchmarkQueue_Get500_100000-4             53563             23206 ns/op              32 B/op          1 allocs/op
    BenchmarkQueue_Get1000_100000-4            54792             22923 ns/op              32 B/op          1 allocs/op
    
    实际工程使用时, 可通过以上基准测试选择最佳参数
```


