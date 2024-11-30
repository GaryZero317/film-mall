package handler

import (
	"mall/common/errorx"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"mall/service/user/api/internal/logic"
	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"
)

func UpdateAdminHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateAdminRequest
		if err := httpx.Parse(r, &req); err != nil {
			logx.Errorf("Parse request error: %v", err)
			httpx.Error(w, errorx.NewDefaultError("请求参数错误"))
			return
		}

		// 打印请求内容用于调试
		logx.Infof("Request body: %+v", req)

		l := logic.NewUpdateAdminLogic(r.Context(), svcCtx)
		resp, err := l.UpdateAdmin(&req)
		if err != nil {
			logx.Errorf("UpdateAdmin error: %v", err)
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
