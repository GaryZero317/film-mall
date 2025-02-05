package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	AdminAuth struct {
		AccessSecret string
		AccessExpire int64
	}
	DB struct {
		DataSource string
	}
}
