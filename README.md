# eth-address-watch

Consume blocks from an Ethereum node as they are mined and keep track of the addresses that have been used in the transactions. Allow subscribing to certain addresses and support notifications when new transactions to/from those addresses are found.

Architecture description in [ARCHITECTURE](ARCHITECTURE.md)


## Prerequisites

1. Install `golangci-lint`
2. Install `pre-commit`

Nix/NixOS users can use the provided `flake.nix` and enter a dev shell with all dependencies in place by running:

```sh
nix develop
```

Tested on Linux, but the above Nix setup should work on macOS too.

### Installing dependencies

To install Go dependencies associated with `eth-address-watch`, run the
command

```sh
make install
```

### Setting up `pre-commit`

Install the git hooks to run checks before committing via `pre-commit`.

```sh
make local-setup
```

### Using Code Formatters

Format code with

```sh
make codestyle
```

### Using Code Linters

Run code linters with

```sh
make lint
```

### Running Tests

Run tests with

```sh
make test
```

### Running the service

Local developer run, assuming you have all the dependencies and the Go toolchain

```sh
make run
```

Docker-based run -- the only dependency is a working Docker installation

```sh
make docker-gen # builds the image
make docker-run # runs the service and exposes port 9000
```

### API requests

Now that you have the service running, you can make requests to it.

Get the current ETH block obtained from the node

```sh
curl 'http://localhost:9000/block'
```

Subscribe for transactions involving a specific address

```sh
curl http://localhost:9000/subscribe \
    --data '{"address":"0xdac17f958d2ee523a2206206994597c13d831ec7"}'
```

Note: The `0xdac17f958d2ee523a2206206994597c13d831ec7` address is the address of the [USDT smart contract](https://etherscan.io/address/0xdac17f958d2ee523a2206206994597c13d831ec7) and is a good candidate for testing since it gets transactions all the time.

Get transactions we have discovered for a subscribed address

```sh
curl 'http://localhost:9000/transactions?address=0xdac17f958d2ee523a2206206994597c13d831ec7' \
    | jq '.data'
```

Note the use of `jq` above to pretty-print the output. Since the USDT smart contract address has a lot of transactions, you could use a `jq` trick to get just the total number of transactions (by fetching the length of the `data` response attribute).

```sh
curl -s 'http://localhost:9000/transactions?address=0xdac17f958d2ee523a2206206994597c13d831ec7' \
    | jq '.data | length'
```
