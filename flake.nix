{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";

    pre-commit-hooks = {
      url = "github:cachix/pre-commit-hooks.nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
  };

  outputs = { self, nixpkgs, flake-utils, pre-commit-hooks }:

    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShellNoCC {

          name = "go";

          buildInputs = with pkgs; [
            go
            gopls
            gotools
            golint
            python3.pkgs.grip
            nixpkgs-fmt
            jq
          ];

          inherit (self.checks.${system}.pre-commit-check) shellHook;
        };

        packages = rec {
          uciml = pkgs.callPackage ./default.nix { };

          uciml-gen = pkgs.writeShellApplication {
            name = "uciml-gen";

            runtimeInputs = with pkgs; [ curl ];

            text = ''
              export FILE="${./data.json}"
              export UCIML="${pkgs.lib.getExe uciml}"

              ${builtins.readFile ./interactive_generator.sh}
            '';
          };
        };

        checks.pre-commit-check = pre-commit-hooks.lib.${system}.run {
          src = ./.;

          hooks = {
            nixpkgs-fmt.enable = true;

            golangci-lint = {
              enable = true;

              # The name of the hook (appears on the report table):
              name = "Golangci lint";

              # The command to execute (mandatory):
              entry = "${pkgs.golangci-lint}/bin/golangci-lint run";

              # The pattern of files to run on (default: "" (all))
              # see also https://pre-commit.com/#hooks-files
              files = "\\.go$";

              # List of file types to run on (default: [ "file" ] (all files))
              # see also https://pre-commit.com/#filtering-files-with-types
              # You probably only need to specify one of `files` or `types`:
              types = [ "text" "go" ];
            };

          };
        };
      }
    );
}
