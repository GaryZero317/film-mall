package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mall/service/statistics/api/internal/logic"
	"mall/service/statistics/api/internal/svc"
	"mall/service/statistics/api/internal/types"
)

func getUserBehaviorHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserBehaviorReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetUserBehaviorLogic(r.Context(), svcCtx)
		resp, err := l.GetUserBehavior(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
