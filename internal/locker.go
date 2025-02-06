package internal

import (
	"fmt"
	"locker/internal/lock"
	"locker/internal/pam"
	"locker/pkg/options"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

func NewUserLocker(opts ...options.Option) *UserLocker {
	u, err := user.Current()
	if err != nil {
		panic(fmt.Errorf("failed to get current user: %w", err))
	}
	return &UserLocker{
		User:    u,
		Options: options.NewLockOptions(opts...),
	}
}

type UserLocker struct {
	User    *user.User
	Options *options.LockOptions
}

func (s UserLocker) Lock() (err error) {

	// Construct the path to the lockfile.
	lockfilePath := filepath.Join(s.User.HomeDir, ".system.lockfile")

	// Check if there is already a lockfile
	if _, err := os.Stat(lockfilePath); err == nil {
		return os.ErrExist
	}

	// Split the allowed users and groups
	allowedUsers := strings.Split(s.Options.UsersAllowed, ",")
	allowedGroups := strings.Split(s.Options.GroupsAllowed, ",")

	// Calculate the unlock time if it is set
	var unlockTime time.Time
	if s.Options.TimeUnlock != "" {
		dur, err := time.ParseDuration(s.Options.TimeUnlock)
		if err != nil {
			return fmt.Errorf("failed to calculate unlock time: %w", err)
		}
		unlockTime = time.Now().Add(dur)
	}

	// Create a new lockfile
	lockfile, err := lock.NewLockFile(s.Options.Reason, unlockTime, s.Options.AutoUnlock, allowedUsers, allowedGroups)
	if err != nil {
		return fmt.Errorf("failed to lock: %w", err)
	}

	// Write the lockfile to disk
	if err := lockfile.WriteLockfile(lockfilePath); err != nil {
		return fmt.Errorf("failed to lock: %w", err)
	}

	return nil
}

func (s UserLocker) Unlock() (err error) {
	// Construct the path to the lockfile.
	lockfilePath := filepath.Join(s.User.HomeDir, ".system.lockfile")

	// Check if the file exists.
	if _, err := os.Stat(lockfilePath); err != nil {
		return fmt.Errorf("failed to unlock: %w", err)
	}

	// Remove the lockfile
	return os.Remove(lockfilePath)
}

func (s UserLocker) CheckUnlockCriteria() (shouldUnlock bool, err error) {
	// Get all lock files
	lockfiles, err := lock.FindAllLockfiles()
	if err != nil {
		return false, fmt.Errorf("failed to check unlock criteria: %w", err)
	}

	for _, lockfile := range lockfiles {
		homedir, err := findUserHomeDir(lockfile.User)
		if err != nil {
			return false, fmt.Errorf("failed to check unlock criteria: %w", err)
		}
		lockfileLocation := filepath.Join(homedir, ".system.lockfile")

		// Check to see if the file should be unlocked based on auto-unlock
		if lockfile.IsSessionFinished() {
			err = lockfile.Remove(lockfileLocation)
			if err != nil {
				return false, fmt.Errorf("failed to remove lock for user %s: %w", lockfile.User, err)
			}
			continue
		}

		// Check to see if the file should be unlocked based on the time
		if lockfile.IsExpired() {
			lockfile.Remove(lockfileLocation)
			if err != nil {
				return false, fmt.Errorf("failed to remove lock for user %s: %w", lockfile.User, err)
			}
			continue
		}

		// Otherwise it must still be active and desired so we return false
		return false, nil
	}

	return true, nil
}

func (s UserLocker) AuthorizeLogin() (loginAllowed bool, err error) {
	// Check login criteria first
	shouldUnlock, err := s.CheckUnlockCriteria()
	if err != nil {
		return false, fmt.Errorf("failed to authorize login: %w", err)
	}

	if shouldUnlock {
		return true, nil
	}

	// Otherwise check if the user can access the system either by being the user with the lock or in the
	// allowed users or groups. We get the information from the PAM environment variables.
	env := pam.GetEnvironment()
	// TODO: Temporary for testing
	if lockfiles, _ := lock.FindAllLockfiles(); len(lockfiles) > 0 {
		for _, lockfile := range lockfiles {
			if lockfile.User == *env.User {
				return true, nil
			}
		}
	}
	return false, nil
}

func findUserHomeDir(username string) (string, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return "", fmt.Errorf("failed to find user home dir: %w", err)
	}
	return u.HomeDir, nil
}
