package internal

//import (
//	"fmt"
//	"github.com/sirupsen/logrus"
//	"locker/internal/lock"
//	"locker/internal/pam"
//	"locker/pkg/options"
//	"os"
//	"os/user"
//	"path/filepath"
//	"strings"
//	"time"
//)
//
//func NewUserLocker(opts ...options.Option) *UserLocker {
//
//	o := options.NewLockOptions(opts...)
//
//	u, err := user.Current()
//	if err != nil {
//		o.Logger.Errorf("failed to get current user: %v", err)
//	}
//	return &UserLocker{
//		User:    u,
//		Options: o,
//		Logger:  o.Logger,
//	}
//}
//
//type UserLocker struct {
//	User    *user.User
//	Options *options.LockOptions
//	Logger  *logrus.Logger
//}
//
//func (s UserLocker) Lock() (err error) {
//
//	// Construct the path to the lockfile.
//	lockfilePath := filepath.Join(s.User.HomeDir, ".system.lockfile")
//
//	// Check if there is already a lockfile
//	if _, err := os.Stat(lockfilePath); err == nil {
//		s.Logger.Errorf("lockfile already exists at %s", lockfilePath)
//		return os.ErrExist
//	}
//
//	// Split the allowed users and groups
//	allowedUsers := strings.Split(s.Options.UsersAllowed, ",")
//	allowedGroups := strings.Split(s.Options.GroupsAllowed, ",")
//	s.Logger.WithFields(logrus.Fields{
//		"allowed_users":  allowedUsers,
//		"allowed_groups": allowedGroups,
//	}).Tracef("allowed users and groups")
//
//	// Calculate the unlock time if it is set
//	var unlockTime time.Time
//	if s.Options.TimeUnlock != "" {
//		dur, err := time.ParseDuration(s.Options.TimeUnlock)
//		if err != nil {
//			s.Logger.Errorf("failed to parse unlock time: %v", err)
//			return fmt.Errorf("failed to calculate unlock time: %w", err)
//		}
//		unlockTime = time.Now().Add(dur)
//		s.Logger.Debugf("unlock time set to %s", unlockTime)
//	}
//
//	// Create a new lockfile
//	lockfile, err := lock.NewLockFile(s.Options.Reason, unlockTime, s.Options.AutoUnlock, allowedUsers, allowedGroups, s.Logger)
//	if err != nil {
//		s.Logger.Errorf("failed to create lockfile: %v", err)
//		return fmt.Errorf("failed to lock: %w", err)
//	}
//
//	// Write the lockfile to disk
//	if err := lockfile.WriteLockfile(lockfilePath); err != nil {
//		s.Logger.Errorf("failed to write lockfile: %v", err)
//		return fmt.Errorf("failed to lock: %w", err)
//	}
//
//	s.Logger.WithFields(logrus.Fields{
//		"user":     s.User.Username,
//		"homedir":  s.User.HomeDir,
//		"lockfile": lockfilePath,
//	}).Debug("lockfile created")
//	return nil
//}
//
//func (s UserLocker) Unlock() (err error) {
//	// Construct the path to the lockfile.
//	lockfilePath := filepath.Join(s.User.HomeDir, ".system.lockfile")
//	s.Logger.Debugf("unlocking lockfile at %s", lockfilePath)
//
//	// Check if the file exists.
//	if _, err := os.Stat(lockfilePath); err != nil {
//		s.Logger.Warnf("lockfile does not exist at %s", lockfilePath)
//		return fmt.Errorf("failed to unlock: %w", err)
//	}
//
//	// RemoveLockfile the lockfile
//	s.Logger.Debugf("removing lockfile at %s", lockfilePath)
//	return os.Remove(lockfilePath)
//}
//
//func (s UserLocker) CheckUnlockCriteria() (shouldUnlock bool, err error) {
//	// Get all lock files
//	lockfiles, err := lock.FindAllLockfiles(s.Logger)
//	s.Logger.WithFields(logrus.Fields{
//		"count": len(lockfiles),
//	}).Debug("found lockfiles")
//	if err != nil {
//		s.Logger.WithFields(logrus.Fields{
//			"error": err,
//		}).Debug("failed to find lockfiles")
//		return false, fmt.Errorf("failed to check unlock criteria: %w", err)
//	}
//
//	for _, lockfile := range lockfiles {
//		homedir, err := findUserHomeDir(lockfile.User)
//		if err != nil {
//			s.Logger.WithFields(logrus.Fields{"homedir": homedir, "error": err}).Debug("failed to find user home dir")
//			return false, fmt.Errorf("failed to check unlock criteria: %w", err)
//		}
//		lockfileLocation := filepath.Join(homedir, ".system.lockfile")
//		s.Logger.WithFields(logrus.Fields{"lockfile": lockfileLocation}).Debug("checking lockfile")
//
//		// Check to see if the file should be unlocked based on auto-unlock
//		if lockfile.IsSessionFinished() {
//			s.Logger.WithFields(logrus.Fields{
//				"user":    lockfile.User,
//				"session": lockfile.Session,
//				"trigger": "session finished",
//			}).Debug("removing lock file")
//			err = lockfile.Remove(lockfileLocation)
//			if err != nil {
//				s.Logger.WithFields(logrus.Fields{"error": err}).Debug("failed to remove lock file")
//				return false, fmt.Errorf("failed to remove lock for user %s: %w", lockfile.User, err)
//			}
//			continue
//		}
//
//		// Check to see if the file should be unlocked based on the time
//		if lockfile.IsExpired() {
//			s.Logger.WithFields(logrus.Fields{
//				"user":    lockfile.User,
//				"session": lockfile.Session,
//				"trigger": "lock expired",
//			}).Debug("removing lock file")
//			lockfile.Remove(lockfileLocation)
//			if err != nil {
//				s.Logger.WithFields(logrus.Fields{"error": err}).Debug("failed to remove lock file")
//				return false, fmt.Errorf("failed to remove lock for user %s: %w", lockfile.User, err)
//			}
//			continue
//		}
//
//		// Otherwise it must still be active and desired so we return false
//		return false, nil
//	}
//
//	return true, nil
//}
//
//func (s UserLocker) AuthorizeLogin() (loginAllowed bool, err error) {
//	env := pam.GetEnvironment()
//	fields := logrus.Fields{
//		"user":     env.User,
//		"service":  env.Service,
//		"tty":      env.TTY,
//		"rhost":    env.RHost,
//		"ruser":    env.RUser,
//		"pam_type": env.PAMType,
//	}
//
//	// Check login criteria first
//	s.Logger.WithFields(fields).Debug("authorizing user login")
//	shouldUnlock, err := s.CheckUnlockCriteria()
//	if err != nil {
//		s.Logger.WithFields(logrus.Fields{"error": err}).Debug("check unlock criteria failed")
//		return false, fmt.Errorf("failed to authorize login: %w", err)
//	}
//
//	if shouldUnlock {
//		fields["unlock"] = true
//		s.Logger.WithFields(fields).Debug("login authorized")
//		return true, nil
//	}
//
//	// Otherwise check if the user can access the system either by being the user with the lock or in the
//	// allowed users or groups. We get the information from the PAM environment variables.
//	// TODO: Temporary for testing
//	if lockfiles, _ := lock.FindAllLockfiles(s.Logger); len(lockfiles) > 0 {
//		for _, lockfile := range lockfiles {
//			if lockfile.User == env.User {
//				fields["lock_user"] = lockfile.User
//				fields["lock_session"] = lockfile.Session
//				fields["lock_tty"] = lockfile.TTY
//				s.Logger.WithFields(fields).Debug("login authorized")
//				return true, nil
//			}
//		}
//	}
//
//	//// Attempt to open /dev/tty for writing
//	//tty, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0)
//	//if err != nil {
//	//	// If there's no TTY (non-interactive SSH?), silently exit
//	//	s.Logger.WithFields(logrus.Fields{
//	//		"error": err,
//	//	}).Debug("failed to open TTY")
//	//}
//	//defer tty.Close()
//	//
//	//// Print message to the TTY
//	//fmt.Fprintln(tty, "hello")
//
//	fmt.Printf("You are not authorized to access this system.\n")
//
//	s.Logger.WithFields(fields).Debug("access denied")
//	return false, nil
//}
//
//// findUserHomeDir finds the home directory for a given user
//func findUserHomeDir(username string) (string, error) {
//	u, err := user.Lookup(username)
//	if err != nil {
//		return "", fmt.Errorf("failed to find user home dir: %w", err)
//	}
//	return u.HomeDir, nil
//}
