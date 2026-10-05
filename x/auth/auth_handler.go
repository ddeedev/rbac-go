package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ddeedev/rbac-go/core/ports"
	"github.com/ddeedev/rbac-go/utils"
	"github.com/ddeedev/rbac-go/validation"
	"github.com/ddeedev/rbac-go/x/auth/types"
	usertypes "github.com/ddeedev/rbac-go/x/user/types"
)

var _ types.AuthServiceServer = (*AuthHandler)(nil)

// handles authentication endpoints.
type AuthHandler struct {
	types.UnimplementedAuthServiceServer
	svc ports.AuthService
}

func NewAuthHandler(svc ports.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(ctx context.Context, req *types.RegisterRequest) (*types.RegisterResponse, error) {
	in := usertypes.CreateUserInput{Name: req.Name, Email: req.Email, Password: req.Password}
	if err := validation.Struct(in); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err := h.svc.Register(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		return nil, utils.ToStatus(err)
	}

	return &types.RegisterResponse{
		Success: true,
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *types.LoginRequest) (*types.LoginResponse, error) {
	in := types.LoginBodyDTO{Username: req.Username, Password: req.Password}
	if err := validation.Struct(in); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res, err := h.svc.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, utils.ToStatus(err)
	}
	return &types.LoginResponse{
		AccessToken: res.AccessToken,
		ExpiresIn:   int32(res.ExpiresIn.Seconds()),
		User: &usertypes.User{
			Id: res.User.ID, Name: res.User.Name, Email: res.User.Email,
		},
	}, nil
}
