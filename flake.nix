{
  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    systems.url = "github:nix-systems/default";
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, flake-parts, systems, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        go = pkgs.go_1_22;
        shell = pkgs.mkShell {
          buildInputs = [
            go
          ] ++ (with pkgs; [
            gnumake
            pre-commit
            golangci-lint
            jq
            curl
          ]);
        };
      in
      rec {
        name = "flake 1";
        description = "Cashu 1";
        devShells.default = shell;
      }
    );
}
