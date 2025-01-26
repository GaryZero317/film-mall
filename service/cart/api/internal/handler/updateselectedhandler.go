package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mall/service/cart/api/internal/logic"
	"mall/service/cart/api/internal/svc"
	"mall/service/cart/api/internal/types"
)

func UpdateSelectedHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateSelectedReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewUpdateSelectedLogic(r.Context(), svcCtx)
		resp, err := l.UpdateSelected(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
