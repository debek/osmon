package main

import "github.com/shirou/gopsutil/process"

// contains checks if a slice contains a specific string.
func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}

// getProcessCount retrieves the total number of running processes.
func getProcessCount() int {
    pids, err := process.Pids()
    if err != nil {
        return 0
    }
    return len(pids)
}
