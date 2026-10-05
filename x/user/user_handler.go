package user

import (
	"context"
	"time"

	"github.com/ddeedev/rbac-go/core/domain"
	"github.com/ddeedev/rbac-go/core/ports"
	"github.com/ddeedev/rbac-go/utils"
	"github.com/ddeedev/rbac-go/x/user/types"
)

var _ types.UserServiceServer = (*UserHandler)(nil)

// UserHandler exposes the user service over gRPC / the HTTP gateway.
type UserHandler struct {
	types.UnimplementedUserServiceServer
	svc ports.UserService
}

func NewUserHandler(svc ports.UserService) *UserHandler { return &UserHandler{svc: svc} }

func (h *UserHandler) Create(ctx context.Context, req *types.CreateUserRequest) (*types.UserResponse, error) {
	u, err := h.svc.Create(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		return nil, utils.ToStatus(err)
	}
	return toResponse(u), nil
}

func (h *UserHandler) Get(ctx context.Context, req *types.GetUserRequest) (*types.UserResponse, error) {
	u, err := h.svc.Get(ctx, req.Id)
	if err != nil {
		return nil, utils.ToStatus(err)
	}
	return toResponse(u), nil
}

func (h *UserHandler) Update(ctx context.Context, req *types.UpdateUserRequest) (*types.UserResponse, error) {
	u, err := h.svc.Update(ctx, req.Id, req.Name, req.Email)
	if err != nil {
		return nil, utils.ToStatus(err)
	}
	return toResponse(u), nil
}

func (h *UserHandler) Delete(ctx context.Context, req *types.DeleteUserRequest) (*types.UserResponse, error) {
	u, err := h.svc.Delete(ctx, req.Id)
	if err != nil {
		return nil, utils.ToStatus(err)
	}
	return toResponse(u), nil
}

func (h *UserHandler) List(ctx context.Context, req *types.ListUserRequest) (*types.ListUserResponse, error) {
	us, err := h.svc.List(ctx)
	if err != nil {
		return nil, utils.ToStatus(err)
	}

	var res []*types.User

	for _, user := range us {
		res = append(res, &types.User{
			Name:      user.Name,
			Email:     user.Email,
			Id:        user.ID,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		})
	}

	return &types.ListUserResponse{
		Users: res,
	}, nil
}

func toResponse(u *domain.User) *types.UserResponse {
	return &types.UserResponse{
		Id:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}
