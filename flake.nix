{
  description = "Engine";

  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "aarch64-darwin"
        "x86_64-darwin"
      ];
      perSystem =
        {
          config,
          self',
          inputs',
          pkgs,
          system,
          ...
        }:
        {
          devShells.default = pkgs.mkShell.override { stdenv = pkgs.clang18Stdenv; } {
            hardeningDisable = [ "fortify" ];
            packages = with pkgs; [
              # Base:
              go
              # IDE Tools:
              gopls # LSP & Formatter
              golangci-lint
              delve
              # Libs:

              # misc:
              atac
            ];
          };
        };
    };
}
