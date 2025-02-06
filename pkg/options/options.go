package options

import "fmt"

// LockOptions holds the configurable options for the lock.
type LockOptions struct {
	AutoUnlock    bool   // Whether to auto-unlock on exit.
	TimeUnlock    string // A date + time string for unlocking.
	UsersAllowed  string // Comma-separated list of allowed users.
	GroupsAllowed string // Comma-separated list of allowed groups.
	Reason        string // Optional reason for locking.
	Email         string // Optional email for notifications or info.
}

// Option represents a configuration function for LockOptions.
type Option func(*LockOptions)

// WithAutoUnlock configures the AutoUnlock option.
func WithAutoUnlock(auto bool) Option {
	return func(lo *LockOptions) {
		lo.AutoUnlock = auto
	}
}

// WithTimeUnlock configures the TimeUnlock option.
func WithTimeUnlock(timeUnlock string) Option {
	return func(lo *LockOptions) {
		lo.TimeUnlock = timeUnlock
	}
}

// WithUsersAllowed configures the UsersAllowed option.
func WithUsersAllowed(users string) Option {
	return func(lo *LockOptions) {
		lo.UsersAllowed = users
	}
}

// WithGroupsAllowed configures the GroupsAllowed option.
func WithGroupsAllowed(groups string) Option {
	return func(lo *LockOptions) {
		lo.GroupsAllowed = groups
	}
}

// WithReason configures the Reason option.
func WithReason(reason string) Option {
	return func(lo *LockOptions) {
		lo.Reason = reason
	}
}

// WithEmail configures the Email option.
func WithEmail(email string) Option {
	return func(lo *LockOptions) {
		lo.Email = email
	}
}

// NewLockOptions creates a new LockOptions instance with defaults and applies any provided options.
func NewLockOptions(opts ...Option) *LockOptions {
	// Set default values.
	lo := &LockOptions{
		AutoUnlock:    true, // Default to true.
		TimeUnlock:    "",
		UsersAllowed:  "",
		GroupsAllowed: "",
		Reason:        "",
		Email:         "",
	}

	// Apply each option.
	for _, opt := range opts {
		opt(lo)
	}

	return lo
}

// For demonstration, String prints the options.
func (lo LockOptions) String() string {
	return fmt.Sprintf(`LockOptions:
  AutoUnlock:    %t
  TimeUnlock:    %q
  UsersAllowed:  %q
  GroupsAllowed: %q
  Reason:        %q
  Email:         %q`,
		lo.AutoUnlock,
		lo.TimeUnlock,
		lo.UsersAllowed,
		lo.GroupsAllowed,
		lo.Reason,
		lo.Email,
	)
}
