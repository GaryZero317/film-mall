package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	AdminAuth struct {
		AccessSecret string
		AccessExpire int64
	}

	Mysql struct {
		DataSource string
	}

	CacheRedis cache.CacheConf

	FileUpload struct {
		UploadPath   string
		UrlPrefix    string
		AllowedExts  []string
		MaxSizeBytes int64
	}
}
