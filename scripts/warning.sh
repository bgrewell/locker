#!/bin/bash
# ANSI escape codes:
#   Bold, Flashing Red: \033[1;5;31m
#   Bold White:         \033[1;37m
#   Reset:              \033[0m

# Display the flashing warning message
echo -e "\033[1;5;31mWARNING: SYSTEM HAS BEEN RESERVED BY 'ben' AND IS LOCKED\033[0m"

# Display the detailed message in white
echo -e "\033[1;37mAccess to this system is restricted to prevent interference with ongoing critical operations.
Please coordinate with 'ben' if you require access.\033[0m"
