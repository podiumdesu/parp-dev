#!/bin/bash

set -exu
set -o pipefail

NETWORK_DIR=./network

GETH_BOOTNODE_PORT=30301

GETH_HTTP_PORT=8000
GETH_WS_PORT=8100
GETH_AUTH_RPC_PORT=8200
GETH_METRICS_PORT=8300
GETH_NETWORK_PORT=8400

GETH_BINARY=./dependencies/go-ethereum/build/bin/geth

NUM_NODES=1
i=1

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

$GETH_BINARY \
    --networkid=${CHAIN_ID:-32382} \
    --http \
    --http.api=eth,net,web3,personal,miner,admin \
    --http.addr=127.0.0.1 \
    --http.corsdomain="*" \
    --http.port=$((GETH_HTTP_PORT + i)) \
    --port=$((GETH_NETWORK_PORT + i)) \
    --metrics.port=$((GETH_METRICS_PORT + i)) \
    --ws \
    --ws.api=eth,net,web3 \
    --ws.addr=127.0.0.1 \
    --ws.origins="*" \
    --ws.port=$((GETH_WS_PORT + i)) \
    --authrpc.vhosts="*" \
    --authrpc.addr=127.0.0.1 \
    --authrpc.jwtsecret=$NODE_DIR/execution/jwtsecret \
    --authrpc.port=$((GETH_AUTH_RPC_PORT + i)) \
    --datadir=$NODE_DIR/execution \
    --password=$geth_pw_file \
    --bootnodes=$bootnode_enode \
    --identity=node-$i \
    --maxpendpeers=$NUM_NODES \
    --verbosity=3 \
    --syncmode=full




