{
  description = "Multiplayer tetris over ssh";

  inputs = { nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable"; };

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
      buildInputs = with pkgs; [ go gopls gotools go-tools just entr fd ];
    in { devShell.${system} = pkgs.mkShell { inherit buildInputs; }; };
}
