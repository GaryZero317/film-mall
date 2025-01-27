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

type DeleteAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAddressLogic {
	return &DeleteAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteAddressLogic) DeleteAddress(req *types.DeleteAddressReq) error {
	uidAny := l.ctx.Value("uid")
	if uidAny == nil {
		return errors.New("未登录")
	}

	var userId int64
	var err error
	if jsonNumber, ok := uidAny.(json.Number); ok {
		userId, err = jsonNumber.Int64()
		if err != nil {
			return err
		}
	} else {
		return errors.New("无效的用户ID")
	}

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

	err = l.svcCtx.AddressModel.Delete(l.ctx, req.Id)
	if err != nil {
		return err
	}

	return nil
}
