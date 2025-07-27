
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


## Usage

After installation, you can run the application from any directory by executing:

`go-pwr` or call it directly `~/go/bin/go-pwr`

Add to Bash Path
 `echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc && source ~/.bashrc`


# To Get Very Latest Changes
```bash
# Clean Up
[ -d go-pwr ] && rm -rf go-pwr
rm -f ~/go/bin/go-pwr
# Install Fresh
git clone https://github.com/rocketpowerinc/go-pwr.git
cd go-pwr
go install
```
Then run `~/go/bin/go-pwr`




## Contributing

Contributions are welcome! If you have suggestions for improvements or new features, feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.