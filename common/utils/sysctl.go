// Package utils provides common utility functions.
// 包 utils 提供了常用的工具函数。
package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/turtacn/chasi-bod/common/errors"
)

// sysctlBaseDir is the base directory for sysctl entries in /proc.
// sysctlBaseDir 是 /proc 中 sysctl 条目的基本目录。
const sysctlBaseDir = "/proc/sys"

// ReadSysctl reads the value of a kernel parameter.
// ReadSysctl 读取内核参数的值。
// param: The sysctl parameter name (e.g., "net.ipv4.ip_forward"). / sysctl 参数名称（例如，“net.ipv4.ip_forward”）。
// Returns the parameter value as a string and an error if reading failed.
// 返回参数值作为字符串，以及读取失败时的错误。
func ReadSysctl(param string) (string, error) {
	filePath := filepath.Join(sysctlBaseDir, strings.ReplaceAll(param, ".", "/"))
	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", errors.New(errors.ErrTypeNotFound, fmt.Sprintf("sysctl parameter %s not found", param))
		}
		return "", errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to read sysctl parameter %s from %s", param, filePath), err)
	}
	// Remove leading/trailing whitespace including newlines
	// 移除前导/尾随的空白字符，包括换行符
	return strings.TrimSpace(string(content)), nil
}

// SetSysctl sets the value of a kernel parameter. Requires appropriate permissions (usually root).
// SetSysctl 设置内核参数的值。需要适当的权限（通常是 root）。
// param: The sysctl parameter name (e.g., "net.ipv4.ip_forward"). / sysctl 参数名称（例如，“net.ipv4.ip_forward”）。
// value: The desired parameter value as a string. / 期望的参数值作为字符串。
// Returns an error if setting failed.
// 设置失败时返回错误。
func SetSysctl(param, value string) error {
	filePath := filepath.Join(sysctlBaseDir, strings.ReplaceAll(param, ".", "/"))
	// Open the file with write permissions
	// 以写入权限打开文件
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644) // Use 0644 as default permissions, adjust if needed
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New(errors.ErrTypeNotFound, fmt.Sprintf("sysctl parameter %s not found for writing", param))
		}
		// Check for permission denied error specifically
		// 专门检查权限拒绝错误
		if os.IsPermission(err) {
			return errors.New(errors.ErrTypeSystem, fmt.Sprintf("permission denied to set sysctl parameter %s. Requires root privileges.", param))
		}
		return errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to open sysctl file %s for writing", filePath), err)
	}
	defer file.Close()

	// Write the new value followed by a newline
	// 写入新值后跟一个换行符
	_, err = file.WriteString(value + "\n")
	if err != nil {
		return errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to write value %q to sysctl file %s", value, filePath), err)
	}

	// Ensure the changes are flushed to disk
	// 确保更改刷新到磁盘
	err = file.Sync()
	if err != nil {
		GetLogger().Printf("Warning: Failed to sync sysctl file %s after writing: %v", filePath, err)
		// Continue, as write might have succeeded partially
		// 继续，因为写入可能已部分成功
	}

	return nil
}

// ApplySysctlConfig applies a map of sysctl parameters. Requires appropriate permissions.
// ApplySysctlConfig 应用 sysctl 参数的映射。需要适当的权限。
// config: A map where keys are sysctl parameter names and values are desired string values. / 键为 sysctl 参数名称，值为期望的字符串值的映射。
// Returns an error if any parameter failed to set.
// 如果任何参数设置失败，则返回错误。
func ApplySysctlConfig(config map[string]string) error {
	if len(config) == 0 {
		GetLogger().Println("No sysctl parameters to apply.")
		return nil
	}

	var failedParams []string
	for param, value := range config {
		GetLogger().Printf("Applying sysctl %s = %q", param, value)
		if err := SetSysctl(param, value); err != nil {
			GetLogger().Printf("Error applying sysctl %s = %q: %v", param, value, err)
			failedParams = append(failedParams, param)
			// Continue to try applying other parameters
			// 继续尝试应用其他参数
		}
	}

	if len(failedParams) > 0 {
		return errors.New(errors.ErrTypeSystem, fmt.Sprintf("failed to apply sysctl parameters: %s", strings.Join(failedParams, ", ")))
	}

	GetLogger().Println("Successfully applied all sysctl parameters.")
	return nil
}
