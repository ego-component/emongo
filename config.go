package emongo

import (
	"time"

	"github.com/gotomicro/ego/core/util/xtime"
)

type config struct {
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
	interceptors               []Interceptor
	keyName                    string
	dbName                     string
}

// DefaultConfig 返回默认配置
func DefaultConfig() *config {
	return &config{
		DSN:                     "",
		Debug:                   false,
		DialTimeout:             xtime.Duration("10s"),
		SocketTimeout:           xtime.Duration("300s"),
		SlowLogThreshold:        xtime.Duration("600ms"),
		MinPoolSize:             0,
		MaxPoolSize:             300,
		EnableMetricInterceptor: true,
		EnableTraceInterceptor:  true,
	}
}
