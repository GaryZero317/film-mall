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

type AdminLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminLoginLogic {
	return &AdminLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 管理员服务
func (l *AdminLoginLogic) AdminLogin(in *user.AdminLoginRequest) (*user.AdminLoginResponse, error) {
	// 查询管理员是否存在
	res, err := l.svcCtx.AdminModel.FindOneByUsername(l.ctx, in.Username)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(100, "管理员不存在")
		}
		return nil, status.Error(500, err.Error())
	}

	l.Logger.Infof("DEBUG AdminLogin - Found admin: ID=%d, Username=%s",
		res.ID, res.Username)

	// 判断密码是否正确
	password := cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, in.Password)
	if password != res.Password {
		return nil, status.Error(100, "密码错误")
	}

	response := &user.AdminLoginResponse{
		Id:       res.ID,
		Username: res.Username,
		Level:    int32(res.Level),
	}

	l.Logger.Infof("DEBUG AdminLogin - Returning response: ID=%d, Username=%s, Level=%d",
		response.Id, response.Username, response.Level)

	return response, nil
}
