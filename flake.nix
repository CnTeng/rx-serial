{
  description = "NixOS Configuration";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

    flake-parts.url = "github:hercules-ci/flake-parts";

    treefmt = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    pre-commit = {
      url = "github:cachix/pre-commit-hooks.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [
        "x86_64-linux"
        "aarch64-linux"
      ];

      imports = [
        inputs.pre-commit.flakeModule
        inputs.treefmt.flakeModule
      ];

      perSystem =
        { config, pkgs, ... }:
        {
          devShells.default = pkgs.mkShell {
            packages = with pkgs; [
              go
              socat
              config.treefmt.build.wrapper
            ];

            CGO_ENABLED = "0";

            shellHook = config.pre-commit.installationScript;
          };

          pre-commit.settings.hooks = {
            commitizen.enable = true;
            treefmt.enable = true;
          };

          treefmt = {
            projectRootFile = "flake.nix";

            programs = {
              gofumpt.enable = true;
              nixfmt.enable = true;
              prettier.enable = true;
            };
          };
        };
    };
}
