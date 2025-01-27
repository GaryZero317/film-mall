package logic

import (
	"context"

	"mall/service/address/internal/svc"
	"mall/service/address/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddressLogic {
	return &AddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddressLogic) Address(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
