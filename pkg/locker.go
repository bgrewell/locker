package pkg

const (
	ACTION_LOCK      = "lock"
	ACTION_UNLOCK    = "unlock"
	ACTION_AUTHORIZE = "authorize"
	ACTION_STATUS    = "status"
)

//func NewLocker(options ...options.Option) Locker {
//	return internal.NewUserLocker(options...)
//}
//
//// Locker is an interface that defines the methods that a locker must implement
//type Locker interface {
//	// Lock creates a system lock for the user that calls it
//	Lock() (err error)
//	// Unlock removes the system lock for the user that calls it
//	Unlock() (err error)
//	// CheckUnlockCriteria checks auto unlock and time unlock criteria and returns if the system should unlock
//	CheckUnlockCriteria() (shouldUnlock bool, err error)
//	// AuthorizeLogin checks if the user is allowed to log into the system
//	AuthorizeLogin() (loginAllowed bool, err error)
//}
//
//// UserHasLockFile is a helper function to determine if the current user has a lock file
//func UserHasLockFile(logger *logrus.Logger) bool {
//	// Get the current user
//	cu, err := user.Current()
//	if err != nil {
//		panic(err)
//	}
//
//	// Check if the user has a lock file
//	lockfiles, err := lock.FindAllLockfiles(logger)
//	if err != nil {
//		panic(err)
//	}
//
//	for _, lockfile := range lockfiles {
//		if lockfile.User == cu.Username {
//			return true
//		}
//	}
//
//	return false
//}
