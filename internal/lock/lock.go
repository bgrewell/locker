package lock

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"
)

type LockStatus int

const (
	StatusUnknown LockStatus = iota
	StatusLocked
	StatusUnlocked
)

// LockFile represents a lockfile on disk
type LockFile struct {
	User          string        `json:"user" yaml:"user"`
	UID           int           `json:"uid" yaml:"uid"`
	Reason        string        `json:"reason,omitempty" yaml:"reason,omitempty"`
	Email         string        `json:"email,omitempty" yaml:"email,omitempty"`
	Session       string        `json:"session,omitempty" yaml:"session,omitempty"`
	TTY           string        `json:"tty,omitempty" yaml:"tty,omitempty"`
	UnlockTime    time.Time     `json:"unlock_time,omitempty" yaml:"unlock_time,omitempty"`
	UnlockOnExit  bool          `json:"unlock_on_exit" yaml:"unlock_on_exit"`
	UnlockOnIdle  time.Duration `json:"unlock_on_idle,omitempty" yaml:"unlock_on_idle,omitempty"`
	AllowedUsers  []string      `json:"allowed_users,omitempty" yaml:"allowed_users,omitempty"`
	AllowedGroups []string      `json:"allowed_groups,omitempty" yaml:"allowed_groups,omitempty"`
	logger        *zap.Logger
}

// ReadLockfile reads a lockfile from disk and parses it into a LockFile struct
func ReadLockfile(path string, logger *zap.Logger) (*LockFile, error) {
	// Lockfiles are written in json format so just unmarshal it
	lockfile := &LockFile{
		logger: logger,
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open lockfile: %w", err)
	}
	defer f.Close()

	// Unmarshal the file
	err = json.NewDecoder(f).Decode(lockfile)
	if err != nil {
		return nil, err
	}

	return lockfile, nil
}

// WriteLockfile writes a lockfile to disk
func (lf LockFile) WriteLockfile(path string) error {
	// Open the file for writing
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create lockfile: %w", err)
	}

	// Marshal the lockfile to json
	data, err := json.MarshalIndent(lf, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal lockfile: %w", err)
	}

	// Write the data to the file
	_, err = f.Write(data)
	return err
}

// RemoveLockfile removes a lockfile from disk
func (lf LockFile) RemoveLockfile(lockfilePath string) error {
	if _, err := os.Stat(lockfilePath); err != nil {
		lf.logger.Warn("failed to remove lockfile. lockfile doesn't exist",
			zap.Error(err),
			zap.String("lockfile", lockfilePath))
		return fmt.Errorf("failed to stat lockfile: %w", err)
	}
	// RemoveLockfile the file.
	lf.logger.Debug("removing lockfile", zap.String("lockfile", lockfilePath))
	return os.Remove(lockfilePath)
}

// String implements the fmt.Stringer interface to provide a pretty-printed output.
func (lf LockFile) String() string {
	var unlockTimeStr string
	if !lf.UnlockTime.IsZero() {
		// Format the unlock time using a readable layout.
		unlockTimeStr = lf.UnlockTime.Format("2006-01-02 15:04:05 MST")
	} else {
		unlockTimeStr = "N/A"
	}

	var email string
	if lf.Email == "" {
		email = "N/A"
	} else {
		email = lf.Email
	}

	// Join slices into comma-separated strings.
	allowedUsers := strings.Join(lf.AllowedUsers, ", ")
	allowedGroups := strings.Join(lf.AllowedGroups, ", ")
	if allowedUsers == "" {
		allowedUsers = "-"
	}
	if allowedGroups == "" {
		allowedGroups = "-"
	}

	return fmt.Sprintf(`LockFile:
  User:          %s
  UID:           %d
  Reason:        %s
  Email:         %s
  Session:       %s
  TTY:           %s
  UnlockTime:    %s
  UnlockOnExit:  %t
  UnlockOnIdle:  %s
  AllowedUsers:  %s
  AllowedGroups: %s`,
		lf.User,
		lf.UID,
		lf.Reason,
		email,
		lf.Session,
		lf.TTY,
		unlockTimeStr,
		lf.UnlockOnExit,
		lf.UnlockOnIdle,
		allowedUsers,
		allowedGroups,
	)
}
