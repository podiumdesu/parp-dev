#!/bin/bash

# Check if the number of instances is provided as an argument
if [ -z "$1" ]; then
  echo "Usage: $0 <number_of_instances>"
  exit 1
fi

# Number of instances to run
num_instances=$1

GO_CMD=$(which go)

for i in $(seq 1 $num_instances); do
  # Run each instance in the background
  $GO_CMD run main.go &
done

# Wait for all background processes to finish
wait
