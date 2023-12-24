package main

import (
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"net"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"time"
)

func displaySystemInfoInInterval(interval time.Duration) {
	for {
		clearScreen() // Pozostawione w tej funkcji
		displaySystemInfo()
		time.Sleep(interval)
	}
}

func displaySystemInfo() {
	// Usunięto wywołanie clearScreen() stąd
	currentUser, err := user.Current()
	if err != nil {
		fmt.Printf("Failed to get current user: %s\n", err)
		return
	}
	userName := currentUser.Username

	currentPTS, err := getCurrentPTS()
	if err != nil {
		fmt.Printf("Failed to get current PTS: %s\n", err)
		return
	}

	activeUsers, err := getActiveUsers()
	if err != nil {
		fmt.Printf("Error getting active users: %s\n", err)
		return
	}

	hostName, err := os.Hostname()
	if err != nil {
		fmt.Printf("Failed to get hostname: %s\n", err)
		return
	}

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Printf("Failed to dial: %s\n", err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := localAddr.IP.String()
	conn.Close()

	uptime, err := getUptime()
	if err != nil {
		fmt.Printf("Failed to get uptime: %s\n", err)
		return
	}
	upDays := uptime / (24 * time.Hour)
	upHours := (uptime % (24 * time.Hour)) / time.Hour
	upMins := (uptime % time.Hour) / time.Minute
	upSecs := (uptime % time.Minute) / time.Second

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		fmt.Printf("Failed to get virtual memory info: %s\n", err)
		return
	}
	memoryUsed := vmStat.Used / 1024 / 1024
	memoryTotal := vmStat.Total / 1024 / 1024
	memoryAvailable := vmStat.Available / 1024 / 1024
	usedMemoryPercent := (float64(memoryUsed) / float64(memoryTotal)) * 100
	freeMemoryPercent := 100 - usedMemoryPercent

	cpuInfo, err := getCPUInfo()
	if err != nil {
		fmt.Printf("Failed to get CPU info: %s\n", err)
		return
	}

	cpuCores, err := getCPUCores()
	if err != nil {
		fmt.Printf("Failed to get CPU cores: %s\n", err)
		return
	}

	loadAvg, err := load.Avg()
	if err != nil {
		fmt.Printf("Failed to get load average: %s\n", err)
		return
	}

	diskStat, err := disk.Usage("/")
	if err != nil {
		fmt.Printf("Failed to get disk usage: %s\n", err)
		return
	}

	memoryUsedGB := float64(memoryUsed) / 1024
	memoryTotalGB := float64(memoryTotal) / 1024
	memoryAvailableGB := float64(memoryAvailable) / 1024

	diskFree := float64(diskStat.Free) / 1024 / 1024 / 1024
	diskTotal := float64(diskStat.Total) / 1024 / 1024 / 1024
	diskUsed := diskTotal - diskFree
	freeDiskPercent := (diskFree / diskTotal) * 100

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
	fmt.Printf("%s%s - Disk space /........:%s %.2f GB / %.2f GB (%.2f GB remaining, %.2f%% free)%s\n", bold, cyan, resetColor, diskUsed, diskTotal, diskFree, freeDiskPercent, resetColor)
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

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Print("\033[H\033[2J")
	}
}

func displayHelp() {
	fmt.Println("OSInfo - Display system information with refresh interval")
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
