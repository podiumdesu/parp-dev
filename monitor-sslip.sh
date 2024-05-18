#!/bin/bash

# Check if server command or partial command and client name are provided
if [ -z "$1" ] || [ -z "$2" ]; then
  echo "Usage: $0 <server_command> <client_name>"
  exit 1
fi

server_command=$1
client_name=$2

# Define the log directory and file with a timestamp
log_dir="logs/$(date +"%Y-%m-%d")"
log_file="${log_dir}/sslip_usage-${client_name}-$(date +"%H%M").log"


# Create the log directory if it doesn't exist
mkdir -p $log_dir

# Function to get the first matching process PID
get_server_pid() {
  ps au | grep -v grep | grep "$server_command" | awk 'NR==1 {print $2}'
}

# Find the PID of the server process (get the first one if multiple are found)
server_pid=$(get_server_pid)

# Check if the server process is running
if [ -z "$server_pid" ]; then
  echo "Error: No running process found for command '$server_command'"
  exit 1
fi

echo "Monitoring server process (PID: $server_pid)..."

# Interval between checks (in seconds)
interval=1

# Clear the log file
> $log_file

while true; do
  # Check if the process is still running
  if ! ps -p $server_pid > /dev/null; then
    echo "Process $server_pid has terminated."
    exit 1
  fi

  # Get current timestamp
  timestamp=$(date +"%Y-%m-%d %H:%M:%S")

  # Get CPU and memory usage using top for more accurate results
  usage=$(top -b -n 2 -p $server_pid | grep $server_pid | awk '{print $9, $10}')
  # test=$(top -b -n 1| grep 262271)
  # echo $test
  # Log the usage
  echo "$timestamp PID: $server_pid CPU: $(echo $usage | awk '{print $1}')% MEM: $(echo $usage | awk '{print $2}')%" | tee -a $log_file

  # Wait for the next interval
  sleep $interval
done
