{
  description = "feedy is rss/atom feed integration for channel talk";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
  };

  outputs = { nixpkgs, flake-utils, gomod2nix, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [
            gomod2nix.overlays.default
          ];
        };
      in
      {
        packages.default = pkgs.buildGoApplication {
          pname = "feedy";
          version = "1.0.0";
          pwd = ./.;
          src = ./.;
          modules = ./gomod2nix.toml;
        };

        devShells.default = pkgs.mkShell {
          packages = [
            pkgs.gomod2nix
            pkgs.atlas
            pkgs.sqlc
            pkgs.just
          ];
        };
      });
}
