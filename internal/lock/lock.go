package lock

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ReadLockfile reads a lockfile from disk and parses it into a LockFile struct
func ReadLockfile(path string) (*LockFile, error) {
	// Lockfiles are written in json format so just unmarshal it
	lockfile := &LockFile{}
	err := readJSONFile(path, lockfile)
	if err != nil {
		return nil, err
	}

	return lockfile, nil
}

func NewLockFile(reason string, unlockTime time.Time, unlockOnExit bool, allowedUsers, allowedGroups []string) (lockfile *LockFile, err error) {

	// Use loginctl to get session information.
	sessionID, tty, userName, uid, err := getSessionStatus()
	if err != nil {
		// If loginctl fails, fall back to os/user.
		cu, err2 := user.Current()
		if err2 != nil {
			return nil, fmt.Errorf("failed to get current user: %w (and loginctl error: %v)", err2, err)
		}
		userName = cu.Username
		uid, _ = strconv.Atoi(cu.Uid)
		tty = os.Getenv("SSH_TTY")
		sessionID = ""
	}

	email, err := getUserEmail(userName)
	if err != nil {
		email = ""
	}

	return &LockFile{
		User:          userName,
		UID:           uid,
		Email:         email,
		Reason:        reason,
		Session:       sessionID,
		TTY:           tty,
		UnlockTime:    unlockTime,
		UnlockOnExit:  unlockOnExit,
		AllowedUsers:  allowedUsers,
		AllowedGroups: allowedGroups,
	}, nil
}

type LockFile struct {
	User          string    `json:"user" yaml:"user"`
	UID           int       `json:"uid" yaml:"uid"`
	Reason        string    `json:"reason,omitempty" yaml:"reason,omitempty"`
	Email         string    `json:"email,omitempty" yaml:"email,omitempty"`
	Session       string    `json:"session,omitempty" yaml:"session,omitempty"`
	TTY           string    `json:"tty,omitempty" yaml:"tty,omitempty"`
	UnlockTime    time.Time `json:"unlock_time,omitempty" yaml:"unlock_time,omitempty"`
	UnlockOnExit  bool      `json:"unlock_on_exit" yaml:"unlock_on_exit"`
	AllowedUsers  []string  `json:"allowed_users,omitempty" yaml:"allowed_users,omitempty"`
	AllowedGroups []string  `json:"allowed_groups,omitempty" yaml:"allowed_groups,omitempty"`
}

func (lf LockFile) Remove(lockfilePath string) error {
	if _, err := os.Stat(lockfilePath); err != nil {
		return fmt.Errorf("failed to stat lockfile: %w", err)
	}
	// Remove the file.
	return os.Remove(lockfilePath)
}

func (lf LockFile) IsSessionFinished() bool {
	// If no session is recorded, assume the session is over. This shouldn't happen
	if lf.Session == "" {
		return true //TODO: align this with the config.FailOpen setting
	}

	cmd := exec.Command("loginctl", "list-sessions")
	output, err := cmd.Output()
	if err != nil {
		// If there's an error, we assume the session is over.
		return true //TODO: align this with the config.FailOpen setting
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
		if sessionID == lf.Session {
			// If the session is active, it's not over.
			if state == "active" {
				return false
			}
			// Otherwise, it is over.
			return true
		}
	}

	// If no matching session is found, assume it is over.
	return true
}

func (lf LockFile) IsExpired() bool {
	return !lf.UnlockTime.IsZero() && time.Now().After(lf.UnlockTime)
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
		allowedUsers,
		allowedGroups,
	)
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

// FindAllLockfiles searches through all users' home directories (as determined by /etc/passwd)
// for a file named ".system.lockfile". It returns a slice of LockFile pointers for any found.
func FindAllLockfiles() ([]*LockFile, error) {
	passwdPath := "/etc/passwd"
	file, err := os.Open(passwdPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", passwdPath, err)
	}
	defer file.Close()

	var lockfiles []*LockFile
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Each line in /etc/passwd is colon-delimited.
		// Format: username:password:UID:GID:GECOS:home_directory:shell
		line := scanner.Text()
		fields := strings.Split(line, ":")
		if len(fields) < 6 {
			continue // skip malformed lines
		}
		homeDir := fields[5]
		// Construct the path to the lockfile.
		lockfilePath := filepath.Join(homeDir, ".system.lockfile")
		// Check if the file exists.
		if _, err := os.Stat(lockfilePath); err == nil {
			lf, err := ReadLockfile(lockfilePath)
			if err != nil {
				// Log the error and continue. You might also want to collect these errors.
				fmt.Fprintf(os.Stderr, "warning: failed to read lockfile %s: %v\n", lockfilePath, err)
				continue
			}
			lockfiles = append(lockfiles, lf)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading %s: %w", passwdPath, err)
	}

	return lockfiles, nil
}

// readJSONFile reads a json file from disk and unmarshals it into the provided interface
func readJSONFile(path string, lockfile *LockFile) (err error) {
	// Read the file and unmarshal it into the lockfile struct
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open lockfile: %w", err)
	}
	defer f.Close()

	// Unmarshal the file
	return json.NewDecoder(f).Decode(lockfile)
}

// getUserEmail attempts to obtain the user's email by invoking dbus-send.
// It executes the command and parses its output for a valid email address.
func getUserEmail(userName string) (string, error) {
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

// getSessionStatus invokes "loginctl session-status --no-pager" and parses its output.
// It returns the session ID, TTY, username, and UID.
func getSessionStatus() (sessionID, tty, userName string, uid int, err error) {
	cmd := exec.Command("loginctl", "session-status", "--no-pager")
	output, err := cmd.Output()
	if err != nil {
		return "", "", "", 0, fmt.Errorf("failed to execute loginctl: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	firstLineFound := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// The first non-empty line should look like: "240 - ben (1001)"
		if !firstLineFound {
			firstLineFound = true
			// Split on the hyphen.
			parts := strings.SplitN(line, "-", 2)
			if len(parts) < 2 {
				return "", "", "", 0, errors.New("unexpected format in session-status output (first line)")
			}
			sessionID = strings.TrimSpace(parts[0])
			// The remainder should contain the username and UID.
			rest := strings.TrimSpace(parts[1]) // e.g., "ben (1001)"
			// Use a regular expression to extract the username and UID.
			re := regexp.MustCompile(`^(.+)\s+\((\d+)\)$`)
			matches := re.FindStringSubmatch(rest)
			if len(matches) != 3 {
				return "", "", "", 0, errors.New("could not parse username and UID from session-status output")
			}
			userName = strings.TrimSpace(matches[1])
			uid, err = strconv.Atoi(matches[2])
			if err != nil {
				return "", "", "", 0, fmt.Errorf("invalid UID in session-status output: %w", err)
			}
			continue
		}

		// Look for the TTY line.
		if strings.HasPrefix(line, "TTY:") {
			// Line example: "TTY: pts/4"
			tty = strings.TrimSpace(strings.TrimPrefix(line, "TTY:"))
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return "", "", "", 0, fmt.Errorf("error reading loginctl output: %w", err)
	}
	if sessionID == "" || tty == "" || userName == "" {
		return "", "", "", 0, errors.New("incomplete session-status information")
	}
	return sessionID, tty, userName, uid, nil
}
