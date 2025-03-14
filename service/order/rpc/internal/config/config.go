package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	CacheRedis cache.CacheConf
	OrderRedis struct {
		Host string
		Type string
		Pass string
	}
	ProductRpc      zrpc.RpcClientConf
	PayRpc          zrpc.RpcClientConf
	OrderLockExpiry int // 订单锁过期时间(秒)
}
