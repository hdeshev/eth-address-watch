# eth-address-watch


## Prerequisites

1. Install `golangci-lint`
2. Install `pre-commit`

Nix/NixOS users can use the provided `flake.nix` and enter a dev shell with all dependencies in place by running:

```sh
nix develop
```


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
