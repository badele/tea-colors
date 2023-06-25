{ pkgs, ... }:

{
  # https://devenv.sh/packages/
  packages = with pkgs; [
    git
    goreleaser
    go-task
    expect # Preserve color with unbuffer
  ];

  # https://devenv.sh/languages/
  languages = {
    nix.enable = true;
    go.enable = true;
  };

  # https://devenv.sh/pre-commit-hooks/
  pre-commit.hooks = {
    shellcheck.enable = true;
  };
}
