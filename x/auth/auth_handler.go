package auth

import (
	"context"

	"github.com/ddeedev/rbac-go/core/ports"
	"github.com/ddeedev/rbac-go/utils"
	"github.com/ddeedev/rbac-go/x/auth/types"
	usertypes "github.com/ddeedev/rbac-go/x/user/types"
)

var _ types.AuthServiceServer = (*AuthHandler)(nil)

// handles authentication endpoints.
type AuthHandler struct {
	types.UnimplementedAuthServiceServer
	svc ports.AuthService
}

func NewAuthHandler(svc ports.AuthService) *AuthHandler { return &AuthHandler{svc: svc} }

func (h *AuthHandler) Login(ctx context.Context, req *types.LoginRequest) (*types.LoginResponse, error) {
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
