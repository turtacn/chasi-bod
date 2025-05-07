// Package runtime provides interfaces and implementations for integrating a container runtime into the platform image.
// 包 runtime 提供了将容器运行时集成到平台镜像的接口和实现。
package runtime

import (
	"context"
	"fmt"
	"strings"

	"github.com/turtacn/chasi-bod/common/errors"
	"github.com/turtacn/chasi-bod/pkg/config/model"
)

// RuntimeInstaller defines the interface for installing and configuring a container runtime.
// RuntimeInstaller 定义了安装和配置容器运行时的接口。
// Implementations will handle specific runtimes (e.g., containerd, cri-o, docker).
// 实现将处理特定的运行时（例如，containerd、cri-o、docker）。
type RuntimeInstaller interface {
	// Install installs the container runtime binaries and dependencies into the root filesystem.
	// Install 将容器运行时二进制文件和依赖项安装到根文件系统中。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// config: The cluster configuration (to get runtime name and version). / 集群配置（获取运行时名称和版本）。
	// rootFS: The path to the root filesystem of the image being built. / 正在构建的镜像的根文件系统路径。
	// Returns an error if installation fails.
	// 如果安装失败则返回错误。
	Install(ctx context.Context, config *model.ClusterConfig, rootFS string) error

	// Configure sets up the container runtime configuration (e.g., cgroup driver, registries, endpoints).
	// Configure 设置容器运行时配置（例如，cgroup 驱动、注册表、端点）。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// config: The cluster configuration. / 集群配置。
	// rootFS: The path to the root filesystem of the image being built. / 正在构建的镜像的根文件系统路径。
	// Returns an error if configuration fails.
	// 如果配置失败则返回错误。
	Configure(ctx context.Context, config *model.ClusterConfig, rootFS string) error

	// EnableService ensures the runtime service is configured to start on boot.
	// EnableService 确保运行时服务配置为在启动时启动。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// config: The cluster configuration. / 集群配置。
	// rootFS: The path to the root filesystem of the image being built. / 正在构建的镜像的根文件系统路径。
	// Returns an error if enabling fails.
	// 如果启用失败则返回错误。
	EnableService(ctx context.Context, config *model.ClusterConfig, rootFS string) error

	// PreloadImages pulls any necessary base container images (optional).
	// PreloadImages 拉取任何必要的基础容器镜像（可选）。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// config: The cluster configuration. / 集群配置。
	// rootFS: The path to the root filesystem of the image being built. / 正在构建的镜像的根文件系统路径。
	// Returns an error if preloading fails.
	// 如果预加载失败则返回错误。
	PreloadImages(ctx context.Context, config *model.ClusterConfig, rootFS string) error
}

// NewRuntimeInstaller creates a new RuntimeInstaller implementation based on the configuration.
// NewRuntimeInstaller 根据配置创建一个新的 RuntimeInstaller 实现。
// config: The cluster configuration. / 集群配置。
// Returns a RuntimeInstaller implementation or an error if the runtime is unsupported.
// 返回 RuntimeInstaller 实现，如果运行时不受支持则返回错误。
func NewRuntimeInstaller(config *model.ClusterConfig) (RuntimeInstaller, error) {
	runtimeLower := strings.ToLower(config.ContainerRuntime)
	switch runtimeLower {
	case "containerd":
		// return &ContainerdInstaller{} // Assuming a ContainerdInstaller exists
		return nil, errors.New(errors.ErrTypeNotImplemented, "containerd installer not implemented yet")
	case "cri-o":
		// return &CrioInstaller{} // Assuming a CrioInstaller exists
		return nil, errors.New(errors.ErrTypeNotImplemented, "cri-o installer not implemented yet")
	case "docker":
		// return &DockerInstaller{} // Assuming a DockerInstaller exists (less common for modern K8s)
		return nil, errors.New(errors.ErrTypeNotImplemented, "docker installer not implemented yet")
	default:
		return nil, errors.New(errors.ErrTypeValidation, fmt.Sprintf("unsupported container runtime '%s'", config.ContainerRuntime))
	}
}

// Example implementation structure (not fully functional)
// 示例实现结构（不完全功能）
// type ContainerdInstaller struct{}
//
// func (i *ContainerdInstaller) Install(ctx context.Context, config *model.ClusterConfig, rootFS string) error {
// 	utils.GetLogger().Printf("Placeholder: Installing containerd into %s", rootFS)
// 	// Download and place containerd binaries and systemd service file in rootFS
// 	return errors.New(errors.ErrTypeNotImplemented, "containerd installation not implemented")
// }
//
// func (i *ContainerdInstaller) Configure(ctx context.Context, config *model.ClusterConfig, rootFS string) error {
// 	utils.GetLogger().Printf("Placeholder: Configuring containerd in %s", rootFS)
// 	// Generate and place containerd config.toml in rootFS
// 	return errors.New(errors.ErrTypeNotImplemented, "containerd configuration not implemented")
// }
//
// func (i *ContainerdInstaller) EnableService(ctx context.Context, config *model.ClusterConfig, rootFS string) error {
// 	utils.GetLogger().Printf("Placeholder: Enabling containerd service in %s", rootFS)
// 	// Enable systemd service using chroot or similar
// 	return errors.New(errors.ErrTypeNotImplemented, "containerd service enabling not implemented")
// }
//
// func (i *ContainerdInstaller) PreloadImages(ctx context.Context, config *model.ClusterConfig, rootFS string) error {
// 	utils.GetLogger().Printf("Placeholder: Preloading images for containerd in %s", rootFS)
// 	// Use the runtime's CLI (crictl or containerd) inside the chroot environment to pull images
// 	return errors.New(errors.ErrTypeNotImplemented, "containerd image preloading not implemented")
// }
