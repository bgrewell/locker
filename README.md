# Locker

---

Locker is an adaptable solution that lets you allocate and revoke system access on the fly, ensuring that specific users or groups can perform exclusive tasks—such as system maintenance, temporary collaboration, or uninterrupted processing—without interference from other accounts. By combining a PAM module, a lightweight background daemon, and a user-friendly CLI, Locker simplifies the process of dynamically granting or denying logins based on real-time requirements. 

Please note that Locker is not intended to replace or enhance any native security infrastructure. It is designed purely as a convenience layer atop existing authentication and authorization systems, making it easier to manage who can log in at any given time for short-term or specialized needs.

---
> **Notice:** This project is in an early development phase and may not yet be fully stable or feature complete. As it evolves, you may encounter significant changes to the API, behavior, and overall functionality.

## Table of Contents

1. [Features](#features)  
2. [Installation](#installation)  
3. [Usage](#usage)  
   - [Locker CLI](#locker-cli)  
   - [Lock/Unlock Examples](#lockunlock-examples)  
4. [Components](#components)  
   - [Locker PAM Module](#locker-pam-module)  
   - [Locker Service](#locker-service)  
   - [Locker CLI (detailed)](#locker-cli-detailed)  
5. [How It Works](#how-it-works)  
6. [Development](#development)  
   - [Building From Source](#building-from-source)  
   - [Compile API](#compile-api)  
   - [Testing](#testing)  
7. [License](#license)

---

## Features

- **System-Level Locking**: Prevents unauthorized logins when the system is locked, allowing only the locking user, admin/sudo accounts, or those explicitly included on the access list. 
- **Flexible Unlock Options**:  
  - Auto-unlock when the locking user's session ends.  
  - Unlock at a specific time (time-unlock).  
  - Unlock after a specified period of system idle time (idle-time-unlock).  
- **Granular Access Control**: Configure lists of allowed users/groups who can bypass the lock.  
- **Warning and Messaging**: Display an MOTD or warning message to users who are allowed to log in under lock conditions.  
- **Easy CLI Commands**: `lock`, `unlock`, and `locker status` for quick control and status checks.

---

## Installation

### Quick Install Script

You can install **Locker** on most Debian/Ubuntu-based systems using the following command:

    curl -sL https://bgrewell.github.io/locker/install.sh | sudo bash

This script will:

1. Download the latest release artifacts from GitHub.  
2. Install binaries to `/opt/locker/bin`.  
3. Configure the PAM module.  
4. Set up and enable the `lockerd.service` systemd service.  
5. Create handy symlinks: `locker`, `lock`, and `unlock`.

**Note**: The install script may require you to have certain prerequisites (e.g., `curl`). If you run into issues, please make sure your system is up-to-date and has `curl` installed.

### From Source

If you prefer to install from source, jump to the [Development](#development) section for build instructions.

---

## Usage

Once installed, you can control the system lock using the following commands:

    locker lock
    locker unlock
    locker status

Or use the shortcut scripts:

    lock
    unlock

### Locker CLI

The `locker` CLI includes subcommands and options:

```bash
Usage:
  locker [OPTIONS] <action> 

Actions:
  lock      Lock the system exclusively for the current user.
  unlock    Unlock the system.
  status    Check whether the system is locked or unlocked.

Options:
  -a, --auto-unlock             Automatically unlock when your session ends (default: true)
  -t, --time-unlock <duration>  Automatically unlock after <duration> (e.g., "15m", "2h")
  -i, --idle-time-unlock <dur>  Automatically unlock after <duration> of inactivity
  -u, --users-allowed <list>    Comma-separated list of users allowed during lock
  -g, --groups-allowed <list>   Comma-separated list of groups allowed during lock
  -r, --reason <text>           Reason for locking the system
  -m, --email <address>         Email address displayed to users (optional)
      --enable                  Enable system locking (TBD future use)
      --disable                 Disable system locking (TBD future use)
  -d, --debug                   Enable debug output
  -h, --help                    Show this help message
```

### Lock/Unlock Examples

- **Lock the system** with auto-unlock on session exit:

      locker lock

- **Lock the system** until a specified time has passed:

      locker lock --time-unlock 1h --reason "Running critical experiments"

- **Unlock the system** explicitly:

      locker unlock

- **Check the system’s current status**:

      locker status

> **TIP**: You can also use the shortcut commands `lock` and `unlock` (both invoke `locker lock`/`locker unlock`).

---

## Components

### Locker PAM Module

The Locker PAM module is a shared object (`pam_locker.so`) loaded by the PAM system on user login attempts. It checks whether the system is locked, and if so, decides whether to allow or deny access based on:

- **Lock status**  
- **User groups**  
- **Allowed users**  
- **Admin/sudo privileges**  

Denied users see a custom message explaining why the system is locked.

### Locker Service

A background service (`lockerd`) runs in the background and handles:

- **auto-unlock**: Monitors sessions and unlocks once the locking session ends.  
- **time-unlock**: Automatically unlocks the system at a specified time.  
- **idle-unlock**: Unlocks the system after a configured period of inactivity.  
- **manual-control**: Coordinates lock/unlock requests from the CLI.

### Locker CLI (detailed)

The `locker` (or `lock` and `unlock`) commands allow users to:

- **Lock**: Create a lock file at `/var/lock/locker.lock` containing metadata like the locking user, lock time, optional email, allowed users/groups, and more.
- **Unlock**: Remove or expire the lock file.
- **Status**: Check if a lock file exists and retrieve its details.

When another user attempts to log in under a locked system, the PAM module checks whether they’re allowed. If they are not in the allowed list or an admin, login is denied.

---

## How It Works

**Locker** is designed for scenarios where you need to ensure uninterrupted system access for a single user or a specific set of users.

1. **Locking**  
   - When a user runs `locker lock`, Locker checks if another user has already locked the system. If locked, it denies the request. Otherwise, a lock file is created at `/var/lock/locker.lock` with details:
     - Locking username  
     - Lock time and optional unlock time  
     - Allowed users/groups  
     - Reason for lock (if any)  
     - Auto unlock flags (time-based, session-based, idle-based)
2. **PAM Interception**  
   - Upon a login attempt, the **PAM module** checks the lock file. If the system is locked and the user is neither on the allowed list nor an admin, access is denied and a message is displayed.
3. **Unlocking**  
   - The user can run `locker unlock` to release the lock.  
   - The **Locker Service** may also automatically unlock the system based on time or session rules.
4. **Admin Override**  
   - Root/sudo users always have access but will see a warning on login when the system is locked.

---

## Development

If you’d like to modify Locker, compile from source, or contribute, please follow these steps:

### Building From Source

1. Clone this repository:

       git clone https://github.com/bgrewell/locker.git
       cd locker

2. Install or verify system dependencies:

       sudo apt update && sudo apt install -y libpam0g-dev

3. Build Locker:

       make

This will compile:
- `bin/locker` (CLI)  
- `bin/lockerd` (Service)  
- `bin/pam_locker.so` (PAM module)

### Compile API

If you make changes to the protobuf definitions or need to regenerate the gRPC stubs, install the protoc plug-ins:

    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

Then compile the proto definitions:

    make proto


### Unit Testing

Locker’s code can be tested with:

    go test ./...

### End-to-End (Functional) Tests

This project includes a suite of end-to-end (E2E) tests—sometimes called *functional tests*—which verify the full workflow of installing, configuring, and using Locker in a real(ish) environment. These tests help ensure that all parts of the system work together correctly under actual usage scenarios.

### What Do the Tests Cover?
- **Environment Setup**: Installs necessary dependencies (e.g., `sshpass`), ensures DNS is operational, and deploys the Locker binaries.
- **User Management**: Creates test users (like `bob`, `jim`, `tom`) with different privileges or roles to check various access control scenarios.
- **Locker Operations**: Tests locking and unlocking the system, verifying that only allowed users can access it while locked.
- **SSH Connectivity**: Ensures that users who *should* be able to SSH in can do so, and those who *shouldn’t* cannot.

A typical test run looks like this:

```
[+] Running test setup
  running setup on localhost ......... done 
  running setup on locker-test ....... done 
  ensure sshpass is installed ........ done 
  ensure dns is working .............. done 
  install locker ..................... done 
  create user bob .................... done 
  create user jim .................... done 
  create user tom .................... done 
  ensure password login is allowed ... done 
  restart ssh ........................ done 

[+] Running tests
  00001: verify locker is installed .................. passed
  00002: test ........................................ passed
  00003: ssh to locker-test as bob ................... passed
  00004: ssh to locker-test as jim ................... passed
  00005: lock system as jim .......................... passed
  00006: ssh to locker-test as disallowed user bob ... passed
  00007: ssh to locker-test as allowed user tom ...... passed
  00008: unlock system as jim ........................ passed
  00009: verify bob can again access the system ...... passed

[+] Running test teardown
  running teardown on localhost ...... done 
  running teardown on locker-test .... done 

[+] Results
  Pass: 00009
  Fail: 00000
```

### How to Run the Tests

1. **Install or Update `dart`**  
   The tests are orchestrated by [Dart](https://github.com/bgrewell/dart), a small testing framework. You can install it via:
   ```bash
   go install github.com/bgrewell/dart/cmd/dart@latest
   ```
   or use the provided Make target (`make install_dart`).

2. **Run E2E Tests**  
   Once Dart is installed, simply run:
   ```bash
   dart -c testing/locker-e2e.yaml
   ```
   or use the Make target:
   ```bash
   make test-e2e
   ```
   This spins up the environment, executes the tests, and tears everything down automatically.

### Customizing or Extending Tests
- You can modify the `testing/locker-e2e.yaml` file to add new scenarios or tweak existing ones (e.g., adding more test cases, different user setups, etc.).
- The underlying approach supports multiple hosts or containers, so you can adapt it to more complex environments if needed.

These end-to-end tests are an essential part of validating that **Locker** behaves correctly across installation, configuration, and usage, giving you confidence in each release.


---

## License

Locker is free software, distributed under the terms of the **GNU General Public License v3 (GPLv3)**.  
Please see the [LICENSE](LICENSE) file for details.

---

**Contributions and Feedback**  
We welcome issues, bug reports, and pull requests! Please open an issue on GitHub with any questions or suggestions.

**Maintainer**: [@bgrewell](https://github.com/bgrewell)

<sub>© 2025 Ben Grewell. All rights reserved.</sub>
