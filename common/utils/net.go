// Package utils provides common utility functions.
// 包 utils 提供了常用的工具函数。
package utils

import (
	"net"
	"strconv"
	//"time"

	"github.com/turtacn/chasi-bod/common/errors"
)

// IsPortOpen checks if a specific port is open on a given host.
// IsPortOpen 检查给定主机上的特定端口是否打开。
// host: The hostname or IP address. / 主机名或 IP 地址。
// port: The port number. / 端口号。
// timeout: The timeout duration for the connection attempt. / 连接尝试的超时时间。
// Returns true if the port is open, false otherwise, and an error if checking failed.
// 如果端口打开则返回 true，否则返回 false，检查失败时返回错误。
func IsPortOpen(host string, port int, timeout time.Duration) (bool, error) {
	address := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
			// Handle timeout specifically if needed, but for IsPortOpen, timeout means not open within time
			// 如果需要，可以专门处理超时，但对于 IsPortOpen，超时意味着在规定时间内未打开
			return false, nil // Timeout is considered as port not open
		}
		// Connection refused or other network errors are also considered as port not open
		// 连接被拒绝或其他网络错误也被视为端口未打开
		// log the error for debugging, but return false
		GetLogger().Printf("Debug: Connection to %s failed: %v", address, err)
		return false, nil
	}
	defer conn.Close()
	return true, nil // Connection successful, port is open
}

// GetLocalIPs returns a slice of non-loopback local IP addresses.
// GetLocalIPs 返回非回环本地 IP 地址的切片。
// Returns a slice of strings representing IP addresses and an error if retrieval failed.
// 返回表示 IP 地址的字符串切片，以及检索失败时的错误。
func GetLocalIPs() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, errors.NewWithCause(errors.ErrTypeNetwork, "failed to get interface addresses", err)
	}

	var ips []string
	for _, addr := range addrs {
		// Check the address type and if it is not a loopback
		// 检查地址类型以及它是否不是回环地址
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
			// Optionally include IPv6 addresses if needed
			// 如果需要，可以选择包含 IPv6 地址
			// else if ipnet.IP.To16() != nil {
			// 	ips = append(ips, ipnet.IP.String())
			// }
		}
	}

	if len(ips) == 0 {
		return nil, errors.New(errors.ErrTypeNetwork, "no non-loopback IPv4 addresses found")
	}

	return ips, nil
}

// IsValidIPAddress checks if the given string is a valid IP address (IPv4 or IPv6).
// IsValidIPAddress 检查给定的字符串是否为有效的 IP 地址（IPv4 或 IPv6）。
// ip: The string to check. / 要检查的字符串。
// Returns true if it's a valid IP address, false otherwise.
// 如果是有效的 IP 地址则返回 true，否则返回 false。
func IsValidIPAddress(ip string) bool {
	return net.ParseIP(ip) != nil
}

// IsValidHostname checks if the given string is a valid hostname (according to RFC 952/1123).
// IsValidHostname 检查给定的字符串是否为有效的主机名（根据 RFC 952/1123）。
// hostname: The string to check. / 要检查的字符串。
// This is a basic check and might not cover all edge cases.
// 这是一个基本检查，可能不涵盖所有边缘情况。
func IsValidHostname(hostname string) bool {
	// Hostname cannot be empty or too long
	// 主机名不能为空或过长
	if len(hostname) == 0 || len(hostname) > 253 {
		return false
	}
	// Check if it matches the standard hostname pattern (letters, digits, hyphen, dot)
	// 检查它是否匹配标准主机名模式（字母、数字、连字符、点）
	// Labels must be 63 characters or less
	// Labels must start and end with an alphanumeric character
	// Labels can contain hyphens but not at the beginning or end
	// Labels are separated by dots
	// This regex is a simplification but covers common cases.
	// 这个正则表达式是一个简化，但涵盖了常见情况。
	// A more robust check might require a dedicated library or stricter regex.
	// 更强大的检查可能需要专门的库或更严格的正则表达式。
	// Simplified check: does not contain invalid characters or start/end with hyphen
	// 简化检查：不包含无效字符或以连字符开头/结尾
	for _, r := range hostname {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '.') {
			return false
		}
	}
	if hostname[0] == '-' || hostname[len(hostname)-1] == '-' {
		return false
	}
	return true
}
