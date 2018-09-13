#!/bin/sh

interrupt_handler() {
    kill -TERM -$$
}

trap interrupt_handler INT

set -e

LOGGER="echo microcosm --> "

DEBUG=$DEBUG

MICROCOSM_DIR=${MICROCOSM_DIR:-/root}
DATA_DIR=$MICROCOSM_DIR/.ethereum
mkdir -p $DATA_DIR

set -u

# If debug is set to anything nonempty, then simply run the command passed to "docker run"
if [ ! -z $DEBUG ] ; then
    $LOGGER "Running container in DEBUG mode with command: $@"
    $@
else
    NUM_ACCOUNTS=${NUM_ACCOUNTS:-11}

    $LOGGER "Creating $NUM_ACCOUNTS accounts"
    # Set up keystore with accounts and create genesis file
    KEYSTORE="$DATA_DIR/keystore"
    PASSWORD="microcosm"
    GENESIS_FILE=$MICROCOSM_DIR/genesis.json
    INITIALIZATION_FILE=$MICROCOSM_DIR/init
    mkdir -p $KEYSTORE

    NEW_ACCOUNTS=$(microcosm accounts -keystore $KEYSTORE -numAccounts $NUM_ACCOUNTS -password $PASSWORD)

    $LOGGER "Configuring geth:"

    if [ ! -f $GENESIS_FILE ]; then
        $LOGGER "Creating genesis file"
        microcosm genesis -genesisFile $GENESIS_FILE $NEW_ACCOUNTS 1>/dev/null
    fi

    if [ ! -f $INITIALIZATION_FILE ] ; then
        $LOGGER "geth initialization"
        # Configure geth to use private chain
        geth --datadir $DATA_DIR init $GENESIS_FILE
        echo "SUCCESS" >> $INITIALIZATION_FILE
    fi

    $LOGGER "geth configuration complete"

    # Prepare microcosm account and password files
    ACCOUNTS_FILE=$MICROCOSM_DIR/accounts.txt
    PASSWORDS_FILE=$MICROCOSM_DIR/passwords.txt
    for account in $NEW_ACCOUNTS; do
        echo "$account" >>$ACCOUNTS_FILE
        echo "$PASSWORD" >>$PASSWORDS_FILE
    done

    # Run a mining node on the private net with the specified accounts unlocked (and with the
    # oldest one as the coinbase)
    ETHERBASE=$(head -n1 $ACCOUNTS_FILE)
    REGULAR_ACCOUNTS=$(tail -n+2 $ACCOUNTS_FILE | tr '\n' ',')
    ACCOUNTS_STRING="$ETHERBASE,$REGULAR_ACCOUNTS"

    $LOGGER "etherbase account: $ETHERBASE"
    $LOGGER "unlocked accounts: $REGULAR_ACCOUNTS"

    $LOGGER "Starting geth"
    geth --datadir $DATA_DIR --mine --minerthreads 1 --unlock $ACCOUNTS_STRING --password $PASSWORDS_FILE --etherbase $ETHERBASE $@
fi
