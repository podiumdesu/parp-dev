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

i=0

MIN_SYNC_PEERS=1

PRYSM_BOOTSTRAP_NODE=


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
    --rpc-host=127.0.0.1 \
    --rpc-port=$((PRYSM_BEACON_RPC_PORT + i)) \
    --grpc-gateway-host=127.0.0.1 \
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


