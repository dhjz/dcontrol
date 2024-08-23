package utils

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return strconv.FormatUint(bytes, 10) + " B"
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func FormatDuration(duration time.Duration) string {
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	return fmt.Sprintf("%d天%d小时%d分钟", days, hours, minutes)
}

func GetAllIPs() string {
	ips := []string{}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
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

		ip = ip.To4()
		if ip != nil && !strings.HasPrefix(ip.String(), "169") && !strings.HasSuffix(ip.String(), ".1") {
			ips = append(ips, ip.String())
		}
	}

	return strings.Join(ips, ", ")
}

func GetMainIP() string {
	var mainIp = ""
	var otherIp = ""
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
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

		ip = ip.To4()
		if ip != nil && !strings.HasPrefix(ip.String(), "169") && !strings.HasSuffix(ip.String(), ".1") {
			if strings.HasPrefix(ip.String(), "192") {
				mainIp = ip.String()
			} else {
				otherIp = ip.String()
			}
		}
	}

	if mainIp != "" {
		return mainIp
	}
	return otherIp
}
