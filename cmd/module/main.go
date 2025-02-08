package main

/*
#cgo LDFLAGS: -lpam
#include <stdlib.h>
#include <security/pam_appl.h>
#include <security/pam_modules.h>

// Helper function to call the PAM conversation function.
// Marked as static inline to ensure internal linkage.
static inline int call_conv(struct pam_conv *conv, int num_msg,
              const struct pam_message **msg,
              struct pam_response **resp) {
    return conv->conv(num_msg, msg, resp, conv->appdata_ptr);
}

// Define typedefs so that the argv parameters in our exported Go functions
// have the same type as in the PAM headers.
typedef const char *cString;
typedef cString *cStringArray;
*/
import "C"
import (
	"fmt"
	"locker/internal/access"
	"time"
	"unsafe"
)

// Define PAM return constants. (Ensure these values match your system’s definitions.)
const (
	PAM_SUCCESS  = 0
	PAM_AUTH_ERR = 1
)

// getArgs converts the C argv array into a Go []string.
// It safely handles the case where no arguments are provided.
func getArgs(argc C.int, argv C.cStringArray) []string {
	count := int(argc)
	if count == 0 || argv == nil {
		return []string{}
	}
	// Create a slice backed by the C array.
	args := (*[1 << 28]*C.char)(unsafe.Pointer(argv))[:count:count]
	out := make([]string, count)
	for i, s := range args {
		out[i] = C.GoString(s)
	}
	return out
}

// debugEnabled checks if the "debug" argument was passed.
func debugEnabled(argc C.int, argv C.cStringArray) bool {
	for _, arg := range getArgs(argc, argv) {
		if arg == "debug" {
			return true
		}
	}
	return false
}

// sendMessage uses the PAM conversation function to send a message to the user.
func sendMessage(pamh *C.pam_handle_t, message string) C.int {
	// Retrieve the conversation pointer from the PAM context.
	var conv *C.struct_pam_conv
	ret := C.pam_get_item(pamh, C.PAM_CONV, (*unsafe.Pointer)(unsafe.Pointer(&conv)))
	if ret != C.int(PAM_SUCCESS) || conv == nil {
		return ret
	}

	// Optionally check that the conversation function pointer is not nil.
	if conv.conv == nil {
		return C.int(-1)
	}

	// Create the C string for the message.
	cMsg := C.CString(message)
	defer C.free(unsafe.Pointer(cMsg))

	// Set up a pam_message with PAM_TEXT_INFO style.
	var pamMsg C.struct_pam_message
	pamMsg.msg_style = C.PAM_TEXT_INFO
	pamMsg.msg = cMsg

	// Allocate C memory for one pointer (avoiding Go pointer to Go memory issues).
	size := C.size_t(unsafe.Sizeof((*C.struct_pam_message)(nil)))
	msgsPtr := C.malloc(size)
	if msgsPtr == nil {
		return C.int(-1) // Allocation error.
	}
	defer C.free(msgsPtr)

	// Store the address of pamMsg into the allocated C memory.
	*(**C.struct_pam_message)(msgsPtr) = &pamMsg

	var resp *C.struct_pam_response = nil

	// Call the conversation function using the C-allocated pointer.
	ret = C.call_conv(conv, 1, (**C.struct_pam_message)(msgsPtr), &resp)

	// If a response was allocated, free it.
	if resp != nil {
		C.free(unsafe.Pointer(resp))
	}
	return ret
}

// Exported PAM function for setting credentials.

//export pam_sm_setcred
func pam_sm_setcred(pamh *C.pam_handle_t, flags C.int, argc C.int, argv C.cStringArray) C.int {
	if debugEnabled(argc, argv) {
		timestamp := time.Now().Format(time.RFC3339)
		fmt.Printf("%s: pam_sm_setcred has been called\n", timestamp)
	}

	return C.int(PAM_SUCCESS)
}

// Exported PAM function for account management.
// This is typically called to enforce account restrictions after authentication.

//export pam_sm_acct_mgmt
func pam_sm_acct_mgmt(pamh *C.pam_handle_t, flags C.int, argc C.int, argv C.cStringArray) C.int {
	if debugEnabled(argc, argv) {
		timestamp := time.Now().Format(time.RFC3339)
		fmt.Printf("%s: pam_sm_acct_mgmt has been called\n", timestamp)
	}

	var pUsername *C.char
	ret := C.pam_get_user(pamh, &pUsername, C.CString("Username: "))
	if ret != C.int(PAM_SUCCESS) {
		return ret
	}

	// Convert the C string to a Go string.
	username := C.GoString(pUsername)

	// Check access policy.
	approved, reason := access.CheckAccess(username)
	if approved {
		return C.int(PAM_SUCCESS)
	}

	// Send a message to the user explaining why access was denied
	if ret := sendMessage(pamh, reason); ret != C.int(PAM_SUCCESS) {
		fmt.Println("Failed to send reason to ssh client")
	}

	// Default to denying access.
	return C.int(PAM_AUTH_ERR)
}

// Exported PAM function for authentication.
// In this sample, it always allows access by returning PAM_SUCCESS.
// after sending a message to the user.

//export pam_sm_authenticate
func pam_sm_authenticate(pamh *C.pam_handle_t, flags C.int, argc C.int, argv C.cStringArray) C.int {
	if debugEnabled(argc, argv) {
		timestamp := time.Now().Format(time.RFC3339)
		fmt.Printf("%s: pam_sm_authenticate has been called\n", timestamp)
	}

	// Nothing happens in here because this function is only called on SSH connections if the user is
	// not using a key to authenticate. Key based authentication bypasses the auth portion of the PAM flow.
	// In order to get consistent behavior, we need to use the account management function to enforce access
	// restrictions.
	return C.int(PAM_SUCCESS)
}

// Exported PAM session function for opening a session.
// This function is called during session setup (after successful authentication).

//export pam_sm_open_session
func pam_sm_open_session(pamh *C.pam_handle_t, flags C.int, argc C.int, argv C.cStringArray) C.int {
	if debugEnabled(argc, argv) {
		timestamp := time.Now().Format(time.RFC3339)
		fmt.Printf("%s: pam_sm_open_session has been called\n", timestamp)
	}

	warning := access.CheckWarning()
	if warning != "" {
		// Send a message to the user explaining why access was denied
		if ret := sendMessage(pamh, warning); ret != C.int(PAM_SUCCESS) {
			fmt.Println("Failed to send warning to ssh client")
		}
	}

	// This is not used either, mostly due to the fact that many other outputs like MOTD etc are printed at this time
	// and so our lock message gets lost in the noise as well as it is a weird user experience since it appears to have
	// connected only to be dropped and locked out which if not paying attention could result in them thinking they are
	// connected when they are not.
	return C.int(PAM_SUCCESS)
}

// Exported PAM session function for closing a session.
// This function is called when the user logs out or the session is otherwise terminated.

//export pam_sm_close_session
func pam_sm_close_session(pamh *C.pam_handle_t, flags C.int, argc C.int, argv C.cStringArray) C.int {
	if debugEnabled(argc, argv) {
		timestamp := time.Now().Format(time.RFC3339)
		fmt.Printf("%s: pam_sm_close_session has been called\n", timestamp)
	}

	// Optionally, perform cleanup tasks here.
	return C.int(PAM_SUCCESS)
}

// main is required for buildmode=c-shared.
func main() {}
