// Code generated by goctl. DO NOT EDIT!
// Source: user.proto

package server

import (
	"context"

	"mall/service/user/rpc/internal/logic"
	"mall/service/user/rpc/internal/svc"
	"mall/service/user/rpc/pb/user"
)

type UserServer struct {
	svcCtx *svc.ServiceContext
	user.UnimplementedUserServer
}

func NewUserServer(svcCtx *svc.ServiceContext) *UserServer {
	return &UserServer{
		svcCtx: svcCtx,
	}
}

func (s *UserServer) Login(ctx context.Context, in *user.LoginRequest) (*user.LoginResponse, error) {
	l := logic.NewLoginLogic(ctx, s.svcCtx)
	return l.Login(in)
}

func (s *UserServer) Register(ctx context.Context, in *user.RegisterRequest) (*user.RegisterResponse, error) {
	l := logic.NewRegisterLogic(ctx, s.svcCtx)
	return l.Register(in)
}

func (s *UserServer) UserInfo(ctx context.Context, in *user.UserInfoRequest) (*user.UserInfoResponse, error) {
	l := logic.NewUserInfoLogic(ctx, s.svcCtx)
	return l.UserInfo(in)
}

//  管理员服务
func (s *UserServer) AdminLogin(ctx context.Context, in *user.AdminLoginRequest) (*user.AdminLoginResponse, error) {
	l := logic.NewAdminLoginLogic(ctx, s.svcCtx)
	return l.AdminLogin(in)
}

func (s *UserServer) CreateAdmin(ctx context.Context, in *user.CreateAdminRequest) (*user.CreateAdminResponse, error) {
	l := logic.NewCreateAdminLogic(ctx, s.svcCtx)
	return l.CreateAdmin(in)
}

func (s *UserServer) UpdateAdmin(ctx context.Context, in *user.UpdateAdminRequest) (*user.UpdateAdminResponse, error) {
	l := logic.NewUpdateAdminLogic(ctx, s.svcCtx)
	return l.UpdateAdmin(in)
}

func (s *UserServer) DeleteAdmin(ctx context.Context, in *user.DeleteAdminRequest) (*user.DeleteAdminResponse, error) {
	l := logic.NewDeleteAdminLogic(ctx, s.svcCtx)
	return l.DeleteAdmin(in)
}

func (s *UserServer) AdminInfo(ctx context.Context, in *user.AdminInfoRequest) (*user.AdminInfoResponse, error) {
	l := logic.NewAdminInfoLogic(ctx, s.svcCtx)
	return l.AdminInfo(in)
}

//  客服服务
func (s *UserServer) SubmitCustomerService(ctx context.Context, in *user.CustomerServiceRequest) (*user.CustomerServiceResponse, error) {
	l := logic.NewSubmitCustomerServiceLogic(ctx, s.svcCtx)
	return l.SubmitCustomerService(in)
}

func (s *UserServer) GetServiceList(ctx context.Context, in *user.ServiceListRequest) (*user.ServiceListResponse, error) {
	l := logic.NewGetServiceListLogic(ctx, s.svcCtx)
	return l.GetServiceList(in)
}

func (s *UserServer) GetServiceDetail(ctx context.Context, in *user.ServiceDetailRequest) (*user.ServiceDetailResponse, error) {
	l := logic.NewGetServiceDetailLogic(ctx, s.svcCtx)
	return l.GetServiceDetail(in)
}

func (s *UserServer) GetFaqList(ctx context.Context, in *user.FaqListRequest) (*user.FaqListResponse, error) {
	l := logic.NewGetFaqListLogic(ctx, s.svcCtx)
	return l.GetFaqList(in)
}
