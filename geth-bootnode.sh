#!/bin/bash

set -exu
set -o pipefail

GETH_BOOTNODE_PORT=30301
GETH_BOOTNODE_BINARY=./dependencies/go-ethereum/build/bin/bootnode

NETWORK_DIR=./network

trap 'echo "Error on line $LINENO"; exit 1' ERR
# Function to handle the cleanup
cleanup() {
    echo "Caught Ctrl+C. Killing active background processes and exiting."
    kill $(jobs -p)  # Kills all background processes started in this script
    exit
}
# Trap the SIGINT signal and call the cleanup function when it's caught
trap 'cleanup' SIGINT

rm -rf "$NETWORK_DIR" || echo "no network directory"
mkdir -p $NETWORK_DIR
pkill geth || echo "No existing geth processes"
pkill beacon-chain || echo "No existing beacon-chain processes"
pkill validator || echo "No existing validator processes"
pkill bootnode || echo "No existing bootnode processes"



PRYSM_CTL_BINARY=./dependencies/prysm/bazel-bin/cmd/prysmctl/prysmctl_/prysmctl

NUM_NODES=2

$PRYSM_CTL_BINARY testnet generate-genesis \
--fork=deneb \
--num-validators=$NUM_NODES \
--chain-config-file=./config.yml \
--geth-genesis-json-in=./genesis.json \
--output-ssz=$NETWORK_DIR/genesis.ssz \
--geth-genesis-json-out=$NETWORK_DIR/genesis.json


mkdir -p $NETWORK_DIR/bootnode

$GETH_BOOTNODE_BINARY -genkey $NETWORK_DIR/bootnode/nodekey

$GETH_BOOTNODE_BINARY \
    -nodekey $NETWORK_DIR/bootnode/nodekey \
    -addr=:$GETH_BOOTNODE_PORT \
    -verbosity=5 | tee -a "$NETWORK_DIR/bootnode/bootnode.log" 


