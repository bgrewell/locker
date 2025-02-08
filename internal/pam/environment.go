package pam

import "os"

func GetEnvironment() *Environment {

	user := os.Getenv("PAM_USER")
	service := os.Getenv("PAM_SERVICE")
	tty := os.Getenv("PAM_TTY")
	rhost := os.Getenv("PAM_RHOST")
	ruser := os.Getenv("PAM_RUSER")
	pamType := os.Getenv("PAM_TYPE")

	return &Environment{
		User:    user,
		Service: service,
		TTY:     tty,
		RHost:   rhost,
		RUser:   ruser,
		PAMType: pamType,
	}
}

type Environment struct {
	User    string
	Service string
	TTY     string
	RHost   string
	RUser   string
	PAMType string
}
