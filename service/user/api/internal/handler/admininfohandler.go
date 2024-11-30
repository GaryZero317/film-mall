package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mall/service/user/api/internal/logic"
	"mall/service/user/api/internal/svc"
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
