// Package phases defines the individual steps involved in deploying the chasi-bod platform.
// 包 phases 定义了部署 chasi-bod 平台的各个步骤。
package phases

import (
	"context"
	"github.com/turtacn/chasi-bod/common/utils" // Assuming logger is here // 假设日志记录器在这里
	"github.com/turtacn/chasi-bod/pkg/config/model"
	// Assuming you have an SSH client utility
	// 假设您有一个 SSH 客户端工具
	// "github.com/turtacn/chasi-bod/pkg/sshutil"
)

// RuntimeConfigPhase defines the interface for configuring and starting the container runtime.
// RuntimeConfigPhase 定义了配置和启动容器运行时的接口。
// This phase assumes the runtime binaries are already present in the image (installed by builder).
// 此阶段假设运行时二进制文件已存在于镜像中（由 builder 安装）。
type RuntimeConfigPhase interface {
	// Run executes the container runtime configuration phase for a single node.
	// Run 为单个节点执行容器运行时配置阶段。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// nodeCfg: The configuration for the specific node. / 特定节点的配置。
	// clusterCfg: The overall cluster configuration. / 整体集群配置。
	// Returns an error if the phase fails for the node.
	// 如果阶段在节点上失败则返回错误。
	Run(ctx context.Context, nodeCfg *model.NodeConfig, clusterCfg *model.ClusterConfig) error
}

// NewRuntimeConfigPhase creates a new RuntimeConfigPhase instance.
// NewRuntimeConfigPhase 创建一个新的 RuntimeConfigPhase 实例。
// Returns a RuntimeConfigPhase implementation.
// 返回 RuntimeConfigPhase 实现。
func NewRuntimeConfigPhase() RuntimeConfigPhase {
	return &defaultRuntimeConfigPhase{}
}

// defaultRuntimeConfigPhase is a default implementation of the RuntimeConfigPhase.
// defaultRuntimeConfigPhase 是 RuntimeConfigPhase 的默认实现。
type defaultRuntimeConfigPhase struct{}

// Run executes the container runtime configuration phase.
// Run 执行容器运行时配置阶段。
func (p *defaultRuntimeConfigPhase) Run(ctx context.Context, nodeCfg *model.NodeConfig, clusterCfg *model.ClusterConfig) error {
	utils.GetLogger().Printf("Running RuntimeConfigPhase for node %s (Runtime: %s)", nodeCfg.Address, clusterCfg.ContainerRuntime)

	// TODO: Get SSH client for the node
	// TODO: 获取节点的 SSH 客户端
	// sshClient, err := sshutil.NewClient(...)
	// if err != nil { return errors.NewWithCause(...) }
	// defer sshClient.Close()

	// Step 1: Generate or update container runtime configuration file
	// 步骤 1：生成或更新容器运行时配置文件
	utils.GetLogger().Printf("Generating runtime configuration for %s on %s...", clusterCfg.ContainerRuntime, nodeCfg.Address)
	// This might involve:
	// - Setting cgroup driver (should match Kubelet)
	// - Configuring registries, proxies, etc.
	// This config file template should ideally be in the image, and this phase applies node-specific values or updates.
	// 这可能涉及：
	// - 设置 cgroup 驱动（应与 Kubelet 匹配）
	// - 配置注册表、代理等。
	// 这个配置文件模板理想情况下应该在镜像中，此阶段应用节点特定的值或更新。
	// TODO: Implement remote config file generation/update via SSH
	// Need to determine config file path based on runtime (e.g., /etc/containerd/config.toml)
	// 需要根据运行时确定配置文件路径（例如，/etc/containerd/config.toml）
	// Might need to use a template engine remotely or generate locally and copy.
	// 可能需要远程使用模板引擎或在本地生成后复制。
	utils.GetLogger().Printf("Placeholder: Runtime configuration generated for %s on %s.", clusterCfg.ContainerRuntime, nodeCfg.Address)

	// Step 2: Reload systemd manager configuration
	// 步骤 2：重新加载 systemd 管理器配置
	utils.GetLogger().Printf("Reloading systemd daemon on %s...", nodeCfg.Address)
	// TODO: sshClient.RunCommand(ctx, "sudo systemctl daemon-reload")
	utils.GetLogger().Printf("Placeholder: Systemd daemon reloaded on %s.", nodeCfg.Address)

	// Step 3: Enable and start the container runtime service
	// 步骤 3：启用并启动容器运行时服务
	runtimeServiceName := clusterCfg.ContainerRuntime // Assuming service name matches runtime name (e.g., containerd, crio, docker)
	utils.GetLogger().Printf("Enabling and starting %s service on %s...", runtimeServiceName, nodeCfg.Address)
	// TODO: sshClient.RunCommand(ctx, fmt.Sprintf("sudo systemctl enable %s", runtimeServiceName))
	// TODO: sshClient.RunCommand(ctx, fmt.Sprintf("sudo systemctl start %s", runtimeServiceName))
	utils.GetLogger().Printf("Placeholder: %s service enabled and started on %s.", runtimeServiceName, nodeCfg.Address)

	// Step 4: Wait for the runtime service to be healthy
	// 步骤 4：等待运行时服务健康
	utils.GetLogger().Printf("Waiting for %s service to be healthy on %s...", runtimeServiceName, nodeCfg.Address)
	// TODO: Implement health check (e.g., check socket file existence, run `crictl info` or `docker info` via SSH)
	// TODO: 实现健康检查（例如，检查 socket 文件是否存在，通过 SSH 运行 `crictl info` 或 `docker info`）
	// Example: checkCmd := "sudo crictl info"
	// err = sshClient.RunCommandUntilSuccess(ctx, checkCmd, 10*time.Second, 1*time.Second) // Assuming a helper exists
	// if err != nil { return errors.NewWithCause(...) }
	utils.GetLogger().Printf("Placeholder: %s service is healthy on %s.", runtimeServiceName, nodeCfg.Address)

	utils.GetLogger().Printf("RuntimeConfigPhase completed successfully for node %s", nodeCfg.Address)
	return nil
}

// TODO: Implement helper functions for remote command execution, file transfer, etc. using sshutil
// TODO: 使用 sshutil 实现远程命令执行、文件传输等辅助函数
