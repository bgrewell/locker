package server

import (
	"context"
	"go.uber.org/zap"
	locker "locker/api/go"
	"locker/internal/config"
	"locker/internal/lock"
	"locker/pkg/controller"
	"time"
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

	lf := &lock.LockFile{
		User:          req.User,
		UID:           int(req.Uid),
		Reason:        req.Reason,
		Email:         req.Email,
		Session:       req.SessionId,
		TTY:           req.Tty,
		UnlockTime:    time.Now().Add(req.UnlockTime.AsDuration()),
		UnlockOnExit:  req.UnlockOnExit,
		UnlockOnIdle:  req.UnlockAfterIdle.AsDuration(),
		AllowedUsers:  req.AllowedUsers,
		AllowedGroups: req.AllowedGroups,
	}

	err := s.ctrl.Lock(lf)
	if err != nil {
		return nil, err
	}

	return &locker.LockResponse{
		Success: true,
		Message: "Lock acquired",
	}, nil
}

// Implement Unlock. No parameters, so simply forward the call.
func (s *LockerServiceServerImpl) Unlock(ctx context.Context, req *locker.UnlockRequest) (*locker.UnlockResponse, error) {
	s.log.Info("Unlock RPC called")

	err := s.ctrl.Unlock()
	if err != nil {
		return nil, err
	}

	return &locker.UnlockResponse{
		Success: true,
		Message: "Lock released",
	}, nil
}

// Implement Status. Return current status and details.
func (s *LockerServiceServerImpl) Status(ctx context.Context, req *locker.StatusRequest) (*locker.StatusResponse, error) {
	s.log.Info("Status RPC called")

	lockStatus, lockfile, err := s.ctrl.Status()
	if err != nil {
		return nil, err
	}

	s.log.Info("Status RPC called", zap.Any("lockStatus", lockStatus), zap.Any("lockfile", lockfile))
	return &locker.StatusResponse{
		State:           locker.StatusResponse_LockState(lockStatus),
		User:            lockfile.User,
		Uid:             int32(lockfile.UID),
		Tty:             lockfile.TTY,
		SessionId:       lockfile.Session,
		AllowedUsers:    lockfile.AllowedUsers,
		AllowedGroups:   lockfile.AllowedGroups,
		Reason:          lockfile.Reason,
		UnlockOnExit:    lockfile.UnlockOnExit,
		UnlockTime:      nil, //TODO: convert unlock time to a duration based on the time from now till unlock
		UnlockAfterIdle: nil, // TODO: Populate with the UnlockAFterIdle value
	}, nil
}

// Implement Authorize. Check if a given username is authorized.
func (s *LockerServiceServerImpl) Authorize(ctx context.Context, req *locker.AuthorizeRequest) (*locker.AuthorizeResponse, error) {
	s.log.Info("Authorize RPC called", zap.String("username", req.Username))

	// TODO: Convert from grpc type to internal type
	// TODO: Send to controller

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

func (s *LockerServiceServerImpl) Email(ctx context.Context, req *locker.EmailRequest) (*locker.EmailResponse, error) {
	s.log.Info("Email RPC called", zap.String("username", req.Username))

	email, err := s.ctrl.GetUserEmail(req.Username)
	if err != nil {
		return &locker.EmailResponse{
			Email: "",
			Found: false,
		}, err
	}

	return &locker.EmailResponse{
		Email: email,
		Found: true,
	}, nil
}
