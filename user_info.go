package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

func getActiveUsers() (map[string][]string, error) {
	out, err := exec.Command("ps", "aux").Output()
	if err != nil {
		return nil, err
	}

	activeUsers := make(map[string][]string)
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 1 {
			user := fields[0]
			tty := fields[6]
			if strings.HasPrefix(tty, "pts") { // Dodano sprawdzenie czy tty zaczyna się od "pts"
				if _, found := activeUsers[tty]; !found {
					activeUsers[tty] = []string{}
				}
				fullUserName, err := getFullUserName(user)
				if err == nil {
					user = fullUserName
				}
				if !contains(activeUsers[tty], user) {
					activeUsers[tty] = append(activeUsers[tty], user)
				}
			}
		}
	}
	return activeUsers, nil
}

func getFullUserName(shortName string) (string, error) {
	file, err := os.Open("/etc/passwd")
	if err != nil {
		return shortName, err
	}
	defer file.Close()

	trimmedShortName := strings.TrimSuffix(shortName, "+")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) > 0 {
			userName := parts[0]
			if strings.HasPrefix(userName, trimmedShortName) {
				return userName, nil
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return shortName, err
	}
	return shortName, nil
}

func getCurrentPTS() (string, error) {
	if runtime.GOOS == "linux" {
		tty, err := os.Readlink("/proc/self/fd/0")
		if err != nil {
			return "", err
		}
		re := regexp.MustCompile(`(pts/\d+|ttyS?\d*)`)
		pts := re.FindString(tty)
		if pts == "" {
			return "", fmt.Errorf("could not parse PTS")
		}
		return pts, nil
	} else if runtime.GOOS == "darwin" {
		ttyNumber := os.Getenv("_P9K_SSH_TTY")
		re := regexp.MustCompile(`(ttys?\d*)`)
		pts := re.FindString(ttyNumber)
		if ttyNumber == "" {
			return "not available", nil
		}
		return pts, nil
	}
	return "", fmt.Errorf("unsupported platform")
}
