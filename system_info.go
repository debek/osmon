package main

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "runtime"
    "strconv"
    "strings"
    "time"
)

// getOSRelease retrieves the operating system release information.
func getOSRelease() string {
    if runtime.GOOS == "darwin" {
        out, err := exec.Command("sw_vers", "-productVersion").Output()
        if err != nil {
            return "Unknown"
        }
        return "macOS " + strings.TrimSpace(string(out))
    }
    if runtime.GOOS == "linux" {
        file, err := os.Open("/etc/os-release")
        if err != nil {
            return "Unknown"
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            line := scanner.Text()
            if strings.HasPrefix(line, "PRETTY_NAME=") {
                return strings.Trim(line[13:], "\"")
            }
        }
    }
    return "Unknown"
}

// getCPUInfo retrieves information about the CPU.
func getCPUInfo() (string, error) {
    if runtime.GOOS == "darwin" {
        out, err := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output()
        if err != nil {
            return "", err
        }
        return strings.TrimSpace(string(out)), nil
    }
    if runtime.GOOS == "linux" {
        out, err := exec.Command("cat", "/proc/cpuinfo").Output()
        if err != nil {
            return "", err
        }
        return parseLinuxCPUInfo(string(out)), nil
    }
    return "", fmt.Errorf("unsupported platform")
}

// parseLinuxCPUInfo parses CPU information from /proc/cpuinfo on Linux.
func parseLinuxCPUInfo(info string) string {
    lines := strings.Split(info, "\n")
    for _, line := range lines {
        if strings.HasPrefix(line, "model name") {
            return strings.TrimSpace(strings.Split(line, ":")[1])
        }
    }
    return "Unknown"
}

// getCPUCores retrieves the number of CPU cores.
func getCPUCores() (int, error) {
    if runtime.GOOS == "darwin" {
        out, err := exec.Command("sysctl", "-n", "hw.ncpu").Output()
        if err != nil {
            return 0, err
        }
        cores, err := strconv.Atoi(strings.TrimSpace(string(out)))
        if err != nil {
            return 0, err
        }
        return cores, nil
    }
    if runtime.GOOS == "linux" {
        out, err := exec.Command("grep", "-c", "processor", "/proc/cpuinfo").Output()
        if err != nil {
            return 0, err
        }
        cores, err := strconv.Atoi(strings.TrimSpace(string(out)))
        if err != nil {
            return 0, err
        }
        return cores, nil
    }
    return 0, fmt.Errorf("unsupported platform")
}

// getUptime retrieves the system uptime.
func getUptime() (time.Duration, error) {
    if runtime.GOOS == "darwin" {
        out, err := exec.Command("sysctl", "-n", "kern.boottime").Output()
        if err != nil {
            return 0, err
        }
        uptimeString := strings.Split(strings.Split(string(out), "=")[1], ",")[0]
        uptimeSec, err := strconv.ParseInt(strings.TrimSpace(uptimeString), 10, 64)
        if err != nil {
            return 0, err
        }
        return time.Since(time.Unix(uptimeSec, 0)), nil
    }
    if runtime.GOOS == "linux" {
        out, err := exec.Command("cat", "/proc/uptime").Output()
        if err != nil {
            return 0, err
        }
        uptimeString := strings.Fields(string(out))[0]
        uptimeSec, err := strconv.ParseFloat(uptimeString, 64)
        if err != nil {
            return 0, err
        }
        return time.Duration(uptimeSec) * time.Second, nil
    }
    return 0, fmt.Errorf("unsupported platform")
}
