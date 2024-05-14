#!/bin/bash

GETH_DIR=./dependencies/geth-sslip

( cd $GETH_DIR && make geth )

set -exu
set -o pipefail

NETWORK_DIR=./network

GETH_BOOTNODE_PORT=30301

GETH_HTTP_PORT=8000
GETH_WS_PORT=8100
GETH_AUTH_RPC_PORT=8200
GETH_METRICS_PORT=8300
GETH_NETWORK_PORT=8400

GETH_BINARY=./dependencies/geth-sslip/build/bin/geth

NUM_NODES=1
i=2

NODE_DIR=$NETWORK_DIR/node-$i

rm -rf "$NODE_DIR/execution" || echo "no network directory"
mkdir -p $NODE_DIR/execution

geth_pw_file="$NODE_DIR/geth_password.txt"
echo "" > "$geth_pw_file"

cp $NETWORK_DIR/genesis.json $NODE_DIR/execution/genesis.json


bootnode_enode=$(head -n 1 $NETWORK_DIR/bootnode/bootnode.log)
# Check if the line begins with "enode"
if [[ "$bootnode_enode" == enode* ]]; then
    echo "bootnode enode is: $bootnode_enode"
else
    echo "The bootnode enode was not found. Exiting."
    exit 1
fi


$GETH_BINARY account new --datadir "$NODE_DIR/execution" --password "$geth_pw_file"

$GETH_BINARY init \
    --datadir=$NODE_DIR/execution \
    $NODE_DIR/execution/genesis.json

echo
echo
echo


$GETH_BINARY \
    --networkid=${CHAIN_ID:-32382} \
    --http \
    --http.api=eth,net,web3,personal,miner,admin \
    --http.addr=0.0.0.0 \
    --http.corsdomain="*" \
    --http.port=$((GETH_HTTP_PORT + i)) \
    --port=$((GETH_NETWORK_PORT + i)) \
    --metrics.port=$((GETH_METRICS_PORT + i)) \
    --ws \
    --ws.api=eth,net,web3 \
    --ws.addr=0.0.0.0 \
    --ws.origins="*" \
    --ws.port=$((GETH_WS_PORT + i)) \
    --authrpc.vhosts="*" \
    --authrpc.addr=0.0.0.0 \
    --authrpc.jwtsecret=$NODE_DIR/execution/jwtsecret \
    --authrpc.port=$((GETH_AUTH_RPC_PORT + i)) \
    --datadir=$NODE_DIR/execution \
    --password=$geth_pw_file \
    --bootnodes=$bootnode_enode \
    --identity=node-$i \
    --maxpendpeers=$NUM_NODES \
    --verbosity=3 \
    --syncmode=full


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
    --grpc-gateway-host=0.0.0.0 \
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


NETWORK_DIR=./network
NODE_DIR=$NETWORK_DIR/node-$i


PRYSM_VALIDATOR_RPC_PORT=7000
PRYSM_VALIDATOR_GRPC_GATEWAY_PORT=7100
PRYSM_VALIDATOR_MONITORING_PORT=7200

PRYSM_VALIDATOR_BINARY=./dependencies/prysm/bazel-bin/cmd/validator/validator_/validator


PRYSM_BEACON_RPC_PORT=4000

$PRYSM_VALIDATOR_BINARY \
    --beacon-rpc-provider=localhost:$((PRYSM_BEACON_RPC_PORT + i)) \
    --datadir=$NODE_DIR/consensus/validatordata \
    --accept-terms-of-use \
    --interop-num-validators=1 \
    --interop-start-index=$i \
    --rpc-port=$((PRYSM_VALIDATOR_RPC_PORT + i)) \
    --grpc-gateway-port=$((PRYSM_VALIDATOR_GRPC_GATEWAY_PORT + i)) \
    --monitoring-port=$((PRYSM_VALIDATOR_MONITORING_PORT + i)) \
    --graffiti="node-$i" \
    --chain-config-file=$NODE_DIR/consensus/config.yml



