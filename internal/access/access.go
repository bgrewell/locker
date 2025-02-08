package access

// CheckAccess is used to determine if a user should be granted acccess to the system. It checks to see if the system
// is locked and if it is not then it returns true (allows the normal authorization process to continue). If the system
// is locked then it will check the access rules and return a result based on the rules.
func CheckAccess(username string) (approved bool, reason string) {

	denyReason := `
┌────────────────────────────────────────────────────────────────────┐
│                        SYSTEM LOCK NOTICE                          │
├────────────────────────────────────────────────────────────────────┤
│ This system is currently locked by 'ben'.                          │
│                                                                    │
│ Access is restricted to the following:                             │
│                                                                    │
│    Users:                                                          │
│       - root                                                       │
│       - admin                                                      │
│       - ben                                                        │
│                                                                    │
│    Groups:                                                         │
│       - wheel                                                      │
│                                                                    │
│ The system is scheduled to automatically unlock at:                │
│    2021-12-31 23:59:59                                             │
│                                                                    │
│ Please contact your administrator for further details.             │
└────────────────────────────────────────────────────────────────────┘`

	if username == "ben" {
		return true, ""
	}

	return false, denyReason
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
