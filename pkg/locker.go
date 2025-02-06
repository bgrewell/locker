package pkg

import (
	"locker/internal"
	"locker/internal/lock"
	"locker/pkg/options"
	"os/user"
)

const (
	ACTION_LOCK      = "lock"
	ACTION_UNLOCK    = "unlock"
	ACTION_AUTHORIZE = "authorize"
	ACTION_STATUS    = "status"
)

func NewLocker(options ...options.Option) Locker {
	return internal.NewUserLocker(options...)
}

type Locker interface {
	Lock() (err error)
	Unlock() (err error)
	CheckUnlockCriteria() (shouldUnlock bool, err error)
	AuthorizeLogin() (loginAllowed bool, err error)
}

func UserHasLockFile() bool {
	// Get the current user
	cu, err := user.Current()
	if err != nil {
		panic(err)
	}

	// Check if the user has a lock file
	lockfiles, err := lock.FindAllLockfiles()
	if err != nil {
		panic(err)
	}

	for _, lockfile := range lockfiles {
		if lockfile.User == cu.Username {
			return true
		}
	}

	return false
}
