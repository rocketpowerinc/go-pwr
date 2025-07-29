# Go-PWR Application

## Overview

`go-pwr` is a simple command-line application built using the Go programming language. It utilizes the Bubble Tea framework to create an interactive user interface.

## Installation

To install the Go Power application, ensure that you have Go installed on your computer. You can download it from the official Go website: https://golang.org/dl/.

Once Go is installed, set your GOPATH correctly. Then, run the following command in your terminal:

```bash
go install -v github.com/rocketpowerinc/go-pwr@latest
```

This command will compile the application and place the executable in your `$GOPATH/bin` directory.

Might have to clean cache first `go clean -modcache`

If still not getting the latest changes use

```bash
[ -d go-pwr ] && rm -rf go-pwr
rm -f ~/go/bin/go-pwr
git clone https://github.com/rocketpowerinc/go-pwr.git
cd go-pwr
go install
```

# Aliases

- Windows

```
$gooAlias = @"
function goo {
    Remove-Item -Recurse -Force go-pwr -ErrorAction SilentlyContinue
    Remove-Item -Force `"$HOME\go\bin\go-pwr.exe`" -ErrorAction SilentlyContinue
    git clone https://github.com/rocketpowerinc/go-pwr.git
    Set-Location go-pwr
    go install
    & "$env:USERPROFILE\go\bin\go-pwr.exe"
}
"@

$gooAlias | Out-File -Append -Encoding UTF8 $PROFILE
```

- Linux

```
cat << 'EOF' >> ~/.bashrc
# Alias to reinstall go-pwr
alias goo='
  rm -rf go-pwr && \
  rm -f ~/go/bin/go-pwr && \
  git clone https://github.com/rocketpowerinc/go-pwr.git && \
  cd go-pwr && \
  go install && \
  ~/go/bin/go-pwr'
EOF
```

- Mac

```
cat << 'EOF' >> ~/.zshrc
# Alias to reinstall go-pwr
alias goo='
  rm -rf go-pwr && \
  rm -f ~/go/bin/go-pwr && \
  git clone https://github.com/rocketpowerinc/go-pwr.git && \
  cd go-pwr && \
  go install && \
  ~/go/bin/go-pwr'
EOF
```

- Then Source
  - `. $PROFILE`
  - `source ~/.bashrc`
  - `source ~/.zshrc`

## Usage/Export

After installation, you can run the application from any directory by executing:

`go-pwr` or call it directly `~/go/bin/go-pwr`

Linux - Add to Bash Path
`echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc && source ~/.bashrc`
Mac - Add to ZSH Path
`echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc && source ~/.zshrc`

## Contributing

Contributions are welcome! If you have suggestions for improvements or new features, feel free to open an issue or submit a pull request.

### License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
