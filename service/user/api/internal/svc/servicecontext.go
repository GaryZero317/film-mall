package svc

import (
	"context"
	"database/sql"
	"mall/service/user/api/internal/config"
	"mall/service/user/model"
	"mall/service/user/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config               config.Config
	UserRpc              user.UserClient
	UserModel            model.UserModel
	AdminModel           *model.GormAdminModel
	CustomerServiceModel model.CustomerServiceModel
	FaqModel             model.FaqModel
	ChatMessageModel     model.ChatMessageModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	sqlDB, err := sql.Open("mysql", c.Mysql.DataSource)
	if err != nil {
		panic(err)
	}

	adminModel, err := model.NewGormAdminModel(sqlDB, c.CacheRedis)
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:               c,
		UserRpc:              user.NewUserClient(zrpc.MustNewClient(c.UserRpc).Conn()),
		UserModel:            model.NewUserModel(conn, c.CacheRedis),
		AdminModel:           adminModel,
		CustomerServiceModel: model.NewCustomerServiceModel(conn),
		FaqModel:             model.NewFaqModel(conn),
		ChatMessageModel:     model.NewChatMessageModel(conn),
	}
}

func (s *ServiceContext) GetChatHistory(ctx context.Context, userId, adminId int64, page, pageSize int) ([]*model.ChatMessage, int64, error) {
	return s.ChatMessageModel.FindByUserAndAdmin(ctx, userId, adminId, page, pageSize)
}
