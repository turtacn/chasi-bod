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

// OSConfigPhase defines the interface for configuring the operating system on target nodes.
// OSConfigPhase 定义了在目标节点上配置操作系统的接口。
// This phase applies node-specific OS configurations like sysctl, disk mounts, user setup not handled during image build.
// 此阶段应用节点特定的操作系统配置，例如 sysctl、磁盘挂载、在镜像构建期间未处理的用户设置。
type OSConfigPhase interface {
	// Run executes the OS configuration phase for a single node.
	// Run 为单个节点执行操作系统配置阶段。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// nodeCfg: The configuration for the specific node. / 特定节点的配置。
	// clusterCfg: The overall cluster configuration. / 整体集群配置。
	// Returns an error if the phase fails for the node.
	// 如果阶段在节点上失败则返回错误。
	Run(ctx context.Context, nodeCfg *model.NodeConfig, clusterCfg *model.ClusterConfig) error
}

// NewOSConfigPhase creates a new OSConfigPhase instance.
// NewOSConfigPhase 创建一个新的 OSConfigPhase 实例。
// Returns an OSConfigPhase implementation.
// 返回 OSConfigPhase 实现。
func NewOSConfigPhase() OSConfigPhase {
	return &DefaultOSConfigPhase{}
}

// DefaultOSConfigPhase is a default implementation of the OSConfigPhase.
// DefaultOSConfigPhase 是 OSConfigPhase 的默认实现。
type DefaultOSConfigPhase struct{}

