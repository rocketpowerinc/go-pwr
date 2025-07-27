
# Go-PWR Application

## Overview
Go Power is a simple command-line application built using the Go programming language. It utilizes the Bubble Tea framework to create an interactive user interface. This project serves as a demonstration of building a terminal application in Go.

## Installation

To install the Go Power application, ensure that you have Go installed on your computer. You can download it from the official Go website: https://golang.org/dl/.

Once Go is installed, set your GOPATH correctly. Then, run the following command in your terminal:

```bash
go install -v github.com/rocketpowerinc/go-pwr@latest
```
Might have to clean cache first `go clean -modcache`

This command will compile the application and place the executable in your `$GOPATH/bin` directory.

## Usage

After installation, you can run the application from any directory by executing:

`go-pwr` or call it directly `~/go/bin/go-pwr`


## Contributing

Contributions are welcome! If you have suggestions for improvements or new features, feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

# Other Option
```bash
git clone https://github.com/rocketpowerinc/go-pwr.git
cd go-pwr
go mod init
go install
```
Then run `~/go/bin/go-pwr`