package svc

import (
	"mall/service/community/api/internal/config"
	"mall/service/community/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config         config.Config
	WorkModel      model.WorkModel
	WorkImageModel model.WorkImageModel
	LikeModel      model.LikeModel
	CommentModel   model.CommentModel
	UserModel      model.UserModel
	Middleware     rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:         c,
		WorkModel:      model.NewWorkModel(conn),
		WorkImageModel: model.NewWorkImageModel(conn),
		LikeModel:      model.NewLikeModel(conn),
		CommentModel:   model.NewCommentModel(conn),
		UserModel:      model.NewUserModel(conn),
	}
}
