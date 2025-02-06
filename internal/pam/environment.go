package pam

import "os"

func GetEnvironment() *Environment {

	user := os.Getenv("PAM_USER")
	service := os.Getenv("PAM_SERVICE")
	tty := os.Getenv("PAM_TTY")
	rhost := os.Getenv("PAM_RHOST")
	ruser := os.Getenv("PAM_RUSER")
	pamType := os.Getenv("PAM_TYPE")

	var env Environment

	if user != "" {
		env.User = &user
	}

	if service != "" {
		env.Service = &service
	}

	if tty != "" {
		env.TTY = &tty
	}

	if rhost != "" {
		env.RHost = &rhost
	}

	if ruser != "" {
		env.RUser = &ruser
	}

	if pamType != "" {
		env.PAMType = &pamType
	}

	return &env
}

type Environment struct {
	User    *string
	Service *string
	TTY     *string
	RHost   *string
	RUser   *string
	PAMType *string
}
