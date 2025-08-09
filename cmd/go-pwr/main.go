// Package main is the entry point for the go-pwr application.
package main

import (
	"fmt"
	"os"

	"github.com/rocketpowerinc/go-pwr/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
