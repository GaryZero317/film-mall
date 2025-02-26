package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	AdminAuth struct {
		AccessSecret string
		AccessExpire int64
	}

	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	Mysql struct {
		DataSource string
	}
}
