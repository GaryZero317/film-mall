package address

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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

func (l *SetDefaultAddressLogic) SetDefaultAddress(req *types.SetDefaultAddressReq) (*types.SetDefaultAddressResp, error) {
	uidAny := l.ctx.Value("uid")
	if uidAny == nil {
		return nil, errors.New("未登录")
	}

	var userId int64
	var err error
	if jsonNumber, ok := uidAny.(json.Number); ok {
		userId, err = jsonNumber.Int64()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("无效的用户ID")
	}

	// 检查地址是否存在且属于当前用户
	address, err := l.svcCtx.AddressModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("地址不存在")
		}
		return nil, err
	}

	if address.UserId != userId {
		return nil, errors.New("无权操作此地址")
	}

	// 将其他地址设置为非默认
	err = l.svcCtx.AddressModel.UpdateDefaultByUserId(l.ctx, userId, false)
	if err != nil {
		return nil, err
	}

	// 将当前地址设置为默认
	address.IsDefault = 1
	err = l.svcCtx.AddressModel.Update(l.ctx, address)
	if err != nil {
		return nil, err
	}

	return &types.SetDefaultAddressResp{
		Code:    0,
		Message: "设置成功",
	}, nil
}
