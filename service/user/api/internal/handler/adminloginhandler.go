package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mall/service/user/api/internal/logic"
	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"
)

func AdminLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminLoginRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewAdminLoginLogic(r.Context(), svcCtx)
		resp, err := l.AdminLogin(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
