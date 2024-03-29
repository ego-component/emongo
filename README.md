# emongo 组件使用指南
[![goproxy.cn](https://goproxy.cn/stats/github.com/ego-component/emongo/badges/download-count.svg)](https://goproxy.cn/stats/github.com/ego-component/emongo)
[![Release](https://img.shields.io/github/v/release/ego-component/emongo.svg?style=flat-square)](https://github.com/ego-component/emongo)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Example](https://img.shields.io/badge/Examples-2ca5e0?style=flat&logo=appveyor)](https://github.com/ego-component/emongo/tree/master/examples)
[![Doc](https://img.shields.io/badge/Docs-1?style=flat&logo=appveyor)](https://ego.gocn.vip/frame/client/gorm.html#_1-%E7%AE%80%E4%BB%8B)


## 1 简介
对 [mongo-driver](https://godoc.org/go.mongodb.org/mongo-driver) 进行了轻量封装，并提供了以下功能：
- 规范了标准配置格式，提供了统一的 Load().Build() 方法。
- 支持自定义拦截器
- 提供了默认的 Debug 拦截器，开启 Debug 后可输出 Request、Response 至终端。
- 提供了默认的 Metric 拦截器，开启后可采集 Prometheus 指标数据

## 2 使用方式
```bash
go get github.com/ego-component/emongo
```

## 3 mongo配置
```go
type Config struct {
    DSN                        string        `json:"dsn" toml:"dsn"`     // DSN DSN地址
    Debug                      bool          `json:"debug" toml:"debug"` // Debug 是否开启debug模式
    DialTimeout                time.Duration // 连接超时
    SocketTimeout              time.Duration `json:"socketTimeout" toml:"socketTimeout"` // SocketTimeout 创建连接的超时时间
    MaxConnIdleTime            time.Duration `json:"maxConnIdleTime"`
    MinPoolSize                int           // MinPoolSize 连接池大小(最小连接数)
    MaxPoolSize                int           `json:"maxPoolSize" toml:"maxPoolSize"`                               // MaxPoolSize 连接池大小(最大连接数)
    EnableMetricInterceptor    bool          `json:"enableMetricInterceptor" toml:"enableMetricInterceptor"`       // EnableMetricInterceptor 是否启用prometheus metric拦截器
    EnableAccessInterceptorReq bool          `json:"enableAccessInterceptorReq" toml:"enableAccessInterceptorReq"` // EnableAccessInterceptorReq 是否启用access req拦截器，此配置只有在EnableAccessInterceptor=true时才会生效
    EnableAccessInterceptorRes bool          `json:"enableAccessInterceptorRes" toml:"enableAccessInterceptorRes"` // EnableAccessInterceptorRes 是否启用access res拦截器，此配置只有在EnableAccessInterceptor=true时才会生效
    EnableAccessInterceptor    bool          `json:"enableAccessInterceptor" toml:"enableAccessInterceptor"`       // EnableAccessInterceptor 是否启用access拦截器
    EnableTraceInterceptor     bool          `json:"enableTraceInterceptor" toml:"enableTraceInterceptor"`         // EnableTraceInterceptor 是否启用trace拦截器
    SlowLogThreshold           time.Duration // SlowLogThreshold 慢日志门限值，超过该门限值的请求，将被记录到慢日志中
    // TLS 支持
    Authentication Authentication
}
```

## 4 优雅的Debug
通过开启``debug``配置和命令行的``export EGO_DEBUG=true``，我们就可以在测试环境里看到请求里的配置名、地址、耗时、请求数据、响应数据
![img.png](https://cdn.gocn.vip/ego/assets/img/mongo1.d5d9680d.png)


## 5 用户配置
```toml
[mongo]
  debug=true
  dsn="mongodb://user:password@localhost:27017,localhost:27018"
  [mongo.authentication]
    [mongo.authentication.tls]
      enabled=false
      CAFile=""
      CertFile="./cert/tls.pem"
      KeyFile="./cert/tls.key"
      insecureSkipVerify=true
```

## 6 用户代码
```go
var stopCh = make(chan bool)
	// 假设你配置的toml如下所示
	conf := `
[mongo]
	debug=true
	dsn="mongodb://user:password@localhost:27017,localhost:27018"
`
	// 加载配置文件
err := econf.LoadFromReader(strings.NewReader(conf), toml.Unmarshal)
if err != nil {
    panic("LoadFromReader fail," + err.Error())
}

// 初始化emongo组件
cmp := emongo.Load("mongo").Build()
coll := cmp.Client.Database("test").Collection("cells")
findOne(coll)

stopCh <- true
```

