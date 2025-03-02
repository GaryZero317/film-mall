package handler

import (
	"net/http"

	"mall/service/user/api/internal/logic"
	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AdminGetChatSessionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAdminGetChatSessionsLogic(r.Context(), svcCtx)
		resp, err := l.AdminGetChatSessions()
		if err != nil {
			// 记录错误但返回空列表
			logx.Errorf("获取聊天会话失败: %v", err)
			// 返回空列表而不是错误
			httpx.OkJson(w, &types.ChatSessionListResponse{
				Total: 0,
				List:  make([]types.ChatSessionItem, 0),
			})
		} else {
			// 确保返回有效的list字段
			if resp == nil {
				resp = &types.ChatSessionListResponse{
					Total: 0,
					List:  make([]types.ChatSessionItem, 0),
				}
			}
			httpx.OkJson(w, resp)
		}
	}
}
