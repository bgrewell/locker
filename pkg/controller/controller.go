package controller

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
	"locker/internal/config"
	"locker/internal/lock"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func NewLockController(configuration *config.Configuration, log *zap.Logger) LockController {
	return &StandardLockController{
		configuration: configuration,
		log:           log,
		pollInterval:  1 * time.Second,
	}
}

type LockController interface {
	Start() error
	Stop() error
	Lock(lf *lock.LockFile) error
	Unlock() error
	Status() (status lock.LockStatus, lf *lock.LockFile, err error)
	Authorize() error
	GetUserEmail(username string) (string, error)
}

// Confirm that StandardLockController implements LockController
var _ LockController = &StandardLockController{}

// StandardLockController is the default implementation of the LockController interface.
type StandardLockController struct {
	running       bool
	log           *zap.Logger
	configuration *config.Configuration
	lockfile      *lock.LockFile
	pollInterval  time.Duration
}

// Start begins the lock controller's main loop.
func (lc *StandardLockController) Start() error {

	lc.running = true

	//TODO: Setup to watch the lockfile on disk so we can update our in-memory copy if changes are made
	go lc.watchLockfile()
	go lc.process()

	return nil
}

// Stop halts the lock controller's main loop.
func (lc *StandardLockController) Stop() error {

	lc.running = false
	return nil

}

// Lock writes the lockfile to disk and starts the lock controller.
func (lc *StandardLockController) Lock(lf *lock.LockFile) error {
	lc.lockfile = lf
	err := lf.WriteLockfile(lc.configuration.LockFileLocation)
	if err != nil {
		return err
	}
	return lc.Start()
}

// Unlock removes the lockfile from disk and stops the lock controller.
func (lc *StandardLockController) Unlock() error {
	err := lc.lockfile.RemoveLockfile(lc.configuration.LockFileLocation)
	if err != nil && err != os.ErrNotExist {
		return err
	}
	return lc.Stop()
}

// Status returns the current lock status and lockfile.
func (lc *StandardLockController) Status() (status lock.LockStatus, lf *lock.LockFile, err error) {
	if lc.lockfile == nil {
		return lock.StatusUnlocked, nil, nil
	}

	return lock.StatusLocked, lc.lockfile, nil
}

// Authorize checks if the user is authorized to unlock the system.
func (lc *StandardLockController) Authorize() error {
	return nil
}

// GetUserEmail attempts to obtain the user's email by invoking dbus-send.
// It executes the command and parses its output for a valid email address.
func (lc *StandardLockController) GetUserEmail(userName string) (string, error) {

	// Validate the username to allow only letters, digits, underscores, and hyphens.
	validUser := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validUser.MatchString(userName) {
		return "", errors.New("invalid username format")
	}

	cmd := exec.Command("sudo", "dbus-send",
		"--print-reply",
		"--system",
		"--dest=org.freedesktop.sssd.infopipe",
		"/org/freedesktop/sssd/infopipe",
		"org.freedesktop.sssd.infopipe.GetUserAttr",
		"string:"+userName,
		"array:string:mail",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run dbus-send command: %w; output: %s", err, output)
	}

	// Regex to find lines like: string "someone@example.com"
	re := regexp.MustCompile(`string "([^"]+)"`)
	matches := re.FindAllStringSubmatch(string(output), -1)
	if len(matches) == 0 {
		return "", fmt.Errorf("no email found in dbus-send output")
	}

	// Return the first match that isn’t the literal "mail".
	for _, m := range matches {
		if len(m) > 1 && strings.ToLower(m[1]) != "mail" {
			return m[1], nil
		}
	}

	return "", fmt.Errorf("no valid email found in dbus-send output")
}

