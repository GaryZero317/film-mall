package address

import (
	"context"
	"database/sql"
	"mall/service/address/api/internal/svc"
	"mall/service/address/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetDefaultAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetDefaultAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetDefaultAddressLogic {
	return &SetDefaultAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetDefaultAddressLogic) SetDefaultAddress(req *types.SetDefaultAddressReq) error {
	userId := l.ctx.Value("userId").(int64)

	// 检查地址是否存在且属于当前用户
	address, err := l.svcCtx.AddressModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	if address.UserId != userId {
		return nil
	}

	// 将其他地址设置为非默认
	err = l.svcCtx.AddressModel.UpdateDefaultByUserId(l.ctx, userId, false)
	if err != nil {
		return err
	}

	// 将当前地址设置为默认
	address.IsDefault = 1
	err = l.svcCtx.AddressModel.Update(l.ctx, address)
	if err != nil {
		return err
	}

	return nil
}
