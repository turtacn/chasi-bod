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

// StorageConfigPhase defines the interface for configuring storage on target nodes.
// StorageConfigPhase 定义了在目标节点上配置存储的接口。
// This includes mounting disks, setting up filesystems, and potentially configuring CSI driver dependencies.
// 这包括挂载磁盘、设置文件系统，以及可能配置 CSI 驱动依赖项。
type StorageConfigPhase interface {
	// Run executes the storage configuration phase for a single node.
	// Run 为单个节点执行存储配置阶段。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// nodeCfg: The configuration for the specific node. / 特定节点的配置。
	// clusterCfg: The overall cluster configuration including storage details. / 整体集群配置，包括存储详细信息。
	// Returns an error if the phase fails for the node.
	// 如果阶段在节点上失败则返回错误。
	Run(ctx context.Context, nodeCfg *model.NodeConfig, clusterCfg *model.ClusterConfig) error
}

// NewStorageConfigPhase creates a new StorageConfigPhase instance.
// NewStorageConfigPhase 创建一个新的 StorageConfigPhase 实例。
// Returns a StorageConfigPhase implementation.
// 返回 StorageConfigPhase 实现。
func NewStorageConfigPhase() StorageConfigPhase {
	return &defaultStorageConfigPhase{}
}

// defaultStorageConfigPhase is a default implementation of the StorageConfigPhase.
// defaultStorageConfigPhase 是 StorageConfigPhase 的默认实现。
type defaultStorageConfigPhase struct{}

// Run executes the storage configuration phase.
// Run 执行存储配置阶段。
func (p *defaultStorageConfigPhase) Run(ctx context.Context, nodeCfg *model.NodeConfig, clusterCfg *model.ClusterConfig) error {
	utils.GetLogger().Printf("Running StorageConfigPhase for node %s", nodeCfg.Address)

	// TODO: Get SSH client for the node
	// TODO: 获取节点的 SSH 客户端
	// sshClient, err := sshutil.NewClient(...)
	// if err != nil { return errors.NewWithCause(...) }
	// defer sshClient.Close()

	// Step 1: Configure node-specific disks (partitioning, formatting, mounting) - Handled in OSConfigPhase
	// 步骤 1：配置节点特定磁盘（分区、格式化、挂载）- 在 OSConfigPhase 中处理
	// Reiterate the importance if disk configs are present
	// 如果存在磁盘配置，重申其重要性
	if len(nodeCfg.DiskConfigs) > 0 {
		utils.GetLogger().Printf("Disk configuration (partitioning, formatting, mounting) was handled in OSConfigPhase for node %s.", nodeCfg.Address)
	} else {
		utils.GetLogger().Printf("No node-specific disk configurations defined for node %s.", nodeCfg.Address)
	}

	// Step 2: Install any node-specific dependencies for the chosen CSI driver (e.g., iscsi-initiator-utils, ceph-common)
	// 步骤 2：安装所选 CSI 驱动的任何节点特定依赖项（例如，iscsi-initiator-utils、ceph-common）
	// This depends on the clusterCfg.Storage.StorageClasses and their provisioners.
	// This requires identifying the OS and using its package manager remotely.
	// 这取决于 clusterCfg.Storage.StorageClasses 及其 provisioners。
	// 这需要识别操作系统并远程使用其包管理器。
	utils.GetLogger().Printf("Installing CSI driver dependencies on %s...", nodeCfg.Address)
	// TODO: Based on required CSI drivers from clusterCfg.Storage, determine necessary OS packages and install them via SSH (e.g., apt-get install, yum install)
	// TODO: 根据 clusterCfg.Storage 中所需的 CSI 驱动，确定必要的操作系统软件包并通过 SSH 进行安装（例如，apt-get install, yum install）
	utils.GetLogger().Printf("Placeholder: CSI driver dependencies installed on %s.", nodeCfg.Address)

	// Step 3: Configure required storage mounts or settings on the node for the CSI driver (if needed)
	// 步骤 3：在节点上为 CSI 驱动配置所需的存储挂载或设置（如果需要）
	// Some CSI drivers might require specific mounts, kernel modules, or configurations on the host nodes.
	// This depends heavily on the specific CSI driver.
	// 某些 CSI 驱动可能需要在主机节点上进行特定的挂载、内核模块或配置。
	// 这高度依赖于特定的 CSI 驱动。
	utils.GetLogger().Printf("Configuring node storage settings for CSI driver on %s...", nodeCfg.Address)
	// TODO: Implement remote configuration via SSH based on the chosen CSI driver
	// TODO: 根据所选的 CSI 驱动，通过 SSH 实现远程配置
	utils.GetLogger().Printf("Placeholder: Node storage settings configured for CSI driver on %s.", nodeCfg.Address)

	// Step 4: Verify storage setup (optional)
	// 步骤 4：验证存储设置（可选）
	utils.GetLogger().Printf("Verifying storage setup on %s...", nodeCfg.Address)
	// This might involve checking mount points, disk usage, or CSI node driver registration (if CSI node agent is deployed).
	// 这可能涉及检查挂载点、磁盘使用情况或 CSI 节点驱动程序注册（如果部署了 CSI 节点代理）。
	utils.GetLogger().Printf("Placeholder: Storage setup verified on %s.", nodeCfg.Address)

	utils.GetLogger().Printf("StorageConfigPhase completed successfully for node %s", nodeCfg.Address)
	return nil
}

// TODO: Implement sshutil package for remote command execution and file transfer
// TODO: 实现 sshutil 包用于远程命令执行和文件传输
// TODO: Implement helper functions to identify required OS packages for different CSI drivers
// TODO: 实现为不同 CSI 驱动识别所需操作系统软件包的辅助函数
