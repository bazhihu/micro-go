# micro-go
```
    高并发和微服务kit实战
```
###

+ 服务注册与发现

```
    简单的字符串组合和相差案例
    1、组合字符串
    2、相差字符串
    3、健康检查
    调用案例
    http://127.0.0.1:10085/op/Concat/asda/asdasd
    http://127.0.0.1:10085/op/Diff/asda/asdasd
```

- 远程过程调用rpc


+ 分布式配置中心
    * spring cloud config
    * etcd
    * 分布式锁 
    
+ 微服务网关
    * Nginx
        * nginx 设置反向代理转发到另外的服务
        * 拥有网关服务 服务接收请求根据服务发现自动转发到正确服务
    * Zuul  
    * Kong

  
+ 微服务的容错处理与负载均衡
    * 服务熔断 
    * 负载均衡
```markdown
    简单的方式
    分别启动调用方use-string-service和被调用方string-service
    curl -X POST http://127.0.0.1:10086/op/Concat/qw/er
    关闭string-service 在发起请求响应{"error":"hystrix: circuit open"}
```

+ 统一认证与授权


+ 分布式链路追踪

```
    分布式追踪系统发展
    1、代码埋点
    2、数据存储
    3、查询展示
    
    Zipkin
    基本概念
    1、Span (基本工作单元)
        一次链路调用创建一个Span, 一个64位ID标识Span
        ParentID 表示Span调用链路的来源
    2、Trace （类似于树结构的Span集合）
        一条完整的调用链路，存在唯一的标识TraceID
    3、Annotation（注解）
        CS:Client Sent ，表示客户端发起请求
        SR:Server Receive ， 表示服务端受到请求
        SS:Server Send ， 表示服务端完成处理，并将结果发送给客户端
        CR:Client Received ， 表示客户端获取到服务端返回信息
    
    TraceID
    SpanID
    ParentID
```

### 目录结构
####transport层: 项目提供服务的方式（HTTP服务）
    主要负责网络传输，例如处理HTTP、gRPC、Thrift 相关逻辑。包含请求的参数格式转换。
####endpoint层: 用于接受请求并返回响应
    主要负责请求响应的request/response格式的转换，以及公用拦截器相关的逻辑
    并且提供对日志、限流、熔断、链路追踪和服务监控等扩展能力
####service层: 业务代码实现层
    主要负责于业务逻辑
```
Go-kit提供一下功能
- 熔断器 Circuit breaker
- 限流器 Rate limiter
- 日志 Logging
- Prometheus 统计 Metrics
- 请求跟踪 Request tracing
- 服务发现和负载均衡 
```


## 综合实战 -- 秒杀系统
  
