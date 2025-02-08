package server

import (
	"context"
	"go.uber.org/zap"
	locker "locker/api/go"
	"locker/internal/config"
	"locker/pkg/controller"
)

// LockerServiceServerImpl implements the generated LockerServiceServer interface.
type LockerServiceServerImpl struct {
	// Embed the generated UnimplementedLockerServiceServer for forward compatibility.
	locker.UnimplementedLockerServiceServer

	// Use your LockController to handle operations.
	ctrl controller.LockController
	log  *zap.Logger
}

// NewLockerServiceServer creates a new gRPC server implementation,
// wiring in the controller.
func NewLockerServiceServer(cfg *config.Configuration, log *zap.Logger) locker.LockerServiceServer {
	ctrl := controller.NewLockController(cfg, log)
	err := ctrl.Start()
	if err != nil {
		log.Fatal("failed to start controller", zap.Error(err))
	}
	return &LockerServiceServerImpl{
		ctrl: ctrl,
		log:  log,
	}
}

// Implement the Lock RPC. You can access fields from req to perform the actual lock.
func (s *LockerServiceServerImpl) Lock(ctx context.Context, req *locker.LockRequest) (*locker.LockResponse, error) {
	s.log.Info("Lock RPC called", zap.Any("request", req))
	// Here, you might call a method on your controller that performs locking.
	// For demonstration, we print the request fields.
	// In a real implementation, you’d use req.User, req.Uid, etc.
	// For example, you might call: s.ctrl.Lock(req)
	return &locker.LockResponse{
		Success: true,
		Message: "Lock acquired",
	}, nil
}

// Implement Unlock. No parameters, so simply forward the call.
func (s *LockerServiceServerImpl) Unlock(ctx context.Context, req *locker.UnlockRequest) (*locker.UnlockResponse, error) {
	s.log.Info("Unlock RPC called")
	// Call your controller's unlock functionality if needed.
	return &locker.UnlockResponse{
		Success: true,
		Message: "Lock released",
	}, nil
}

// Implement Status. Return current status and details.
func (s *LockerServiceServerImpl) Status(ctx context.Context, req *locker.StatusRequest) (*locker.StatusResponse, error) {
	s.log.Info("Status RPC called")
	// For demonstration, return a dummy status.
	// In a real implementation, retrieve these details from your controller.
	return &locker.StatusResponse{
		State:           locker.StatusResponse_LOCKED, // or UNLOCKED, etc.
		User:            "johndoe",
		Uid:             "1001",
		Tty:             "/dev/tty1",
		SessionId:       "session123",
		AllowedUsers:    []string{"johndoe", "janedoe"},
		AllowedGroups:   []string{"admin", "users"},
		Reason:          "Routine maintenance",
		UnlockOnExit:    true,
		UnlockTime:      nil, // or set a google.protobuf.Duration value
		UnlockAfterIdle: nil,
	}, nil
}

// Implement Authorize. Check if a given username is authorized.
func (s *LockerServiceServerImpl) Authorize(ctx context.Context, req *locker.AuthorizeRequest) (*locker.AuthorizeResponse, error) {
	s.log.Info("Authorize RPC called", zap.String("username", req.Username))
	// Insert your authorization logic here. For demo, we'll always authorize.
	authorized := true
	reason := ""
	if !authorized {
		reason = "User not allowed"
	}
	return &locker.AuthorizeResponse{
		Authorized: authorized,
		Reason:     reason,
	}, nil
}
