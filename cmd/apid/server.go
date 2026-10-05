package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ddeedev/rbac-go/client"
	"github.com/ddeedev/rbac-go/config"
	"github.com/ddeedev/rbac-go/core/service"
	appcrypto "github.com/ddeedev/rbac-go/utils/crypto"
	"github.com/ddeedev/rbac-go/x/auth"
	"github.com/ddeedev/rbac-go/x/auth/jwt"
	authtypes "github.com/ddeedev/rbac-go/x/auth/types"
	"github.com/ddeedev/rbac-go/x/user"
	usertypes "github.com/ddeedev/rbac-go/x/user/types"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// publicMethods lists gRPC methods the auth interceptor should let through
// without a bearer token.
var publicMethods = map[string]bool{
	"/auth.AuthService/Login": true,
}

func start() error {
	// load server env config
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	// process clean-up
	defer stop()

	// adapters
	// database conntection
	mongoClient, err := client.ConnectToMongoDB(&cfg.Database)
	if err != nil {
		return err
	}

	// database connection finalize clean-up
	// graceful disconnect db after signal
	defer func() {
		disconnectCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := mongoClient.Disconnect(disconnectCtx); err != nil {
			log.Printf("failed to disconnect mongo client: %v", err)
		}
	}()

	// crate db connection instance
	db := mongoClient.Database(cfg.Database.Name)

	// init user db repository
	userRepo, err := user.NewUserRepo(ctx, db)
	if err != nil {
		return err
	}

	// password and jwt manager
	hasher := appcrypto.NewHasher()
	tokenMngr, err := jwt.NewManager(cfg.JWTSecret, cfg.JWTTokenTTL)
	if err != nil {
		return err
	}

	// register service
	userSvc := service.NewUserService(userRepo, hasher)
	authSvc := service.NewAuthService(userRepo, hasher, tokenMngr)

	// init grpc service
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(auth.AuthMiddleware(authSvc, publicMethods)),
	)
	authtypes.RegisterAuthServiceServer(grpcServer, auth.NewAuthHandler(authSvc))
	usertypes.RegisterUserServiceServer(grpcServer, user.NewUserHandler(userSvc))

	// add listener
	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		return err
	}

	// create channel with buffered 2
	// one slot per server
	errCh := make(chan error, 2)

	go func() {
		log.Printf("gRPC server listening on :%s", cfg.GRPCPort)
		if err := grpcServer.Serve(listener); err != nil {
			errCh <- fmt.Errorf("grpc serve: %w", err)
		}
	}()

	// http gateway
	// own context outlives the signal ctx so in-flight REST calls can finish.
	gwCtx, gwCancel := context.WithCancel(context.Background())
	defer gwCancel()

	mux := runtime.NewServeMux()
	endpoint := "localhost:" + cfg.GRPCPort
	dialOps := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := authtypes.RegisterAuthServiceHandlerFromEndpoint(gwCtx, mux, endpoint, dialOps); err != nil {
		return err
	}

	if err := usertypes.RegisterUserServiceHandlerFromEndpoint(gwCtx, mux, endpoint, dialOps); err != nil {
		return err
	}

	httpServer := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("HTTP gateway listening on :%s", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("http serve: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		// now kills the process immediately
		stop()
		log.Println("shutting down...")
	case err := <-errCh:
		// deferred Disconnect still runs
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown: %v", err)
	}

	// safe to drop the gateway grpc connections now
	gwCancel()

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
	case <-shutdownCtx.Done():
		log.Println("graceful stop timed out, forcing")
		grpcServer.Stop()
	}

	return nil
}
