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

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterRequest) (*user.RegisterResponse, error) {
	// 验证输入
	if in.Mobile == "" {
		return nil, status.Error(100, "手机号不能为空")
	}
	if in.Password == "" {
		return nil, status.Error(100, "密码不能为空")
	}
	if in.Name == "" {
		return nil, status.Error(100, "用户名不能为空")
	}

	encryptedPassword := cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, in.Password)
	l.Logger.Slowf("DEBUG Registration - Password: [%s], Salt: [%s], Encrypted: [%s]",
		in.Password, l.svcCtx.Config.Salt, encryptedPassword)

	newUser := model.User{
		Name:     in.Name,
		Gender:   in.Gender,
		Mobile:   in.Mobile,
		Password: encryptedPassword,
	}

	res, err := l.svcCtx.UserModel.Insert(l.ctx, &newUser)
	if err != nil {
		if err == model.ErrDuplicateEntry {
			return nil, status.Error(100, "该用户已存在")
		}
		return nil, status.Error(500, err.Error())
	}

	newUser.Id, err = res.LastInsertId()
	if err != nil {
		return nil, status.Error(500, err.Error())
	}

	return &user.RegisterResponse{
		Id:     newUser.Id,
		Name:   newUser.Name,
		Gender: newUser.Gender,
		Mobile: newUser.Mobile,
	}, nil
}
