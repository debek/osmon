package main

import (
	"flag"
	"fmt"
	"strconv"
)

type intervalFlag struct {
	set   bool
	value int
}

// Set is a method called by the flag package when the flag is set.
func (f *intervalFlag) Set(s string) error {
	f.set = true
	var err error
	f.value, err = strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("invalid format for interval: %v", err)
	}
	return nil
}

// String returns the string representation of the flag
func (f *intervalFlag) String() string {
	if !f.set {
		return "not set"
	}
	return fmt.Sprintf("%d", f.value)
}

func DefineFlags() (intervalFlag, *bool, *bool) {
	var interval intervalFlag
	var helpFlag = flag.Bool("h", false, "Display help information (shorthand)")
	var versionFlag = flag.Bool("v", false, "Display the version of the application (shorthand)")

	flag.Var(&interval, "i", "Set interval for refreshing the display in seconds (shorthand)")
	flag.Var(&interval, "interval", "Set interval for refreshing the display in seconds")

	return interval, helpFlag, versionFlag
}
