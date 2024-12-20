package handler

import (
	"net/http"

	"mall/service/user/api/internal/logic"
	"mall/service/user/api/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func AdminInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAdminInfoLogic(r.Context(), svcCtx)
		resp, err := l.AdminInfo()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
