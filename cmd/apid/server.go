package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ddeedev/rbac-go/client"
	"github.com/ddeedev/rbac-go/config"
	"github.com/ddeedev/rbac-go/core/service"
	"github.com/ddeedev/rbac-go/docs"
	"github.com/ddeedev/rbac-go/middleware"
	appcrypto "github.com/ddeedev/rbac-go/utils/crypto"
	"github.com/ddeedev/rbac-go/x/auth"
	"github.com/ddeedev/rbac-go/x/auth/jwt"
	authtypes "github.com/ddeedev/rbac-go/x/auth/types"
	"github.com/ddeedev/rbac-go/x/user"
	usertypes "github.com/ddeedev/rbac-go/x/user/types"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// publicMethods lists gRPC methods the auth interceptor should let through
// without a bearer token.
var publicMethods = map[string]bool{
	"/auth.AuthService/Register": true,
	"/auth.AuthService/Login":    true,
}

func interceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func Live() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
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

	// logger middleware
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	logOpts := []logging.Option{
		logging.WithLogOnEvents(logging.FinishCall),
		logging.WithDisableLoggingFields(
			"protocol", "grpc.component", "grpc.method_type",
			"grpc.start_time", "peer.address",
		),
	}

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

	// create background worker ctx
	workerCtx, workerCancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	defer func() {
		workerCancel()
		wg.Wait()
	}()

	// using waitgroup with go routine
	wg.Go(func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-workerCtx.Done():
				return
			case <-ticker.C:
				// ticker every 10 secs
				countCtx, cancel := context.WithTimeout(workerCtx, 10*time.Second)
				n, err := userRepo.Count(countCtx)
				cancel()
				if err != nil {
					if workerCtx.Err() == nil {
						log.Printf("user count: %v", err)
					}
					continue
				}
				log.Printf("total users: %d", n)
			}
		}
	})

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
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(interceptorLogger(logger), logOpts...),
			middleware.AuthMiddleware(authSvc, publicMethods),
		),
	)

	authtypes.RegisterAuthServiceServer(grpcServer, auth.NewAuthHandler(authSvc))
	usertypes.RegisterUserServiceServer(grpcServer, user.NewUserHandler(userSvc))

	// expose method on dev
	reflection.Register(grpcServer)

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

	// root must for register swagger
	root := http.NewServeMux()
	root.Handle("/swagger/", docs.Handler())

	// regisger health check
	root.HandleFunc("GET /health", Live())

	root.Handle("/", mux)

	httpServer := &http.Server{
		Addr: ":" + cfg.Port,
		// bulti-in log middleware for rest reqeuest
		Handler:           middleware.LogMiddleware(root),
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
