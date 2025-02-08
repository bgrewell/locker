package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"locker/internal/lock"
	"locker/pkg"
	"locker/pkg/options"
	"log"
	"os"

	"github.com/bgrewell/usage"
)

var (
	version   string
	builddate string
	commit    string
	branch    string
)

func main() {

	// Setup application useage
	usageBuilder := usage.NewUsage(
		usage.WithApplicationName("locker"),
		usage.WithApplicationVersion(version),
		usage.WithApplicationBuildDate(builddate),
		usage.WithApplicationCommitHash(commit),
		usage.WithApplicationBranch(branch),
		usage.WithApplicationDescription("A utility to manage exclusive system access with custom access rules to allow access to specified users and/or groups."),
	)

	// Add flags and arguments
	help := usageBuilder.AddBooleanOption("h", "help", false, "Show this help message", "optional", nil)
	debug := usageBuilder.AddBooleanOption("d", "debug", false, "Enable debug output", "optional", nil)
	enable := usageBuilder.AddBooleanOption("e", "enable", false, "Enable locking on the system", "", nil)
	disable := usageBuilder.AddBooleanOption("D", "disable", false, "Disable locking on the system", "", nil)
	optgroup := usageBuilder.AddGroup(1, "Locking Options", "Options for locking the system")
	autoUnlock := usageBuilder.AddBooleanOption("a", "auto-unlock", true, "Automatically unlock the system when your session ends", "optional", optgroup)
	timeUnlock := usageBuilder.AddStringOption("t", "time-unlock", "", "Automatically unlock the system after a specified time", "optional", optgroup)
	usersAllowed := usageBuilder.AddStringOption("u", "users-allowed", "", "Users allowed to unlock the system", "optional", optgroup)
	groupsAllowed := usageBuilder.AddStringOption("g", "groups-allowed", "", "Groups allowed to unlock the system", "optional", optgroup)
	reason := usageBuilder.AddStringOption("r", "reason", "", "Reason for locking the system", "optional", optgroup)
	email := usageBuilder.AddStringOption("m", "email", "", "Email address to show users that try to access the system", "optional", optgroup)
	action := usageBuilder.AddArgument(1, "action", "The action to perform", "lock/unlock/authorize/status")

	_ = enable
	_ = disable

	// Parse the usage
	parsed := usageBuilder.Parse()
	if !parsed {
		usageBuilder.PrintError(fmt.Errorf("Error parsing usage"))
		os.Exit(1)
	}

	if action == nil || *action == "" {
		usageBuilder.PrintError(fmt.Errorf("Action is required. Use one of lock, unlock, authorize, or status"))
		os.Exit(1)
	}

	if *help {
		usageBuilder.PrintUsage()
		os.Exit(0)
	}

	// Create a new logrus logger
	logger := logrus.New()
	if *debug {
		// Enable human-readable logging to stdout
		logger.SetOutput(os.Stdout)
		logger.SetLevel(logrus.DebugLevel)
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		})
	} else {
		// Discard all logs
		logger.SetOutput(io.Discard)
	}

	// Find any lockfiles
	lockfiles, err := lock.FindAllLockfiles(logger)
	if err != nil {
		log.Printf("Error finding lockfiles: %v", err)
	}

	// Create a locker
	locker := pkg.NewLocker(
		options.WithAutoUnlock(*autoUnlock),
		options.WithTimeUnlock(*timeUnlock),
		options.WithUsersAllowed(*usersAllowed),
		options.WithGroupsAllowed(*groupsAllowed),
		options.WithReason(*reason),
		options.WithEmail(*email),
		options.WithLogger(logger),
	)

	// Handle the actions
	switch *action {
	case pkg.ACTION_LOCK:
		if len(lockfiles) > 0 {
			fmt.Println("The system is already locked")
			os.Exit(1)
		}
		err = locker.Lock()
		if err == os.ErrExist {
			fmt.Println("Lockfile already exists")
		} else if err != nil {
			panic(err)
		}
		fmt.Println("The system has been locked")
	case pkg.ACTION_UNLOCK:
		if len(lockfiles) == 0 {
			fmt.Println("The system is already unlocked")
			os.Exit(1)
		}
		if !pkg.UserHasLockFile(logger) {
			fmt.Println("You do not have a lock on the system")
			os.Exit(1)
		} else {
			err := locker.Unlock()
			if err != nil {
				log.Printf("Error unlocking the system: %v", err)
				os.Exit(1)
			}
		}
		fmt.Println("Unlocking the system")
	case pkg.ACTION_AUTHORIZE:
		// 0 = authorized, anything else = not authorized
		if authorized, err := locker.AuthorizeLogin(); err != nil || !authorized {
			if err != nil {
				fmt.Printf("failed to get authorization: %v\n", err)
			}
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	case pkg.ACTION_STATUS:
		if len(lockfiles) == 0 {
			fmt.Println("The system is unlocked")
		} else {
			fmt.Println("The system is locked. The following users have locks:")
			for i, lockfile := range lockfiles {
				fmt.Printf("%d | %s\n", i, lockfile)
			}
		}
	default:
		fmt.Printf("[ERROR] Unknown action: %s\n", *action)
	}
}
