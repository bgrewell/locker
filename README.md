# Locker

## Components

### Locker PAM Module

The PAM module is a shared object that is loaded by the PAM system when a user tries to log in. The PAM module is
designed to check the lock status of the system and deny the user access if the system is locked unless the user is in 
the list of allowed users or groups.

The primary purpose of the PAM module is to control access to the system based on the lock status of the system and to
present useful messages to uses when it is.

### Locker Service

- auto-unlock: Monitors sessions and unlocks when the locking session ends
- time-unlock: Unlocks the system at a specified time
- idle-unlock: Unlocks the system after a specified period of inactivity
- manual-control: Locks/Unlocks the system when a user issues the unlock command

### Locker Command Line Interface

## How it works

Locker allows users to request locking the system so that other users can not log in during the lock period. This tool
was designed for systems being used for research and development purposes, where users may need to lock the system to
prevent others from modifying or using the system during periods where they are running experiments or measurements
that could be interrupted by other users.

Locker provides commands 'lock' and 'unlock' that the user can use to request a lock on the system. When the system is
locked by a user other users will receive a message if they try to log into the system and then their login request will
be denied. Users with root access (or sudo permissions) will still be able to log into the system but will receive a
message upon login warning them that the system is currently locked.

When a user issues the 'lock' command the system first checks to ensure it is not already locked by another user. If it
is then the lock command will fail. If it is not then the lock command will create a lock file at 
`/usr/local/lib/locker/lockfile` that contains the username of the user who locked the system, the time the system 
was locked, the email address (if accessible) of the user who locked the system, a flag for if the system will auto
unlock when the users sessions have terminated, and optional fields including an unlock date/time, a lock reason, and
a list of any additional users or groups who are allowed to log in while the system is locked.

When a user tries to log in when the system is locked locker will first check to see if the auto unlock feature is set
and if it is it will check to see if the users session is no longer active and unlock if it isn't active. If it is still
active then it will check to see if the user is in the list of users who are allowed to log in while the system is 
locked or if the user is an admin user (root or sudo user). If the user is not in the list of users who are allowed to
log in while the system is locked and the user is not an admin user then the user will receive a message that the system
is locked and their login request will be denied.

If the user is allowed to log in and the system is locked they will see a MOTD which contains a warning that the system
is locked and the reason.

## Building

### Compile API

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```