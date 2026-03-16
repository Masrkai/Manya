{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    # Go toolchain
    go
    gopls
    
    # SCSS compiler (Dart Sass, not the deprecated Ruby version)
    dart-sass
    
    # Live reload for development
    entr
  ];

}