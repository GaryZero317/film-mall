package logic

import (
	"context"
	"mall/service/statistics/api/internal/svc"
	"mall/service/statistics/api/internal/types"
	"mall/service/statistics/api/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserActivityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserActivityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserActivityLogic {
	return &GetUserActivityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserActivityLogic) GetUserActivity(req *types.UserActivityReq) (resp *types.UserActivityResp, err error) {
	// 获取时间范围
	start, end := utils.GetTimeRange(req.TimeRange)

	// 准备返回数据
	hours := utils.GetHoursList()
	resp = &types.UserActivityResp{
		Code: 0,
		Msg:  "success",
		Data: types.UserActivityData{
			Hours:    hours,
			Activity: make([][]int, 0),
		},
	}

	// 查询用户活跃度数据
	activities, err := l.svcCtx.UserActivityLogModel.FindUserActivities(l.ctx, start, end)
	if err != nil {
		l.Logger.Errorf("查询用户活跃度数据失败: %v", err)
		// 生成默认数据
		generateDefaultActivityData(resp)
		return resp, nil
	}

	// 初始化活跃度矩阵 [7][24]
	activityMatrix := make([][]int, 7)
	for i := range activityMatrix {
		activityMatrix[i] = make([]int, 24)
	}

	// 填充数据
	for _, a := range activities {
		weekday := int(a.ActivityTime.Weekday())
		if weekday == 0 {
			weekday = 7 // 将周日从0改为7
		}
		weekday-- // 转换为0-6的索引

		hour := a.ActivityTime.Hour()
		activityMatrix[weekday][hour]++
	}

	// 转换数据格式为 [hour, weekday, value]
	for weekday := 0; weekday < 7; weekday++ {
		for hour := 0; hour < 24; hour++ {
			resp.Data.Activity = append(resp.Data.Activity, []int{
				hour,
				weekday,
				activityMatrix[weekday][hour],
			})
		}
	}

	// 如果没有数据，添加默认数据
	if len(resp.Data.Activity) == 0 {
		generateDefaultActivityData(resp)
	}

	return resp, nil
}

// 生成默认的活跃度数据
func generateDefaultActivityData(resp *types.UserActivityResp) {
	for i := range resp.Data.Hours {
		for j := 0; j < 7; j++ { // 周一到周日
			resp.Data.Activity = append(resp.Data.Activity, []int{i, j, 0})
		}
	}
}
