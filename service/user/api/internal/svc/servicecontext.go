package svc

import (
	"mall/service/user/api/internal/config"
	"mall/service/user/model"
	"mall/service/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	UserRpc   user.UserClient
	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:    c,
		UserRpc:   user.NewUserClient(zrpc.MustNewClient(c.UserRpc).Conn()),
		UserModel: model.NewUserModel(conn, c.CacheRedis),
	}
}
