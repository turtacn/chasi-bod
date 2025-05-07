// Package base provides interfaces and base implementations for building the operating system foundation layer.
// 包 base 提供了构建基础操作系统镜像的接口和基础实现。
package base

import (
	"context"
	"fmt"
	"strings" // Added for string operations // 添加用于字符串操作

	"github.com/turtacn/chasi-bod/common/errors"
	"github.com/turtacn/chasi-bod/pkg/config/model"
)

// OSBuilder defines the interface for building the base operating system image.
// OSBuilder 定义了构建基础操作系统镜像的接口。
// Implementations will handle specific OS distributions (e.g., Debian, CentOS).
// 实现将处理特定的操作系统发行版（例如，Debian、CentOS）。
type OSBuilder interface {
	// PrepareBaseImage downloads or prepares the initial OS image.
	// PrepareBaseImage 下载或准备初始操作系统镜像。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// config: The base OS configuration. / 基础操作系统配置。
	// buildDir: The directory where build artifacts are stored. / 构建 artifact 存储的目录。
	// Returns the path to the root filesystem of the prepared image and an error.
	// 返回准备好的镜像的根文件系统路径和错误。
	PrepareBaseImage(ctx context.Context, config *model.BaseOSConfig, buildDir string) (string, error)

	// InstallPackages installs necessary OS packages into the base image filesystem.
	// InstallPackages 在基础镜像文件系统中安装必要的操作系统软件包。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// config: The base OS configuration. / 基础操作系统配置。
	// rootFS: The path to the root filesystem of the image being built. / 正在构建的镜像的根文件系统路径。
	// Returns an error if installation fails.
	// 如果安装失败则返回错误。
	InstallPackages(ctx context.Context, config *model.BaseOSConfig, rootFS string) error

	// ConfigureSystem sets up basic system configurations (hostname, networking, users, etc.).
	// ConfigureSystem 设置基本系统配置（主机名、网络、用户等）。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// config: The base OS configuration. / 基础操作系统配置。
	// rootFS: The path to the root filesystem of the image being built. / 正在构建的镜像的根文件系统路径。
	// Returns an error if configuration fails.
	// 如果配置失败则返回错误。
	ConfigureSystem(ctx context.Context, config *model.BaseOSConfig, rootFS string) error

	// CustomizeFiles copies custom files and sets permissions.
	// CustomizeFiles 复制自定义文件并设置权限。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// config: The base OS configuration. / 基础操作系统配置。
	// rootFS: The path to the root filesystem of the image being built. / 正在构建的镜像的根文件系统路径。
	// Returns an error if customization fails.
	// 如果自定义失败则返回错误。
	CustomizeFiles(ctx context.Context, config *model.BaseOSConfig, rootFS string) error

	// RunCommands executes custom shell commands within the base image environment.
	// RunCommands 在基础镜像环境中执行自定义 shell 命令。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// config: The base OS configuration. / 基础操作系统配置。
	// rootFS: The path to the root filesystem of the image being built. / 正在构建的镜像的根文件系统路径。
	// Returns an error if command execution fails.
	// 如果命令执行失败则返回错误。
	RunCommands(ctx context.Context, config *model.BaseOSConfig, rootFS string) error

	// Cleanup performs any necessary cleanup after building the base OS.
	// Cleanup 在构建基础操作系统后执行任何必要的清理工作。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// rootFS: The path to the root filesystem. / 根文件系统的路径。
	// Returns an error if cleanup fails.
	// 如果清理失败则返回错误。
	Cleanup(ctx context.Context, rootFS string) error
}

// NewOSBuilder creates a new OSBuilder implementation based on the configuration.
// NewOSBuilder 根据配置创建一个新的 OSBuilder 实现。
// config: The base OS configuration. / 基础操作系统配置。
// Returns an OSBuilder implementation or an error if the OS type is unsupported.
// 返回 OSBuilder 实现，如果操作系统类型不受支持则返回错误。
func NewOSBuilder(config *model.BaseOSConfig) (OSBuilder, error) {
	// Determine the OS type from config.Image or a dedicated field
	// 从 config.Image 或专用字段确定操作系统类型
	// For now, let's do a basic check
	// 现在，让我们做一个基本检查
	imageLower := strings.ToLower(config.Image)
	if strings.Contains(imageLower, "ubuntu") || strings.Contains(imageLower, "debian") {
		// return NewDebianOSBuilder() // Assuming a Debian/Ubuntu builder exists
		return nil, errors.New(errors.ErrTypeNotImplemented, "Debian/Ubuntu OS builder not implemented yet")
	} else if strings.Contains(imageLower, "centos") || strings.Contains(imageLower, "rhel") || strings.Contains(imageLower, "rockylinux") || strings.Contains(imageLower, "almalinux") {
		// return NewRedhatOSBuilder() // Assuming a RedHat-based builder exists
		return nil, errors.New(errors.ErrTypeNotImplemented, "RedHat-based OS builder not implemented yet")
	} else {
		return nil, errors.New(errors.ErrTypeValidation, fmt.Sprintf("unsupported base OS image '%s'", config.Image))
	}
}

// Example implementation structure (not fully functional)
// 示例实现结构（不完全功能）
// type DebianOSBuilder struct{}
//
// func (b *DebianOSBuilder) PrepareBaseImage(ctx context.Context, config *model.BaseOSConfig, buildDir string) (string, error) {
// 	utils.GetLogger().Printf("Placeholder: Preparing Debian base image %s into %s", config.Image, buildDir)
// 	// Use debootstrap or docker export to get a base filesystem
// 	rootFSPath := filepath.Join(buildDir, "rootfs")
// 	// Create dummy rootfs dir
// 	if err := utils.MkdirAll(rootFSPath, 0755); err != nil {
// 		return "", err
// 	}
// 	return rootFSPath, errors.New(errors.ErrTypeNotImplemented, "Debian base image preparation not implemented")
// }
//
// func (b *DebianOSBuilder) InstallPackages(ctx context.Context, config *model.BaseOSConfig, rootFS string) error {
// 	utils.GetLogger().Printf("Placeholder: Installing packages into %s", rootFS)
// 	// Use chroot and apt-get to install packages
// 	return errors.New(errors.ErrTypeNotImplemented, "Debian package installation not implemented")
// }
// // ... implement other methods
