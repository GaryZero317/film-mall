package svc

import (
	"database/sql"
	"mall/service/user/model"
	"mall/service/user/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                   config.Config
	UserModel                model.UserModel
	AdminModel               model.AdminModel
	CustomerServiceModel     model.CustomerServiceModel
	GormCustomerServiceModel *model.GormCustomerServiceModel
	FaqModel                 model.FaqModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	// 创建SQL连接
	sqlDB, err := sql.Open("mysql", c.Mysql.DataSource)
	if err != nil {
		logx.Errorf("打开MySQL连接失败: %v", err)
		panic(err)
	}

	// 初始化GORM版本的CustomerServiceModel
	logx.Info("开始初始化GORM版本的CustomerServiceModel")
	gormCustomerServiceModel, err := model.NewGormCustomerServiceModel(sqlDB)
	if err != nil {
		logx.Errorf("初始化GORM版本的CustomerServiceModel失败: %v", err)
		panic(err)
	}
	logx.Info("GORM版本的CustomerServiceModel初始化成功")

	return &ServiceContext{
		Config:                   c,
		UserModel:                model.NewUserModel(conn, c.CacheRedis),
		AdminModel:               model.NewAdminModel(conn, c.CacheRedis),
		CustomerServiceModel:     model.NewCustomerServiceModel(conn),
		GormCustomerServiceModel: gormCustomerServiceModel,
		FaqModel:                 model.NewFaqModel(conn),
	}
}
