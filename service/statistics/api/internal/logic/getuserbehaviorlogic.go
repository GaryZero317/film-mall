package logic

import (
	"context"
	"mall/service/statistics/api/internal/svc"
	"mall/service/statistics/api/internal/types"
	"mall/service/statistics/api/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserBehaviorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserBehaviorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserBehaviorLogic {
	return &GetUserBehaviorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserBehaviorLogic) GetUserBehavior(req *types.UserBehaviorReq) (resp *types.UserBehaviorResp, err error) {
	// 获取时间范围
	start, end := utils.GetTimeRange(req.TimeRange)

	// 获取日期列表
	dates := utils.GetDatesBetween(start, end)

	// 准备返回数据
	resp = &types.UserBehaviorResp{
		Code: 0,
		Msg:  "success",
		Data: types.UserBehaviorData{
			Dates:  dates,
			Views:  make([]int, len(dates)),
			Carts:  make([]int, len(dates)),
			Orders: make([]int, len(dates)),
		},
	}

	// 构建日期映射
	dateMap := make(map[string]int)
	for i, date := range dates {
		dateMap[date] = i
	}

	// 查询用户行为数据
	behaviors, err := l.svcCtx.UserActivityLogModel.FindUserBehaviors(l.ctx, start, end)
	if err != nil {
		l.Logger.Errorf("查询用户行为数据失败: %v", err)
		// 出错时直接返回空数据(已经初始化为0了)
		return resp, nil
	}

	// 填充数据
	for _, b := range behaviors {
		date := b.ActivityTime.Format("2006-01-02")
		if idx, ok := dateMap[date]; ok {
			switch b.ActivityType {
			case "view":
				resp.Data.Views[idx]++
			case "cart":
				resp.Data.Carts[idx]++
			case "order":
				resp.Data.Orders[idx]++
			}
		}
	}

	return resp, nil
}
