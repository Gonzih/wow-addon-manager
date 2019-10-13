with import <nixpkgs> {};
pkgs.mkShell {
  buildInputs = [
    gnumake go_1_12 chromium
  ];
}
