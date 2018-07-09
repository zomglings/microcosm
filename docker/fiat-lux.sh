#!/bin/sh

interrupt_handler() {
    kill -TERM -$$
}

trap interrupt_handler INT

set -e -o pipefail

LOGGER="echo -e microcosm --> "

DEBUG=$DEBUG

set -u

# If debug is set to anything nonempty, then simply run the command passed to "docker run"
if [ ! -z $DEBUG ] ; then
    $LOGGER "Running container in DEBUG mode with command: $@"
    $@
else
    NUM_ACCOUNTS=${1:-11}

    $LOGGER "Creating $NUM_ACCOUNTS accounts"
    # Set up keystore with accounts and create genesis file
    KEYSTORE="/root/.ethereum/keystore"
    PASSWORD="microcosm"
    GENESIS_FILE=/root/genesis.json
    mkdir -p $KEYSTORE
    microcosm accounts -keystore $KEYSTORE -numAccounts $NUM_ACCOUNTS -password $PASSWORD | microcosm genesis -genesisFile $GENESIS_FILE 1>/dev/null

    $LOGGER "Configring geth"
    # Configure geth to use private chain
    geth init $GENESIS_FILE

    # Prepare to unlock generated accounts in geth
    USER_STRING=$(microcosm addresses $KEYSTORE/* | tr ' ' ',' | sed 's/,$//')
    PASSWORD_STRING=""
    for account in $(microcosm addresses $KEYSTORE/*) ; do
        PASSWORD_STRING="${PASSWORD_STRING},$PASSWORD"
    done
    PASSWORD_STRING=$(echo $PASSWORD_STRING | sed 's/^,//')

    PASSWORD_FILE=/root/passwords.txt
    echo $PASSWORD_STRING | tr ',' '\n' > $PASSWORD_FILE

    # Run a mining node on the private net with the specified accounts unlocked (and with the
    # oldest one as the coinbase)
    ETHERBASE=$(echo $USER_STRING | awk -F',' '{print $1}')

    REGULAR_ACCOUNTS=$(echo $USER_STRING | sed 's/^[^,]*,//')
    $LOGGER "etherbase account: $ETHERBASE"
    $LOGGER "unlocked accounts: $REGULAR_ACCOUNTS"

    $LOGGER "Starting geth"
    geth --mine --minerthreads 1 --unlock $USER_STRING --password $PASSWORD_FILE --etherbase $ETHERBASE
fi
