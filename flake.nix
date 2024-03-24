{
  description = "A flake for developing brain.test";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system};
      in {
        devShells.default = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            #mdbook
            mdbook
            mdbook-mermaid
            # nodejs
            nodejs_21
            # golang
            go
            go-task
            air
            govulncheck
            gotestsum
          ];
        };
      });
}
