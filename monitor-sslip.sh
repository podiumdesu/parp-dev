#!/bin/bash

# Check if server command or partial command is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <server_command>"
  exit 1
fi

server_command=$1

# Find the PID of the server process (get the first one if multiple are found)
server_pid=$(pgrep -f "$server_command" | head -n 1)

# Check if the server process is running
if [ -z "$server_pid" ]; then
  echo "Error: No running process found for command '$server_command'"
  exit 1
fi

echo "Monitoring server process (PID: $server_pid)..."

# Interval between checks (in seconds)
interval=1

# Log file
log_file="sslip_usage.log"

# Clear the log file
> $log_file

while true; do
  # Get current timestamp
  timestamp=$(date +"%Y-%m-%d %H:%M:%S")

  # Get CPU and memory usage
  usage=$(ps -p $server_pid -o %cpu,%mem --no-headers)

  # Log the usage
  echo "$timestamp $usage" | tee -a $log_file

  # Wait for the next interval
  sleep $interval
done