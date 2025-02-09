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
	// TODO: Check if the lockfile exists when we start (i.e. after a reboot) and restore the lock
	return &LockerServiceServerImpl{
		ctrl: ctrl,
		log:  log,
	}
}

// Implement the Lock RPC. You can access fields from req to perform the actual lock.
func (s *LockerServiceServerImpl) Lock(ctx context.Context, req *locker.LockRequest) (*locker.LockResponse, error) {
	s.log.Info("Lock RPC called", zap.Any("request", req))

	unlockTime := time.Time{}
	if req.UnlockTime != nil && req.UnlockTime.AsDuration() > 0 {
		unlockTime = time.Now().Add(req.UnlockTime.AsDuration())
	}

	lf := &lock.LockFile{
		User:          req.User,
		UID:           int(req.Uid),
		Reason:        req.Reason,
		Email:         req.Email,
		Session:       req.SessionId,
		TTY:           req.Tty,
		UnlockTime:    unlockTime,
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
	if lockStatus == lock.StatusUnlocked {
		return &locker.StatusResponse{
			State: locker.StatusResponse_LockState(lockStatus),
		}, nil
	} else if lockStatus == lock.StatusLocked {
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
	} else {
		return &locker.StatusResponse{State: locker.StatusResponse_UNKNOWN}, nil
	}
}

// Implement Authorize. Check if a given username is authorized.
func (s *LockerServiceServerImpl) Authorize(ctx context.Context, req *locker.AuthorizeRequest) (*locker.AuthorizeResponse, error) {
	s.log.Info("Authorize RPC called", zap.String("username", req.Username))

	// TODO: This shouldn't be needed as the PAM module is the only thing that checks authorization and it does so
	//       by accessing the lockfile directly in order to improve responsiveness. This is just a placeholder for now.
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
