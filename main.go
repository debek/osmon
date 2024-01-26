// OSMON
// Author: Daniel Debny - github.com/debek
// Year: 24.12.2023
// Description: An application for displaying system information in the terminal.

package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var interval intervalFlag
	var helpFlag bool
	var versionFlag bool

	flag.Usage = displayHelp
	flag.Var(&interval, "i", "Set interval for refreshing the display in seconds (shorthand)")
	flag.Var(&interval, "interval", "Set interval for refreshing the display in seconds")
	flag.BoolVar(&helpFlag, "h", false, "Display help information (shorthand)")
	flag.BoolVar(&helpFlag, "help", false, "Display help information")
	flag.BoolVar(&versionFlag, "v", false, "Display the version of the application (shorthand)")
	flag.BoolVar(&versionFlag, "version", false, "Display the version of the application")
	flag.Parse()

	if versionFlag {
		appVersion := os.Getenv("APP_VERSION")
		if appVersion == "" {
			appVersion = "development"
		}
		fmt.Printf("Osmon version %s\n", appVersion)
		os.Exit(0)
	}

	if helpFlag {
		displayHelp()
		os.Exit(0)
	}

	if interval.set {
		// Display system information in interval
		displaySystemInfoInInterval(interval.value)
	} else {
		// Display system information once
		displaySystemInfo()
	}
}
