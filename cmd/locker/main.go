package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	api "locker/api/go"
	"locker/pkg"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bgrewell/usage"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
)

var (
	version   string
	builddate string
	commit    string
	branch    string
)

//// NewLockFile creates a new LockFile struct with the provided parameters.
//func NewLockFile(reason string, unlockTime time.Time, unlockOnExit bool, allowedUsers, allowedGroups []string, logger *zap.Logger) (lockfile *LockFile, err error) {
//
//	// Use loginctl to get session information.
//	sessionID, tty, userName, uid, err := getSessionStatus()
//	logger.Debug("session information",
//		zap.String("session_id", sessionID),
//		zap.String("tty", tty),
//		zap.String("user", userName),
//		zap.Int("uid", uid))
//	if err != nil {
//		// If loginctl fails, fall back to os/user.
//		cu, err2 := user.Current()
//		if err2 != nil {
//			return nil, fmt.Errorf("failed to get current user: %w (and loginctl error: %v)", err2, err)
//		}
//		userName = cu.Username
//		uid, _ = strconv.Atoi(cu.Uid)
//		tty = os.Getenv("SSH_TTY")
//		sessionID = ""
//	}
//
//	email, err := getUserEmail(userName)
//	if err != nil {
//		email = ""
//		logger.Warn("failed to get email for user",
//			zap.String("user", userName),
//			zap.Error(err))
//	} else {
//		logger.Debug("email found",
//			zap.String("email", email),
//			zap.String("user", userName))
//	}
//
//	return &LockFile{
//		User:          userName,
//		UID:           uid,
//		Email:         email,
//		Reason:        reason,
//		Session:       sessionID,
//		TTY:           tty,
//		UnlockTime:    unlockTime,
//		UnlockOnExit:  unlockOnExit,
//		AllowedUsers:  allowedUsers,
//		AllowedGroups: allowedGroups,
//		logger:        logger,
//	}, nil
//}

// getSessionStatus invokes "loginctl session-status --no-pager" and parses its output.
// It returns the session ID, TTY, username, and UID.
func getSessionStatus() (sessionID, tty, userName string, uid int, err error) {
	cmd := exec.Command("/usr/bin/loginctl", "session-status", "--no-pager")
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

func getUserEmail(username string) (string, error) {
	grpcAddr := "localhost:5128"
	// TODO: This function is a mess, clean it and some of the main function up. Refactor and improve reuse
	// Set a timeout duration for the RPC calls.
	timeout := 2 * time.Second

	dialServer := func() *grpc.ClientConn {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		conn, err := grpc.DialContext(ctx, grpcAddr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			panic(err)
		}
		return conn
	}

	conn := dialServer()
	defer conn.Close()
	client := api.NewLockerServiceClient(conn)

	req := &api.EmailRequest{
		Username: username,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resp, err := client.Email(ctx, req)
	if err != nil || !resp.Found {
		return "", err
	}

	return resp.Email, nil
}

func main() {

	// Setup application usage
	usageBuilder := usage.NewUsage(
		usage.WithApplicationName("api"),
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
	idleTimeUnlock := usageBuilder.AddStringOption("i", "idle-time-unlock", "", "Automatically unlock the system after a specified idle time", "optional", optgroup)
	usersAllowed := usageBuilder.AddStringOption("u", "users-allowed", "", "Users allowed to unlock the system", "optional", optgroup)
	groupsAllowed := usageBuilder.AddStringOption("g", "groups-allowed", "", "Groups allowed to unlock the system", "optional", optgroup)
	reason := usageBuilder.AddStringOption("r", "reason", "", "Reason for locking the system", "optional", optgroup)
	email := usageBuilder.AddStringOption("m", "email", "", "Email address to show users that try to access the system", "optional", optgroup)
	action := usageBuilder.AddArgument(1, "action", "The action to perform", "lock/unlock/status")

	_ = debug   //TODO: Figure out if this should still be here
	_ = enable  //TODO: Implement this, should control the service and the PAM module entries
	_ = disable //TODO: Implement this, should control the service and the PAM module entries
	//_ = autoUnlock
	//_ = timeUnlock
	//_ = usersAllowed
	//_ = groupsAllowed
	//_ = reason
	//_ = email

	// Parse the usage
	parsed := usageBuilder.Parse()
	if !parsed {
		usageBuilder.PrintError(fmt.Errorf("Error parsing usage"))
		os.Exit(1)
	}

	if action == nil || *action == "" {
		usageBuilder.PrintError(fmt.Errorf("Action is required. Use one of lock, unlock, or status"))
		os.Exit(1)
	}

	if *help {
		usageBuilder.PrintUsage()
		os.Exit(0)
	}

	// Create a new logrus logger
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	// Define the gRPC server address; adjust as needed.
	grpcAddr := "localhost:5128"

	// Set a timeout duration for the RPC calls.
	timeout := 2 * time.Second

	// Common function to dial the gRPC server.
	dialServer := func() *grpc.ClientConn {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		conn, err := grpc.DialContext(ctx, grpcAddr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			logger.Fatalf("Failed to connect to gRPC server at %s: %v", grpcAddr, err)
		}
		return conn
	}

	switch *action {
	case pkg.ACTION_LOCK:
		fmt.Println("Locking the system")
		conn := dialServer()
		defer conn.Close()
		client := api.NewLockerServiceClient(conn)

		session, tty, user, uid, err := getSessionStatus()
		if err != nil {
			fmt.Printf("Failed to get session information: %v\n", err)
			os.Exit(1)
		}

		if *email == "" {
			*email, err = getUserEmail(user)
		}

		unlockTime, err := time.ParseDuration(*timeUnlock)
		if err != nil {
			unlockTime = 0
		}

		idleTime, err := time.ParseDuration(*idleTimeUnlock)
		if err != nil {
			idleTime = 0
		}

		// Build a minimal LockRequest with dummy values.
		req := &api.LockRequest{
			User:            user,
			Uid:             int32(uid),
			Tty:             tty,
			SessionId:       session,
			AllowedUsers:    strings.Split(*usersAllowed, ","),
			AllowedGroups:   strings.Split(*groupsAllowed, ","),
			Reason:          *reason,
			UnlockOnExit:    *autoUnlock,
			UnlockTime:      durationpb.New(unlockTime),
			UnlockAfterIdle: durationpb.New(idleTime),
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		resp, err := client.Lock(ctx, req)
		if err != nil {
			fmt.Printf("Lock RPC failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Lock RPC response: Success=%v, Message=%s\n", resp.Success, resp.Message)

	case pkg.ACTION_UNLOCK:
		fmt.Println("Unlocking the system")
		conn := dialServer()
		defer conn.Close()
		client := api.NewLockerServiceClient(conn)

		req := &api.UnlockRequest{}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		resp, err := client.Unlock(ctx, req)
		if err != nil {
			fmt.Printf("Unlock RPC failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Unlock RPC response: Success=%v, Message=%s\n", resp.Success, resp.Message)

	case pkg.ACTION_STATUS:
		fmt.Println("Checking the status of the system")
		conn := dialServer()
		defer conn.Close()
		client := api.NewLockerServiceClient(conn)

		req := &api.StatusRequest{}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		resp, err := client.Status(ctx, req)
		if err != nil {
			fmt.Printf("Status RPC failed: %v\n", err)
			os.Exit(1)
		}
		// For minimal output, print the state.
		fmt.Printf("Status RPC response: State=%v\n", resp.State)

	default:
		fmt.Printf("[ERROR] Unknown action: %s\n", *action)
	}
}
