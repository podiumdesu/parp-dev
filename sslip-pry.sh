#!/bin/bash

set -exu
set -o pipefail

NETWORK_DIR=./network

GETH_AUTH_RPC_PORT=8200

PRYSM_BEACON_RPC_PORT=4000
PRYSM_BEACON_GRPC_GATEWAY_PORT=4100
PRYSM_BEACON_P2P_TCP_PORT=4200
PRYSM_BEACON_P2P_UDP_PORT=4300
PRYSM_BEACON_MONITORING_PORT=4400

PRYSM_CTL_BINARY=./dependencies/prysm/bazel-bin/cmd/prysmctl/prysmctl_/prysmctl
PRYSM_BEACON_BINARY=./dependencies/prysm/bazel-bin/cmd/beacon-chain/beacon-chain_/beacon-chain

i=2

MIN_SYNC_PEERS=1

PRYSM_BOOTSTRAP_NODE=

if [[ -z "${PRYSM_BOOTSTRAP_NODE}" ]]; then
    # sleep 5 # sleep to let the prysm node set up
    # If PRYSM_BOOTSTRAP_NODE is not set, execute the command and capture the result into the variable
    # This allows subsequent nodes to discover the first node, treating it as the bootnode
    PRYSM_BOOTSTRAP_NODE=$(curl -s localhost:4100/eth/v1/node/identity | jq -r '.data.enr')
        # Check if the result starts with enr
    if [[ $PRYSM_BOOTSTRAP_NODE == enr* ]]; then
        echo "PRYSM_BOOTSTRAP_NODE is valid: $PRYSM_BOOTSTRAP_NODE"
    else
        echo "PRYSM_BOOTSTRAP_NODE does NOT start with enr"
        exit 1
    fi
fi


NODE_DIR=$NETWORK_DIR/node-$i
mkdir -p $NODE_DIR/consensus

cp ./config.yml $NODE_DIR/consensus/config.yml
cp $NETWORK_DIR/genesis.ssz $NODE_DIR/consensus/genesis.ssz

$PRYSM_BEACON_BINARY \
    --datadir=$NODE_DIR/consensus/beacondata \
    --min-sync-peers=$MIN_SYNC_PEERS \
    --genesis-state=$NODE_DIR/consensus/genesis.ssz \
    --bootstrap-node=$PRYSM_BOOTSTRAP_NODE \
    --interop-eth1data-votes \
    --chain-config-file=$NODE_DIR/consensus/config.yml \
    --contract-deployment-block=0 \
    --chain-id=${CHAIN_ID:-32382} \
    --rpc-host=0.0.0.0 \
    --rpc-port=$((PRYSM_BEACON_RPC_PORT + i)) \
    --grpc-gateway-host=0.0.0.0\
    --grpc-gateway-port=$((PRYSM_BEACON_GRPC_GATEWAY_PORT + i)) \
    --execution-endpoint=http://localhost:$((GETH_AUTH_RPC_PORT + i)) \
    --accept-terms-of-use \
    --jwt-secret=$NODE_DIR/execution/jwtsecret \
    --suggested-fee-recipient=0x123463a4b065722e99115d6c222f267d9cabb524 \
    --minimum-peers-per-subnet=0 \
    --p2p-tcp-port=$((PRYSM_BEACON_P2P_TCP_PORT + i)) \
    --p2p-udp-port=$((PRYSM_BEACON_P2P_UDP_PORT + i)) \
    --monitoring-port=$((PRYSM_BEACON_MONITORING_PORT + i)) \
    --verbosity=info \
    --slasher \
    --enable-debug-rpc-endpoints 


