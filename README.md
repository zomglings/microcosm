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
docker pull fuzzyfrog/microcosm:v0.1.0
```

Create a microcosm container, bind-mounting a volume onto `/root`:

```
MICROCOSM_DIR=$(mktemp -d)
docker run \
    -e NUM_ACCOUNTS=<number of accounts to provision> \
    -v $MICROCOSM_DIR:/root \
    fuzzyfrog/microcosm:v0.1.0 \
    <geth arguments>
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


## geth arguments

As indicated above, you can directly pass in arguments for `geth` when you run the `microcosm`
docker container. For example, if you want to expose the management APIs over the JSON RPC
interface, you can run:
```
docker run -p 8545:8545 -e NUM_ACCOUNTS=0 -v $MICROCOSM_DIR:/root \
    fuzzyfrog/microcosm:test --rpc --rpcaddr 0.0.0.0 --rpcapi eth,web3
```

Note: It is important to use `--rpcaddr 0.0.0.0` because of how docker handles loopbacks within
containers -- using the default of `127.0.0.1` means you will be unable to connect to the RPC API
from outside the container.


## Deploying microcosm to a kubernetes cluster using Helm

This repository also provides a [helm](https://helm.sh/) chart that you can use to deploy
`microcosm` to a kubernetes cluster.

This creates a `StatefulSet` resource provisioned with a 100 GB persistent disk in the standard
storage class.

If you are already set up with `helm`, getting microcosm running is a simple as:
```
helm install ./helm/
```
(from this repository's root directory).

To get up and running with `helm`, follow the instructions [here](https://github.com/helm/helm).

### Modifying storage

You can deploy a custom storage class to your kubernetes cluster following these
[instructions](https://kubernetes.io/docs/concepts/storage/storage-classes/).

You can modify the size of your microcosm volume in your custom `values.yaml` file.

### Caveats

1. If you do not set `--networkid` in your geth args, state will not persist between pod restarts.
See the [`helm/values.yaml`](./helm/values.yaml) for an example.
