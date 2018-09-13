# microcosm

`microcosm` allows you to spin up a single, mining Ethereum node that you can use to test:

1. Smart-contract-based decentralized applications

1. Code which interacts with an Ethereum blockchain even if it doesn't live on the blockchain

1. Modifications to Ethereum node implementations

In the first two capacities, it is a complement to [Ganache](https://truffleframework.com/ganache),
which diverges from Ethereum clients like [`geth`](https://github.com/ethereum/go-ethereum/wiki/geth)
and [`parity`](https://www.parity.io/) in its implementation of the JSON-RPC specification.

With `microcosm`, currently, you get a `geth` node to test with.


## Requirements

The only requirement is [Docker](https://www.docker.com/get-docker).


## Getting started

Pull the latest `microcosm` image from DockerHub:

```
docker pull fuzzyfrog/microcosm
```

Create a microcosm container, bind-mounting a volume onto `/root`:

```
MICROCOSM_DIR=$(mktemp -d)
docker run -v $MICROCOSM_DIR:/root fuzzyfrog/microcosm <number of accounts to provision>
```

If you look in `$MICROCOSM_DIR`, you will see the `microcosm` directory. This directory
contains the `geth` data directory as a subdirectory -- `$MICROCOSM_DIR/.ethereum`.

It also contains the following files:

1. `genesis.json` - Genesis file used to initialize the `microcosm` network being run

2. `init` - File denoting that the network initialization was successful

3. `accounts.txt` - File listing the addresses of accounts created by `microcosm`

4. `passwords.txt` - File listing the passwords corresponding to each account in `accounts.txt`

The items in `$MICROCOSM_DIR` are owned by `root`. To take ownership of them, from outside the
container, run
```
sudo chown -R $USER:$USER $MICROCOSM_DIR
```

Now, you will be able to use the IPC socket `$MICROCOSM_DIR/geth.ipc` as a
[`web3`](https://github.com/ethereum/web3.js/) provider.


For a side-by-side view of the `microcosm`-generated accounts and passwords, you can run:
```
pr -w 100 -m -t $MICROCOSM_DIR/accounts.txt $MICROCOSM_DIR/passwords.txt
```
