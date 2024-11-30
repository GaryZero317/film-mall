package logic

import (
	"context"

	"mall/common/cryptx"
	"mall/service/user/model"
	"mall/service/user/rpc/internal/svc"
	"mall/service/user/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

type CreateAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAdminLogic {
	return &CreateAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAdminLogic) CreateAdmin(in *user.CreateAdminRequest) (*user.CreateAdminResponse, error) {
	// 检查用户名是否已存在
	_, err := l.svcCtx.AdminModel.FindOneByUsername(l.ctx, in.Username)
	if err == nil {
		return nil, status.Error(100, "管理员已存在")
	}

	if err != model.ErrNotFound {
		return nil, status.Error(500, err.Error())
	}

	l.Logger.Infof("DEBUG CreateAdmin - Creating new admin: Username=%s", in.Username)

	// 创建新管理员
	admin := &model.GormAdmin{
		Username: in.Username,
		Password: cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, in.Password),
		Level:    int(in.Level),
	}

	err = l.svcCtx.AdminModel.Insert(l.ctx, admin)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}

	response := &user.CreateAdminResponse{
		Id:       admin.ID,
		Username: admin.Username,
		Level:    int32(admin.Level),
	}

	l.Logger.Infof("DEBUG CreateAdmin - Created admin: ID=%d, Username=%s, Level=%d",
		response.Id, response.Username, response.Level)

	return response, nil
}
