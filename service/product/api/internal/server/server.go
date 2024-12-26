package server

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterStaticHandlers(server *rest.Server) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/uploads/:file",
				Handler: http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))).ServeHTTP,
			},
		},
	)
}
