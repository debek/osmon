package main

import (
    "fmt"
    "github.com/shirou/gopsutil/disk"
    "github.com/shirou/gopsutil/load"
    "github.com/shirou/gopsutil/mem"
    "os"
    "os/exec"
    "os/user"
    "runtime"
    "strings"
    "time"
)

// displaySystemInfoAtInterval displays system information at specified intervals.
func displaySystemInfoAtInterval(interval time.Duration) {
    for {
        clearScreen()
        displaySystemInfo()
        time.Sleep(interval)
    }
}

// displaySystemInfo collects and displays system information.
func displaySystemInfo() {
    currentUser, err := user.Current()
    if err != nil {
        fmt.Printf("Failed to get current user: %s\n", err)
        return
    }
    userName := currentUser.Username

    currentPTS, err := getCurrentPTS()
    if err != nil {
        fmt.Printf("Failed to get current PTS: %s\n", err)
        currentPTS = "Unknown"
    }

    activeUsers, err := getActiveUsers()
    if err != nil {
        fmt.Printf("Error getting active users: %s\n", err)
        activeUsers = make(map[string][]string)
    }

    hostName, err := os.Hostname()
    if err != nil {
        fmt.Printf("Failed to get hostname: %s\n", err)
        hostName = "Unknown"
    }

    ip, err := getLocalIP()
    if err != nil {
        fmt.Println("No network access")
        ip = "No connection"
    }

    uptime, err := getUptime()
    if err != nil {
        fmt.Printf("Failed to get uptime: %s\n", err)
        uptime = 0
    }
    upDays := uptime / (24 * time.Hour)
    upHours := (uptime % (24 * time.Hour)) / time.Hour
    upMins := (uptime % time.Hour) / time.Minute
    upSecs := (uptime % time.Minute) / time.Second

    vmStat, err := mem.VirtualMemory()
    if err != nil {
        fmt.Printf("Failed to get virtual memory info: %s\n", err)
        vmStat = &mem.VirtualMemoryStat{}
    }
    memoryUsedGB := float64(vmStat.Used) / 1024 / 1024 / 1024
    memoryTotalGB := float64(vmStat.Total) / 1024 / 1024 / 1024
    memoryAvailableGB := float64(vmStat.Available) / 1024 / 1024 / 1024
    usedMemoryPercent := (float64(vmStat.Used) / float64(vmStat.Total)) * 100
    freeMemoryPercent := 100 - usedMemoryPercent

    cpuInfo, err := getCPUInfo()
    if err != nil {
        fmt.Printf("Failed to get CPU info: %s\n", err)
        cpuInfo = "Unknown"
    }

    cpuCores, err := getCPUCores()
    if err != nil {
        fmt.Printf("Failed to get CPU cores: %s\n", err)
        cpuCores = 0
    }

    loadAvg, err := load.Avg()
    if err != nil {
        fmt.Printf("Failed to get load average: %s\n", err)
        loadAvg = &load.AvgStat{}
    }

    diskStat, err := disk.Usage("/")
    if err != nil {
        fmt.Printf("Failed to get disk usage: %s\n", err)
        diskStat = &disk.UsageStat{}
    }
    diskFreeGB := float64(diskStat.Free) / 1024 / 1024 / 1024
    diskTotalGB := float64(diskStat.Total) / 1024 / 1024 / 1024
    diskUsedGB := diskTotalGB - diskFreeGB
    freeDiskPercent := (diskFreeGB / diskTotalGB) * 100

    const (
        resetColor = "\033[0m"
        bold       = "\033[1m"
        cyan       = "\033[36m"
    )

    fmt.Println(cyan + "======================================OSMON======================================" + resetColor)
    fmt.Printf("%s%s - CPU.................:%s %s (%d cores)%s\n", bold, cyan, resetColor, cpuInfo, cpuCores, resetColor)
    fmt.Printf("%s%s - Load................:%s %.2f %.2f %.2f%s\n", bold, cyan, resetColor, loadAvg.Load1, loadAvg.Load5, loadAvg.Load15, resetColor)
    fmt.Printf("%s%s - Memory..............:%s %.2f GB / %.2f GB (%.2f GB remaining, %.2f%% free)%s\n",
        bold, cyan, resetColor, memoryUsedGB, memoryTotalGB, memoryAvailableGB, freeMemoryPercent, resetColor)
    fmt.Printf("%s%s - Disk space /........:%s %.2f GB / %.2f GB (%.2f GB remaining, %.2f%% free)%s\n",
        bold, cyan, resetColor, diskUsedGB, diskTotalGB, diskFreeGB, freeDiskPercent, resetColor)
    fmt.Printf("%s%s - Processes...........:%s %d running%s\n", bold, cyan, resetColor, getProcessCount(), resetColor)
    fmt.Printf("%s%s - System uptime.......:%s %d days %d hours %d minutes %d seconds%s\n", bold, cyan, resetColor, upDays, upHours, upMins, upSecs, resetColor)
    fmt.Printf("%s%s - Hostname / IP.......:%s %s / %s%s\n", bold, cyan, resetColor, hostName, ip, resetColor)
    fmt.Printf("%s%s - Release.............:%s %s%s\n", bold, cyan, resetColor, getOSRelease(), resetColor)
    fmt.Printf("%s%s - Current user........:%s [%s: %s]%s\n", bold, cyan, resetColor, currentPTS, userName, resetColor)
    if runtime.GOOS == "linux" {
        fmt.Printf("%s%s - Active users........:%s ", bold, cyan, resetColor)
        for terminal, users := range activeUsers {
            fmt.Printf("[%s: %s] ", terminal, strings.Join(users, ", "))
        }
        fmt.Println()
    }
    fmt.Println(cyan + "=================================================================================" + resetColor)
}

// clearScreen clears the terminal screen.
func clearScreen() {
    if runtime.GOOS == "windows" {
        cmd := exec.Command("cmd", "/c", "cls")
        cmd.Stdout = os.Stdout
        cmd.Run()
    } else {
        fmt.Print("\033[H\033[2J")
    }
}

// displayHelp displays the help information.
func displayHelp() {
    fmt.Println("OSMON - Display system information with refresh interval")
    fmt.Println("\nUsage:")
    fmt.Println("  ./osmon                   Display system information once")
    fmt.Println("  ./osmon -i <interval>     Refresh system information every <interval> seconds")
    fmt.Println("  ./osmon -h/--help         Display this help information")
    fmt.Println("  ./osmon -v/--version      Display the version of the application")
    fmt.Println("\nFlags:")
    fmt.Println("  -i, --interval <interval>  Set interval for refreshing the display in seconds")
    fmt.Println("  -h, --help                 Display help information")
    fmt.Println("  -v, --version              Display the version of the application")
}
