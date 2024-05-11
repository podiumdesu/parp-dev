i=0

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