// IsSessionFinished checks if the session is still active.
func (lc *StandardLockController) IsSessionFinished() bool {
	// If no session is recorded, assume the session is over. This shouldn't happen
	if lc.lockfile.Session == "" {
		if lc.configuration.FailOpen {
			lc.log.Debug("no session recorded. unlocking...",
				zap.Bool("fail_open", true))
			return true
		} else {
			lc.log.Warn("no session recorded. not unlocking...",
				zap.Bool("fail_open", false))
			return false
		}
	}

	// Get the session status.
	cmd := exec.Command("/usr/bin/loginctl", "list-sessions")
	output, err := cmd.Output()
	if err != nil {
		// If there's an error, we assume the session is over.
		lc.log.Debug("failed to list sessions. using fail_open condition...",
			zap.Error(err),
			zap.String("session", lc.lockfile.Session),
			zap.Bool("fail_open", lc.configuration.FailOpen))
		return lc.configuration.FailOpen
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		// Skip the header line.
		if strings.HasPrefix(line, "SESSION") {
			continue
		}

		// Split the line into fields.
		// Example line:
		// "240 1001 ben - pts/4 active no -"
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		sessionID := fields[0]
		state := fields[5] // Expecting the state to be in the sixth column.

		// If we find a matching session...
		if sessionID == lc.lockfile.Session {
			lc.log.Debug("found session in loginctl output",
				zap.String("session", lc.lockfile.Session),
				zap.String("session_id", sessionID),
				zap.String("state", state))
			// If the session is active, it's not over.
			if state == "active" {
				lc.log.Debug("session is active. not unlocking.")
				return false
			}
			// Otherwise, it is over.
			lc.log.Debug("session is not active. unlocking.")
			return true
		}
	}

	// If no matching session is found, assume it is over.
	lc.log.Debug("session not found in loginctl output. unlocking.")
	return true
}

// IdleExceeded checks if the idle time has been exceeded.
func (lc *StandardLockController) IdleExceeded() bool {
	//TODO: Implement idle time check
	return false
}

// ExitDetected checks if the exit condition has been met.
func (lc *StandardLockController) ExitDetected() bool {
	//TODO: Implement exit detection
	return false
}

// process is a goroutine that checks the lockfile for unlock conditions and unlocks the system if any are met.
func (lc *StandardLockController) process() {

	if lc.lockfile == nil {
		lc.Stop()
		return
	}

	// Exit early if all the conditional unlock options are disabled
	if lc.lockfile.UnlockOnExit == false &&
		lc.lockfile.UnlockTime.IsZero() &&
		lc.lockfile.UnlockOnIdle == 0 {

		// With no auto-unlock conditions there is no need to loop watching the lockfile
		// we just set running to false and exit
		lc.log.Debug("No unlock conditions are enabled, ending processing loop")
		lc.running = false
		return
	}

	for lc.running {

		if !lc.lockfile.UnlockTime.IsZero() && time.Now().After(lc.lockfile.UnlockTime) {
			lc.log.Info("Unlocking system due to unlock time condition")
			lc.Unlock()
		}

		if lc.lockfile.UnlockOnIdle != 0 && lc.IdleExceeded() {
			lc.log.Info("Unlocking system due to idle time condition")
			lc.Unlock()
		}

		if lc.lockfile.UnlockOnExit && lc.ExitDetected() {
			lc.log.Info("Unlocking system due to exit condition")
			lc.Unlock()
		}

		time.Sleep(lc.pollInterval)
	}
}

// watchLockfile is a goroutine that watches the lockfile for changes and updates the in-memory copy if there is a write
// event. If the lockfile is removed, then the in-memory copy is cleared and the running flag is set to false.
func (lc *StandardLockController) watchLockfile() {
	lockfile := lc.configuration.LockFileLocation

	// Create a new file watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		lc.log.Error("failed to create file watcher", zap.Error(err))
		return
	}
	defer watcher.Close()

	// Add the lockfile to the watcher.
	if err := watcher.Add(lockfile); err != nil {
		lc.log.Error("failed to add lockfile to watcher", zap.String("lockfile", lockfile), zap.Error(err))
		return
	}

	lc.log.Info("Started watching lockfile", zap.String("lockfile", lockfile))

	// Process events inline (assuming this function is already running in its own goroutine).
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				lc.log.Warn("watcher events channel closed")
				return
			}
			lc.log.Info("lockfile event", zap.Any("event", event))
			// For example, if the file is written to:
			if event.Op&fsnotify.Write == fsnotify.Write {
				lc.log.Info("lockfile was modified", zap.String("lockfile", lockfile))
				lf, err := lock.ReadLockfile(lockfile, lc.log)
				if err != nil {
					lc.log.Error("failed to read lockfile", zap.String("lockfile", lockfile), zap.Error(err))
				} else {
					lc.lockfile = lf
					lc.log.Info("lockfile updated", zap.Any("lockfile", lc.lockfile))
				}
			}
			// If the file is removed, then clear out the in-memory copy also
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				lc.log.Info("lockfile was removed", zap.String("lockfile", lockfile))
				lc.lockfile = nil
				lc.running = false
				return
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				lc.log.Warn("watcher errors channel closed")
				return
			}
			lc.log.Error("watcher error", zap.Error(err))
		}
	}
}
