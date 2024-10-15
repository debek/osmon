package main

import (
    "fmt"
    "net"
)

// getLocalIP retrieves the local IP address without connecting to an external server.
func getLocalIP() (string, error) {
    interfaces, err := net.Interfaces()
    if err != nil {
        return "", err
    }
    for _, iface := range interfaces {
        if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
            addrs, err := iface.Addrs()
            if err != nil {
                continue
            }
            for _, addr := range addrs {
                var ip net.IP
                switch v := addr.(type) {
                case *net.IPNet:
                    ip = v.IP
                case *net.IPAddr:
                    ip = v.IP
                }
                if ip == nil || ip.IsLoopback() {
                    continue
                }
                if ip.To4() != nil {
                    return ip.String(), nil
                }
            }
        }
    }
    return "", fmt.Errorf("No network connection available")
}
