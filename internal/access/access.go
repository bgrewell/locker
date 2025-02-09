package access

import (
	"bytes"
	"fmt"
	"locker/internal/config"
	"locker/internal/lock"
	"os"
	"os/user"
	"text/template"
	"time"
)

// LockNoticeData holds the dynamic data for the lock notice.
type LockNoticeData struct {
	LockedBy      string
	AllowedUsers  []string
	AllowedGroups []string
	// UnlockTime is optional; if nil, the auto-unlock section is omitted.
	UnlockTime *time.Time
	// Email is optional; if provided, it will be used as the contact, otherwise LockedBy is used.
	Email string
	// Reason is optional; if provided, it will be included in the notice.
	Reason string
}

// generateDenyReason produces a lock notice message using template variables.
func generateDenyReason(data LockNoticeData) (string, error) {
	const tmpl = `
┌────────────────────────────────────────────────────────────────────┐
│                        SYSTEM LOCK NOTICE                          │
├────────────────────────────────────────────────────────────────────┤
│ This system is currently locked by '{{.LockedBy}}'.                
{{- if .Reason }}
│ Reason: {{.Reason}}
{{- end }}
{{- if or (gt (len .AllowedUsers) 0) (gt (len .AllowedGroups) 0) }}
│                                                                    │
│ Access is restricted to the following:
{{- if gt (len .AllowedUsers) 0 }}
│    Users:
{{- range .AllowedUsers }}
│       - {{.}}
{{- end }}
{{- end }}
{{- if gt (len .AllowedGroups) 0 }}
│                                                                    │
│    Groups:
{{- range .AllowedGroups }}
│       - {{.}}
{{- end }}
{{- end }}
{{- end }}
{{- if .UnlockTime }}
│                                                                    │
│ The system is scheduled to automatically unlock at:                
│    {{.UnlockTime.Format "2006-01-02 15:04:05"}}                                            
{{- end }}
│                                                                    │
│ Please contact {{if .Email}}{{.Email}}{{else}}{{.LockedBy}}{{end}} for further details.             
└────────────────────────────────────────────────────────────────────┘
`
	t, err := template.New("denyNotice").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.String(), nil
}

// CheckAccess is used to determine if a user should be granted acccess to the system. It checks to see if the system
// is locked and if it is not then it returns true (allows the normal authorization process to continue). If the system
// is locked then it will check the access rules and return a result based on the rules.
func CheckAccess(username string) (approved bool, reason string) {

	//	denyReason := `
	//┌────────────────────────────────────────────────────────────────────┐
	//│                        SYSTEM LOCK NOTICE                          │
	//├────────────────────────────────────────────────────────────────────┤
	//│ This system is currently locked by 'ben'.                          │
	//│                                                                    │
	//│ Access is restricted to the following:                             │
	//│                                                                    │
	//│    Users:                                                          │
	//│       - root                                                       │
	//│       - admin                                                      │
	//│       - ben                                                        │
	//│                                                                    │
	//│    Groups:                                                         │
	//│       - wheel                                                      │
	//│                                                                    │
	//│ The system is scheduled to automatically unlock at:                │
	//│    2021-12-31 23:59:59                                             │
	//│                                                                    │
	//│ Please contact your administrator for further details.             │
	//└────────────────────────────────────────────────────────────────────┘`

	// Pull the lockfile data
	var lf *lock.LockFile
	var err error
	cfg, _ := config.ReadConfiguration()
	if cfg != nil {
		lf, err = lock.ReadLockfile(cfg.LockFileLocation, nil)
		if err == os.ErrNotExist {
			// If there is no lockfile then we can allow access
			return true, ""
		}
		if err != nil {
			return cfg.FailOpen, fmt.Sprintf("failed to read lockfile: %v", err) //TODO: Add a default message
		}
	}

	// Otherwise if the lockfile exists we need to see if the user is allowed to access the system
	if lf.User == username {
		return true, ""
	}

	// Check if the username is in the allowed users list
	for _, user := range lf.AllowedUsers {
		if user == username {
			return true, ""
		}
	}

	// Check if the user is a member of the allowed groups
	for _, group := range lf.AllowedGroups {
		present, _ := UserInGroup(username, group)
		if present {
			return true, ""
		}
	}

	// Build the template data structure
	data := LockNoticeData{
		LockedBy:      lf.User,
		AllowedUsers:  lf.AllowedUsers,
		AllowedGroups: lf.AllowedGroups,
		UnlockTime:    &lf.UnlockTime,
		Email:         lf.Email,
		Reason:        lf.Reason,
	}

	// Generate the deny reason
	notice, err := generateDenyReason(data)
	if err != nil {
		return false, fmt.Sprintf("failed to generate deny reason: %v", err) //TODO: Add a default message if we failed to generate a dynamic message
	}

	return false, notice
}

// CheckWarning is used to return a warning message to users that are authorized to use the system during the lock
// to ensure that they are aware the system is currently locked.
func CheckWarning() string {
	// TODO: This should be templated and customizable by users.
	return "\033[1;31m┌────────────────────────────────────────────────────────────────────┐\033[0m\n" +
		"\033[1;31m│\033[0m\033[1;33m                        ADMIN WARNING                               \033[0m\033[1;31m│\033[0m\n" +
		"\033[1;31m├────────────────────────────────────────────────────────────────────┤\033[0m\n" +
		"\033[1;31m│\033[0m\033[1;37m WARNING: This system is currently locked by 'ben'.                 \033[0m\033[1;31m│\033[0m\n" +
		"\033[1;31m│\033[0m\033[1;37m Please be advised that the system is under a temporary lock.       \033[0m\033[1;31m│\033[0m\n" +
		"\033[1;31m│\033[0m\033[1;37m All critical operations and changes are restricted until the       \033[0m\033[1;31m│\033[0m\n" +
		"\033[1;31m│\033[0m\033[1;37m system is unlocked.                                                \033[0m\033[1;31m│\033[0m\n" +
		"\033[1;31m│\033[0m\033[1;37m Locked by: \033[1;32mben\033[0m\033[1;37m                                                     \033[1;31m│\033[0m\n" +
		"\033[1;31m│\033[0m\033[1;37m Unlock scheduled at: \033[1;32m2021-12-31 23:59:59\033[0m\033[1;37m                           \033[1;31m│\033[0m\n" +
		"\033[1;31m│\033[0m\033[1;37m Proceed with extreme caution and contact the locking authority if  \033[0m\033[1;31m│\033[0m\n" +
		"\033[1;31m│\033[0m\033[1;37m you believe immediate action is required.                          \033[0m\033[1;31m│\033[0m\n" +
		"\033[1;31m└────────────────────────────────────────────────────────────────────┘\033[0m"
}

// UserInGroup returns true if the given user is a member of the specified group.
func UserInGroup(userName, groupName string) (bool, error) {
	u, err := user.Lookup(userName)
	if err != nil {
		return false, err
	}

	// Get the list of group IDs for the user.
	groupIDs, err := u.GroupIds()
	if err != nil {
		return false, err
	}

	// Iterate over the user's groups.
	for _, gid := range groupIDs {
		g, err := user.LookupGroupId(gid)
		if err != nil {
			continue // skip groups we can't look up
		}
		if g.Name == groupName {
			return true, nil
		}
	}

	return false, nil
}
