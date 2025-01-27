package address

import (
	"context"
	"database/sql"

	"mall/service/address/api/internal/svc"
	"mall/service/address/api/internal/types"
	"mall/service/address/api/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAddressLogic {
	return &UpdateAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAddressLogic) UpdateAddress(req *types.UpdateAddressReq) error {
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

	// 如果设置为默认地址，需要将其他地址设置为非默认
	if req.IsDefault {
		err = l.svcCtx.AddressModel.UpdateDefaultByUserId(l.ctx, userId, false)
		if err != nil {
			return err
		}
	}

	// 更新地址信息
	if req.Name != "" {
		address.Name = req.Name
	}
	if req.Phone != "" {
		address.Phone = req.Phone
	}
	if req.Province != "" {
		address.Province = req.Province
	}
	if req.City != "" {
		address.City = req.City
	}
	if req.District != "" {
		address.District = req.District
	}
	if req.DetailAddr != "" {
		address.DetailAddr = req.DetailAddr
	}
	address.IsDefault = utils.BoolToInt64(req.IsDefault)

	err = l.svcCtx.AddressModel.Update(l.ctx, address)
	if err != nil {
		return err
	}

	return nil
}
