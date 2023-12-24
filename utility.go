package main

import "github.com/shirou/gopsutil/process"

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func getProcessCount() int {
	pids, err := process.Pids()
	if err != nil {
		return 0
	}
	return len(pids)
}
