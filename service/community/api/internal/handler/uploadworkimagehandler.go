package handler

import (
	"net/http"

	"mall/service/community/api/internal/logic"
	"mall/service/community/api/internal/svc"
	"mall/service/community/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func uploadWorkImageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadWorkImageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewUploadWorkImageLogic(r.Context(), svcCtx, r)
		resp, err := l.UploadWorkImage(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
