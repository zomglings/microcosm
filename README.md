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

Create a microcosm container, bind-mounting a volume onto `/root/.ethereum`:

```
mkdir /tmp/microcosm-test
docker run -v /tmp/microcosm-test:/root/.ethereum fuzzyfrog/microcosm <number of accounts to provision>
```

If you look in `/tmp/microcosm-test`, you will see the `microcosm` data directory. From outside the
container, run

```
sudo chown -R $USER:$USER /tmp/microcosm-test
```

and you will be able to use the IPC socket `/tmp/microcosm-test/geth.ipc` as a
[`web3`](https://github.com/ethereum/web3.js/) provider.