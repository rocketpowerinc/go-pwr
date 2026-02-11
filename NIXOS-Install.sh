# Install required runtime deps
nix-env -iA nixos.git
nix-env -iA nixos.gnome-terminal

# Build + install + launch go-pwr
nix-shell -p go git --run '
set -euo pipefail

TMPDIR=$(mktemp -d)
trap "rm -rf \"$TMPDIR\"" EXIT

cd "$TMPDIR"
git clone https://github.com/rocketpowerinc/go-pwr.git
cd go-pwr

go build -o go-pwr ./cmd/go-pwr

mkdir -p "$HOME/.local/bin"
install -m755 go-pwr "$HOME/.local/bin/go-pwr"

# Add to PATH for this run and launch
export PATH="$HOME/.local/bin:$PATH"
go-pwr
'
