#!/bin/bash

# Compile on Linux
echo "Building binary..."
go build ./cmd/servnode

# Clear home folder
HOMEDIR="$HOME/.serv"

echo "Clearing home folder..."
rm -rf "$HOMEDIR"

KEY="dev0"
CHAINID="serv_43970-1"
MONIKER="localtestnet"
BINARY="servnode"
KEYRING="test"
KEYALGO="eth_secp256k1"
LOGLEVEL="info"
MIN_DENOM="aservo"
TRACE=""
echo $HOMEDIR
ETHCONFIG="$HOMEDIR/config/config.toml"
GENESIS="$HOMEDIR/config/genesis.json"
TMPGENESIS="$HOMEDIR/config/tmp_genesis.json"

# Build binary
echo "Building binary..."
go build ./cmd/$BINARY

# Set keyring-backend and chain-id
echo "Configuring $BINARY..."
./$BINARY config keyring-backend $KEYRING
./$BINARY config chain-id $CHAINID

# Add key
echo "Adding key $KEY..."
./$BINARY keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO

# Set moniker for the node
echo "Initializing $BINARY..."
./$BINARY init $MONIKER --chain-id $CHAINID

# Change parameter token denominations to native coin
echo "Updating genesis file..."
jq '.app_state.staking.params.bond_denom = $min_denom | 
    .app_state.crisis.constant_fee.denom = $min_denom | 
    .app_state.gov.deposit_params.min_deposit[0].denom = $min_denom | 
    .app_state.gov.params.min_deposit[0].denom = $min_denom | 
    .app_state.mint.params.mint_denom = $min_denom' \
  --arg min_denom "$MIN_DENOM" $GENESIS > $TMPGENESIS && mv $TMPGENESIS $GENESIS

# Increase block time
echo "Increasing block time..."
jq '.consensus_params.block.time_iota_ms = "30000"' $GENESIS > $TMPGENESIS && mv $TMPGENESIS $GENESIS

# Gas limit in genesis
echo "Setting gas limit in genesis..."
jq '.consensus_params.block.max_gas = "10000000"' $GENESIS > $TMPGENESIS && mv $TMPGENESIS $GENESIS

# Setup
echo "Setting up..."
sed -i 's/create_empty_blocks = true/create_empty_blocks = false/g' $ETHCONFIG

# Allocate genesis accounts (cosmos formatted addresses)
echo "Allocating genesis accounts..."
./$BINARY add-genesis-account $KEY 100000000000000000000000000$MIN_DENOM --keyring-backend $KEYRING

# Sign genesis transaction
echo "Signing genesis transaction..."
./$BINARY gentx $KEY 1000000000000000000000$MIN_DENOM --keyring-backend $KEYRING --chain-id $CHAINID

# Collect genesis tx
echo "Collecting genesis tx..."
./$BINARY collect-gentxs

# Run this to ensure everything worked and that the genesis file is set up correctly
echo "Validating genesis..."
./$BINARY validate-genesis

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
echo "Starting the node..."
./$BINARY start --pruning=nothing $TRACE --log_level $LOGLEVEL --minimum-gas-prices=0.0001$MIN_DENOM --json-rpc.api eth,txpool,personal,net,debug,web3 --api.enable --grpc.enable true --chain-id $CHAINID
