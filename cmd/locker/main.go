package main

import (
	"context"
	"fmt"
	api "locker/api/go"
	"locker/pkg"
	"os"
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
	usersAllowed := usageBuilder.AddStringOption("u", "users-allowed", "", "Users allowed to unlock the system", "optional", optgroup)
	groupsAllowed := usageBuilder.AddStringOption("g", "groups-allowed", "", "Groups allowed to unlock the system", "optional", optgroup)
	reason := usageBuilder.AddStringOption("r", "reason", "", "Reason for locking the system", "optional", optgroup)
	email := usageBuilder.AddStringOption("m", "email", "", "Email address to show users that try to access the system", "optional", optgroup)
	action := usageBuilder.AddArgument(1, "action", "The action to perform", "lock/unlock/status")

	_ = debug
	_ = enable
	_ = disable
	_ = autoUnlock
	_ = timeUnlock
	_ = usersAllowed
	_ = groupsAllowed
	_ = reason
	_ = email

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

		// Build a minimal LockRequest with dummy values.
		req := &api.LockRequest{
			User:            "testuser",
			Uid:             "1001",
			Tty:             "/dev/pts/0",
			SessionId:       "session123",
			AllowedUsers:    []string{"testuser"},
			AllowedGroups:   []string{"testgroup"},
			Reason:          "Testing lock",
			UnlockOnExit:    true,
			UnlockTime:      durationpb.New(60 * time.Minute),
			UnlockAfterIdle: durationpb.New(5 * time.Minute),
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
