package auth

import (
	"context"
	"strings"

	"github.com/ddeedev/rbac-go/core/domain"
	"github.com/ddeedev/rbac-go/core/ports"
	"github.com/ddeedev/rbac-go/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthMiddleware(auth ports.AuthService, public map[string]bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo,
		next grpc.UnaryHandler) (any, error) {

		if public[info.FullMethod] {
			return next(ctx, req)
		}

		md, _ := metadata.FromIncomingContext(ctx)
		vals := md.Get("authorization")
		if len(vals) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		scheme, token, ok := strings.Cut(vals[0], " ")
		if !ok || !strings.EqualFold(scheme, "Bearer") || strings.TrimSpace(token) == "" {
			return nil, status.Error(codes.Unauthenticated, "expected: Authorization: Bearer <token>")
		}

		claims, err := auth.ValidateToken(ctx, strings.TrimSpace(token))
		if err != nil {
			return nil, utils.ToStatus(err)
		}

		return next(domain.WithClaims(ctx, claims), req)
	}
}