// Run executes the OS configuration phase.
// Run 执行操作系统配置阶段。
func (p *DefaultOSConfigPhase) Run(ctx context.Context, nodeCfg *model.NodeConfig, clusterCfg *model.ClusterConfig) error {
	utils.GetLogger().Printf("Running OSConfigPhase for node %s", nodeCfg.Address)

	// TODO: Get SSH client for the node
	// TODO: 获取节点的 SSH 客户端
	// sshClient, err := sshutil.NewClient(...)
	// if err != nil { return errors.NewWithCause(...) }
	// defer sshClient.Close()

	// Step 1: Apply node-specific Sysctl configurations
	// 步骤 1：应用节点特定的 Sysctl 配置
	if len(nodeCfg.SysctlConfig) > 0 {
		utils.GetLogger().Printf("Applying node-specific sysctl configurations on %s...", nodeCfg.Address)
		// TODO: Use sshClient to execute sysctl commands remotely
		// TODO: 使用 sshClient 远程执行 sysctl 命令
		// for param, value := range nodeCfg.SysctlConfig {
		// 	cmd := fmt.Sprintf("sudo sysctl -w %s=\"%s\"", param, value)
		// 	_, err := sshClient.RunCommand(ctx, cmd)
		// 	if err != nil {
		// 		utils.GetLogger().Printf("Warning: Failed to set sysctl %s=%s on %s: %v", param, value, nodeCfg.Address, err)
		// 		// Decide if this is a fatal error or just a warning
		// 		// 决定这是否是致命错误或只是警告
		// 	}
		// }
		utils.GetLogger().Printf("Placeholder: Applied node-specific sysctl configurations on %s.", nodeCfg.Address)

		// Optionally, persist sysctl changes (e.g., to /etc/sysctl.d/)
		// 可选地，持久化 sysctl 更改（例如，到 /etc/sysctl.d/）
		utils.GetLogger().Printf("Placeholder: Persisting node-specific sysctl configurations on %s.", nodeCfg.Address)
	} else {
		utils.GetLogger().Printf("No node-specific sysctl configurations to apply on %s.", nodeCfg.Address)
	}

	// Step 2: Configure disks and mount points
	// 步骤 2：配置磁盘和挂载点
	if len(nodeCfg.DiskConfigs) > 0 {
		utils.GetLogger().Printf("Configuring disks on %s...", nodeCfg.Address)
		for _, diskCfg := range nodeCfg.DiskConfigs {
			utils.GetLogger().Printf("Configuring disk %s (MountPoint: %s, Filesystem: %s, Format: %t) on %s",
				diskCfg.Device, diskCfg.MountPoint, diskCfg.Filesystem, diskCfg.Format, nodeCfg.Address)
			// TODO: Implement remote disk partitioning, formatting, and mounting via SSH
			// TODO: 通过 SSH 实现远程磁盘分区、格式化和挂载
			// This is OS-dependent and requires tools like fdisk, parted, mkfs, mount.
			// 这与操作系统相关，需要 fdisk, parted, mkfs, mount 等工具。
			// Example:
			// if diskCfg.Format { // Format the disk if requested
			//     formatCmd := fmt.Sprintf("sudo mkfs.%s %s", diskCfg.Filesystem, diskCfg.Device)
			//     _, err := sshClient.RunCommand(ctx, formatCmd)
			//     if err != nil { return errors.NewWithCause(...) }
			// }
			// if diskCfg.MountPoint != "" { // Create mount point and mount
			//     mkdirCmd := fmt.Sprintf("sudo mkdir -p %s", diskCfg.MountPoint)
			//     _, err := sshClient.RunCommand(ctx, mkdirCmd)
			//     if err != nil { return errors.NewWithCause(...) }
			//     mountCmd := fmt.Sprintf("sudo mount %s %s", diskCfg.Device, diskCfg.MountPoint)
			//     _, err := sshClient.RunCommand(ctx, mountCmd)
			//     if err != nil { return errors.NewWithCause(...) }
			//     // Add fstab entry to make it persistent
			//     fstabEntry := fmt.Sprintf("%s %s %s defaults 0 0", diskCfg.Device, diskCfg.MountPoint, diskCfg.Filesystem)
			//     echoCmd := fmt.Sprintf("echo '%s' | sudo tee -a /etc/fstab", fstabEntry) // Use tee with sudo to append to a root-owned file
			//     _, err := sshClient.RunCommand(ctx, echoCmd)
			//     if err != nil { return errors.NewWithCause(...) }
			// }
			utils.GetLogger().Printf("Placeholder: Disk %s configured on %s.", diskCfg.Device, nodeCfg.Address)
		}
		utils.GetLogger().Printf("Disks configured successfully on %s.", nodeCfg.Address)
	} else {
		utils.GetLogger().Printf("No disk configurations to apply on %s.", nodeCfg.Address)
	}

	// Step 3: Ensure required users and groups exist and have correct permissions (if not handled by image build)
	// 步骤 3：确保所需用户和组存在并具有正确的权限（如果在镜像构建中未处理）
	// Users defined in BaseOSConfig should ideally be handled during image build.
	// This step might be for runtime user creation or modification based on node-specific needs.
	// Users defined in BaseOSConfig 理想情况下应在镜像构建期间处理。
	// 此步骤可能用于根据节点特定需求在运行时创建或修改用户。
	// utils.GetLogger().Printf("Configuring users and groups on %s...", nodeCfg.Address)
	// TODO: Implement remote user/group creation/modification via SSH using useradd, groupadd, usermod commands
	// TODO: 通过 SSH 使用 useradd, groupadd, usermod 命令实现远程用户/组创建/修改
	// For files defined in config.BaseOSConfig.Files, ensure correct permissions are set.
	// This might be done during CustomizeFiles in the builder, but re-checking/setting here is also an option.
	// 对于 config.BaseOSConfig.Files 中定义的文件，确保设置了正确的权限。
	// 这可以在 builder 的 CustomizeFiles 中完成，但在此处重新检查/设置也是一个选项。
	// for _, fileCfg := range clusterCfg.BaseOS.Files {
	// 	if fileCfg.Mode != "" {
	// 		// Convert octal string to file mode if needed
	// 		mode, err := strconv.ParseUint(fileCfg.Mode, 8, 32)
	// 		if err != nil {
	// 			utils.GetLogger().Printf("Warning: Invalid file permission mode '%s' for %s on %s: %v", fileCfg.Mode, fileCfg.Dest, nodeCfg.Address, err)
	// 			continue // Skip if mode is invalid
	// 		}
	// 		chmodCmd := fmt.Sprintf("sudo chmod %o %s", mode, fileCfg.Dest) // Use %o for octal format
	// 		_, err = sshClient.RunCommand(ctx, chmodCmd)
	// 		if err != nil {
	// 			utils.GetLogger().Printf("Warning: Failed to set permissions for %s on %s: %v", fileCfg.Dest, nodeCfg.Address, err)
	// 		}
	// 	}
	// 	// TODO: Handle chown if user/group is specified
	// }
	utils.GetLogger().Printf("Placeholder: Users and groups configured and file permissions set on %s.", nodeCfg.Address)

	utils.GetLogger().Printf("OSConfigPhase completed successfully for node %s", nodeCfg.Address)
	return nil
}
