// Package lifecycle provides functionality for managing the full lifecycle of the chasi-bod platform and hosted applications.
// 包 lifecycle 提供了管理 chasi-bod 平台和托管应用程序完整生命周期的功能。
// This includes platform upgrade, scaling, backup, restore, etc.
// 这包括平台升级、扩缩容、备份、恢复等。
package lifecycle

import (
	"context"
	"fmt"
	"time" // Added for timeouts etc. // 添加用于超时等

	"strings"

	"github.com/turtacn/chasi-bod/common/errors"
	"github.com/turtacn/chasi-bod/common/types/enum"
	"github.com/turtacn/chasi-bod/common/utils" // Assuming logger is here // 假设日志记录器在这里
	// "github.com/turtacn/chasi-bod/pkg/builder"  // Assuming a builder orchestrator interface exists // 假设构建器协调器接口存在
	"github.com/turtacn/chasi-bod/pkg/config/model"
	"github.com/turtacn/chasi-bod/pkg/deployer" // Assuming deployer orchestrator interface exists // 假设部署器协调器接口存在
	"github.com/turtacn/chasi-bod/pkg/storage"  // Import storage package for backup/restore utilities // 导入 storage 包用于备份/恢复工具
	"github.com/turtacn/chasi-bod/pkg/vcluster" // Assuming vcluster manager interface exists // 假设 vcluster 管理器接口存在
	// You might need application deployer as well for app lifecycle management
	// 您可能也需要应用程序 deployer 用于应用程序生命周期管理
	// "github.com/turtacn/chasi-bod/pkg/application"
	// Import Kubernetes client-go if needed for status checks etc.
	// 如果需要用于状态检查等，则导入 Kubernetes client-go
	"k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/tools/clientcmd" // Needed to load kubeconfig // 需要它来加载 kubeconfig
)

// Manager defines the interface for managing the platform lifecycle.
// Manager 定义了管理平台生命周期的接口。
type Manager interface {
	// UpgradePlatform upgrades the chasi-bod platform to a new version defined by the new configuration.
	// UpgradePlatform 将 chasi-bod 平台升级到新配置定义的新版本。
	// This might involve building a new image and performing a rolling update or replacement.
	// 这可能涉及构建新镜像并执行滚动更新或替换。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// currentConfig: The currently active platform configuration. / 当前活动的平台配置。
	// newConfig: The configuration for the target platform version. / 目标平台版本的配置。
	// Returns an error if the upgrade fails.
	// 如果升级失败则返回错误。
	UpgradePlatform(ctx context.Context, currentConfig *model.PlatformConfig, newConfig *model.PlatformConfig) error

	// ScaleHostCluster scales the underlying Host Kubernetes cluster by adding or removing nodes.
	// ScaleHostCluster 通过添加或删除节点来扩缩容底层的 Host Kubernetes 集群。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// config: The current platform configuration. / 当前的平台配置。
	// nodesToAdd: List of node configurations to add. / 要添加的节点配置列表。
	// nodesToRemove: List of node configurations to remove. / 要删除的节点配置列表。
	// hostK8sClient: A Kubernetes client connected to the Host Cluster. / 连接到 Host 集群的 Kubernetes 客户端。
	// Returns an error if scaling fails.
	// 如果扩缩容失败则返回错误。
	ScaleHostCluster(ctx context.Context, config *model.PlatformConfig, nodesToAdd []model.NodeConfig, nodesToRemove []model.NodeConfig, hostK8sClient kubernetes.Interface) error

	// BackupPlatform backs up the platform state (e.g., etcd, configurations).
	// BackupPlatform 备份平台状态（例如，etcd、配置）。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// config: The current platform configuration. / 当前的平台配置。
	// hostK8sClient: A Kubernetes client connected to the Host Cluster (needed for ETCD backup).
	// hostK8sClient: 连接到 Host 集群的 Kubernetes 客户端（ETCD 备份需要）。
	// Returns an error if backup fails.
	// 如果备份失败则返回错误。
	BackupPlatform(ctx context.Context, config *model.PlatformConfig, hostK8sClient kubernetes.Interface) error

	// RestorePlatform restores the platform state from a backup.
	// RestorePlatform 从备份中恢复平台状态。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// config: The platform configuration corresponding to the backup. / 对应于备份的平台配置。
	// backupLocation: The location of the backup. / 备份的位置。
	// hostK8sClient: A Kubernetes client connected to the Host Cluster (needed for cleanup/verification).
	// hostK8sClient: 连接到 Host 集群的 Kubernetes 客户端（清理/验证需要）。
	// Returns an error if restoration fails.
	// 如果恢复失败则返回错误。
	RestorePlatform(ctx context.Context, config *model.PlatformConfig, backupLocation string, hostK8sClient kubernetes.Interface) error

	// // ManageApplicationLifecycle performs lifecycle operations on a specific application (e.g., upgrade, rollback).
	// // ManageApplicationLifecycle 对特定应用程序执行生命周期操作（例如，升级、回滚）。
	// ManageApplicationLifecycle(ctx context.Context, appName string, vclusterName string, operation LifecycleOperation, options interface{}, hostK8sClient kubernetes.Interface) error // Define LifecycleOperation enum and options // 定义 LifecycleOperation 枚举和选项
}

// defaultManager is a default implementation of the Lifecycle Manager.
// defaultManager 是 Lifecycle Manager 的默认实现。
type defaultManager struct {
	platformDeployer deployer.Deployer // The platform deployer orchestrator / 平台 deployer 协调器
	// imageBuilder     builder.Builder   // The platform image builder orchestrator / 平台镜像构建器协调器
	vclusterManager vcluster.Manager // The vcluster manager / vcluster 管理器
	// appDeployer application.Deployer // The application deployer / 应用程序 deployer
	// Add other dependencies like backup/restore tools
	// 添加其他依赖项，例如备份/恢复工具
}

// NewManager creates a new Lifecycle Manager.
// NewManager 创建一个新的 Lifecycle Manager。
// platformDeployer: The platform deployer orchestrator. / 平台 deployer 协调器。
// imageBuilder: The platform image builder orchestrator. / 平台镜像构建器协调器。
// vclusterManager: The vcluster manager. / vcluster 管理器。
// appDeployer: The application deployer (optional). / 应用程序 deployer（可选）。
// Returns a Lifecycle Manager implementation.
// 返回 Lifecycle Manager 实现。
func NewManager(platformDeployer deployer.Deployer, vclusterManager vcluster.Manager /*, appDeployer application.Deployer*/) Manager {
	return &defaultManager{
		platformDeployer: platformDeployer,
		// imageBuilder:     imageBuilder,
		vclusterManager: vclusterManager,
		// appDeployer: appDeployer,
	}
}

// UpgradePlatform upgrades the platform.
// UpgradePlatform 升级平台。
func (m *defaultManager) UpgradePlatform(ctx context.Context, currentConfig *model.PlatformConfig, newConfig *model.PlatformConfig) error {
	utils.GetLogger().Printf("Starting platform upgrade from config '%s' to '%s'", currentConfig.Metadata.Name, newConfig.Metadata.Name)

	// Step 1: Build the new platform image based on the new configuration
	// 步骤 1：根据新配置构建新的平台镜像
	utils.GetLogger().Println("Building new platform image...")
	// TODO: Use m.imageBuilder to build the new image
	// TODO: 使用 m.imageBuilder 构建新镜像
	// newImagePath, err := m.imageBuilder.Build(ctx, newConfig) // Assuming Build method exists on builder
	// if err != nil { return errors.NewWithCause(errors.ErrTypeBuilder, "failed to build new platform image", err) }
	// utils.GetLogger().Printf("New platform image built at: %s", newImagePath)
	utils.GetLogger().Println("Placeholder: New platform image built.")
	// newImagePathPlaceholder := "placeholder/new-image.iso" // Simulate path // 模拟路径

	// Step 2: Implement upgrade strategy (e.g., rolling update, replace nodes)
	// This is complex and depends on desired downtime, infrastructure type etc.
	// This involves using the platform deployer to update the nodes with the new image and configuration.
	// 步骤 2：实现升级策略（例如，滚动更新，替换节点）
	// 这很复杂，取决于期望的停机时间、基础设施类型等。
	// 这涉及使用平台 deployer 更新节点以使用新镜像和配置。
	utils.GetLogger().Println("Applying platform upgrade strategy...")

	// Example: Rolling replace strategy for worker nodes (simplified)
	// 示例：工作节点的滚动替换策略（简化）
	// Identify nodes to upgrade (e.g., worker nodes that have config changes or need image update)
	// 识别要升级的节点（例如，有配置更改或需要镜像更新的工作节点）
	// This requires comparing currentConfig and newConfig to find nodes that need updating.
	// 这需要比较 currentConfig 和 newConfig 以查找需要更新的节点。
	// For simplicity, let's assume we upgrade all nodes in the new config that were also in the old config.
	// 为简单起见，我们假设升级新配置中所有也在旧配置中的节点。

	// TODO: Get Host K8s client using currentConfig (needed to drain/uncordon nodes during rolling upgrade)
	// TODO: 使用 currentConfig 获取 Host K8s 客户端（滚动升级期间需要排空/取消封锁节点）
	// hostK8sClient, err := getHostK8sClient(currentConfig) // Need to implement this function
	// if err != nil { return errors.NewWithCause(errors.ErrTypeSystem, "failed to get host K8s client for upgrade", err) }
	// var hostK8sClient interface{} // Placeholder - should be kubernetes.Interface

	// Need logic to iterate through nodes, drain, reimage/reconfigure, uncordon/join
	// 需要逻辑来遍历节点，排空，重新镜像/重新配置，取消封锁/加入

	utils.GetLogger().Println("Placeholder: Implementing rolling upgrade for nodes.")
	// for _, newNodeCfg := range newConfig.Cluster.Nodes {
	// 	// Find the corresponding node in currentConfig
	// 	currentNodeCfg := findNodeConfig(currentConfig, newNodeCfg.Address) // Need helper
	//
	// 	// Decide if the node needs upgrade (e.g., if config changed or image needs replacing)
	// 	needsUpgrade := nodesNeedUpgrade(currentNodeCfg, &newNodeCfg) // Need helper to compare configs
	//
	// 	if needsUpgrade {
	// 		utils.GetLogger().Printf("Upgrading node: %s", newNodeCfg.Address)
	// 		// 1. Drain node using hostK8sClient
	// 		// 2. Reimage/reconfigure node using the new image/config (requires deployer phases)
	// 		// 3. Uncordon/Join node using hostK8sClient/deployer
	// 	} else {
	// 		utils.GetLogger().Printf("Node %s does not require upgrade.", newNodeCfg.Address)
	// 	}
	// }

	// Handling master node upgrades is particularly complex and version-specific.
	// 处理主节点升级尤其复杂且与版本相关。
	utils.GetLogger().Println("Placeholder: Master node upgrade logic needs careful implementation.")

	utils.GetLogger().Println("Platform upgrade completed successfully (placeholder).")
	return nil
}

// ScaleHostCluster scales the cluster.
// ScaleHostCluster 扩缩容集群。
func (m *defaultManager) ScaleHostCluster(ctx context.Context, config *model.PlatformConfig, nodesToAdd []model.NodeConfig, nodesToRemove []model.NodeConfig, hostK8sClient kubernetes.Interface) error {
	utils.GetLogger().Printf("Scaling Host Cluster: Adding %d nodes, Removing %d nodes", len(nodesToAdd), len(nodesToRemove))

	// Step 1: Add new nodes using the platform deployer's AddNode method
	// 步骤 1：使用平台 deployer 的 AddNode 方法添加新节点
	if len(nodesToAdd) > 0 {
		utils.GetLogger().Printf("Adding %d nodes: %s", len(nodesToAdd), formatNodeAddresses(nodesToAdd))
		for _, nodeCfg := range nodesToAdd {
			nodeCtx, cancel := context.WithTimeout(ctx, 20*time.Minute) // Context for adding a single node
			defer cancel()
			// The deployer's AddNode should handle the phases required for a new node.
			// deployer 的 AddNode 应该处理新节点所需的阶段。
			if err := m.platformDeployer.AddNode(nodeCtx, config, &nodeCfg); err != nil {
				return errors.NewWithCause(errors.ErrTypeSystem, fmt.Sprintf("failed to add node %s", nodeCfg.Address), err)
			}
		}
		utils.GetLogger().Println("New nodes added and joined to cluster successfully.")
	} else {
		utils.GetLogger().Println("No nodes to add.")
	}

	// Step 2: Remove nodes using the platform deployer's RemoveNode method
	// 步骤 2：使用平台 deployer 的 RemoveNode 方法移除节点
	if len(nodesToRemove) > 0 {
		utils.GetLogger().Printf("Removing %d nodes: %s", len(nodesToRemove), formatNodeAddresses(nodesToRemove))
		for _, nodeCfg := range nodesToRemove {
			nodeCtx, cancel := context.WithTimeout(ctx, 10*time.Minute) // Context for removing a single node
			defer cancel()
			// The deployer's RemoveNode should handle draining, deleting from K8s, and OS cleanup.
			// deployer 的 RemoveNode 应该处理排空、从 K8s 删除以及操作系统清理。
			if err := m.platformDeployer.RemoveNode(nodeCtx, &nodeCfg, hostK8sClient); err != nil { // Requires Host K8s client
				return errors.NewWithCause(errors.ErrTypeSystem, fmt.Sprintf("failed to remove node %s", nodeCfg.Address), err)
			}
		}
		utils.GetLogger().Println("Nodes drained and removed from cluster successfully.")
	} else {
		utils.GetLogger().Println("No nodes to remove.")
	}

	// TODO: Update platform configuration state to reflect the new node list after successful scaling
	// TODO: 扩缩容成功后，更新平台配置状态以反映新的节点列表

	utils.GetLogger().Println("Host Cluster scaling completed successfully.")
	return nil
}

// BackupPlatform backs up platform state.
// BackupPlatform 备份平台状态。
func (m *defaultManager) BackupPlatform(ctx context.Context, config *model.PlatformConfig, hostK8sClient kubernetes.Interface) error {
	utils.GetLogger().Printf("Starting platform backup for config '%s'...", config.Metadata.Name)

	// Step 1: Backup Host Cluster ETCD (requires connecting to master nodes and K8s client for endpoints)
	// 步骤 1：备份 Host Cluster ETCD（需要连接到主节点和 K8s 客户端获取端点）
	if config.DFXConfig.Reliability.ETCDBackup.Enabled {
		utils.GetLogger().Println("Backing up Host Cluster ETCD...")
		// Find master nodes in the configuration
		// 在配置中找到主节点
		masterNodes := []model.NodeConfig{}
		for _, node := range config.Cluster.Nodes {
			for _, role := range node.Roles {
				if role == enum.RoleMaster {
					masterNodes = append(masterNodes, node)
					break
				}
			}
		}
		if len(masterNodes) == 0 {
			utils.GetLogger().Println("No master nodes found in config, cannot perform ETCD backup.")
			// Decide if this is an error or warning
			// 决定这是错误还是警告
		} else {
			// TODO: Get ETCD endpoints from Host K8s cluster (e.g., by checking static pods or master node status)
			// TODO: 从 Host K8s 集群获取 ETCD 端点（例如，通过检查静态 pod 或主节点状态）
			etcdEndpoints := "placeholder-etcd-endpoints" // Replace with actual endpoints // 占位符 - 替换为实际端点

			// Perform backup on one of the master nodes
			// 在其中一个主节点上执行备份
			masterNodeToBackup := masterNodes[0]        // For simplicity, backup from the first master
			backupPathOnNode := "/tmp/etcd-snapshot.db" // Temporary path on the node // 节点上的临时路径

			// Need to determine the correct etcdctl command, including certs.
			// The certs are usually in /etc/kubernetes/pki/etcd/ on master nodes.
			// 需要确定正确的 etcdctl 命令，包括证书。
			// 证书通常位于主节点上的 /etc/kubernetes/pki/etcd/。
			etcdctlCmd := fmt.Sprintf("sudo etcdctl --endpoints=%s --cacert=/etc/kubernetes/pki/etcd/ca.crt --cert=/etc/kubernetes/pki/etcd/server.crt --key=/etc/kubernetes/pki/etcd/server.key", etcdEndpoints)
			// Assuming cert paths are standard after kubeadm install // 假设 kubeadm 安装后证书路径是标准的

			err := storage.BackupETCD(ctx, &masterNodeToBackup, etcdctlCmd, backupPathOnNode) // Use storage.BackupETCD
			if err != nil {
				return errors.NewWithCause(errors.ErrTypeReliability, "ETCD backup failed", err)
			}

			// Step 1a: Transfer the ETCD snapshot from the node to the backup location
			// 步骤 1a：将 ETCD 快照从节点传输到备份位置
			utils.GetLogger().Printf("Transferring ETCD snapshot from node %s to backup location %s...", masterNodeToBackup.Address, config.DFXConfig.Reliability.ETCDBackup.Location)
			// TODO: Implement transfer logic (SSH/SFTP to local, then upload to S3/NFS etc.)
			// TODO: 实现传输逻辑（SSH/SFTP 到本地，然后上传到 S3/NFS 等）
			utils.GetLogger().Println("Placeholder: ETCD snapshot transferred.")

		}
	} else {
		utils.GetLogger().Println("ETCD backup is disabled by configuration. Skipping.")
	}

	// Step 2: Backup chasi-bod configuration files
	// 步骤 2：备份 chasi-bod 配置文件
	if config.DFXConfig.Reliability.ConfigBackup.Enabled {
		utils.GetLogger().Printf("Backing up chasi-bod configuration files from '%s'...", "local-config-paths") // Specify config paths // 指定配置路径
		// Identify configuration files to backup (e.g., the main config file path)
		// 识别要备份的配置文件（例如，主配置文件路径）
		// configFilesToBackup := []string{ /* constants.DefaultConfigPath, */ "other/config/files" /* ... */} // TODO: Get actual config file paths // 获取实际配置文件路径
		localBackupTempDir := "/tmp/chasi-bod-config-backup" // Temporary local dir // 临时本地目录

		// This step might backup local config files or files from master nodes.
		// If configs are on master nodes, use storage.BackupConfigFiles.
		// 如果配置在主节点上，使用 storage.BackupConfigFiles。
		// If configs are local (like the main config.yaml), just copy them locally.
		// 如果配置是本地的（如主 config.yaml），只需在本地复制它们。
		utils.GetLogger().Printf("Placeholder: Backing up local configuration files.")
		// err := utils.CopyFile(constants.DefaultConfigPath, filepath.Join(localBackupTempDir, filepath.Base(constants.DefaultConfigPath))) // Example local copy

		// Step 2a: Transfer the configuration files to the backup location
		// 步骤 2a：将配置文件传输到备份位置
		utils.GetLogger().Printf("Transferring configuration files from %s to location %s...", localBackupTempDir, config.DFXConfig.Reliability.ConfigBackup.Location)
		// TODO: Implement transfer logic
		utils.GetLogger().Println("Placeholder: Configuration files transferred.")

		// TODO: Clean up local temporary backup directory
		// TODO: 清理本地临时备份目录
		// utils.RemovePath(localBackupTempDir)

	} else {
		utils.GetLogger().Println("Configuration backup is disabled by configuration. Skipping.")
	}

	// TODO: Add backup of other state if necessary (e.g., vcluster metadata not stored in K8s)
	// TODO: 如果需要，添加其他状态的备份（例如，未存储在 K8s 中的 vcluster 元数据）

	utils.GetLogger().Println("Platform backup completed successfully (placeholder).")
	return nil
}

// RestorePlatform restores platform state.
// RestorePlatform 从备份中恢复平台状态。
func (m *defaultManager) RestorePlatform(ctx context.Context, config *model.PlatformConfig, backupLocation string, hostK8sClient kubernetes.Interface) error {
	utils.GetLogger().Printf("Starting platform restoration from backup location '%s' for config '%s'...", backupLocation, config.Metadata.Name)

	// Restoration is a complex process, potentially requiring rebuilding the Host OS and then restoring ETCD.
	// This outline is highly simplified.
	// 恢复是一个复杂的过程，可能需要重建 Host OS，然后恢复 ETCD。
	// 这个大纲被高度简化了。

	// Step 0: Potentially prepare target nodes (e.g., reimage if necessary based on restore strategy)
	// 步骤 0：可能准备目标节点（例如，如果需要，根据恢复策略重新镜像）
	utils.GetLogger().Println("Placeholder: Preparing target nodes for restoration.")
	// This might involve using the builder/deployer to deploy a base OS image.
	// 这可能涉及使用 builder/deployer 部署基础操作系统镜像。

	// Step 1: Restore chasi-bod configuration files
	// 步骤 1：恢复 chasi-bod 配置文件
	if config.DFXConfig.Reliability.ConfigBackup.Enabled { // Check if config backup was enabled in the backup's config
		utils.GetLogger().Printf("Restoring chasi-bod configuration files from '%s'...", backupLocation)
		// TODO: Implement transfer logic to get config files from backupLocation to a temporary local dir
		// TODO: 实现传输逻辑以从 backupLocation 获取配置文件到临时本地目录
		// localRestoreTempDir := "/tmp/chasi-bod-config-restore" // Temporary local dir // 临时本地目录
		// TODO: Copy files from localRestoreTempDir to their intended destinations (e.g., constants.DefaultConfigPath)
		// If config files were from master nodes, use storage.RestoreConfigFiles.
		// 如果配置文件来自主节点，使用 storage.RestoreConfigFiles。
		// If config files are local, just copy them locally.
		// 如果配置文件是本地的，只需在本地复制它们。
		utils.GetLogger().Println("Placeholder: Chasi-bod configuration files restored.")
		// TODO: Clean up local temporary restore directory
		// TODO: 清理本地临时恢复目录
		// utils.RemovePath(localRestoreTempDir)
	} else {
		utils.GetLogger().Println("Configuration backup was not enabled in the backup's config. Skipping config file restoration.")
	}

	// Step 2: Restore Host Cluster ETCD (requires connecting to master nodes)
	// 步骤 2：恢复 Host Cluster ETCD（需要连接到主节点）
	if config.DFXConfig.Reliability.ETCDBackup.Enabled { // Check if ETCD backup was enabled in the backup's config
		utils.GetLogger().Println("Restoring Host Cluster ETCD...")
		// Find master nodes in the configuration
		// 在配置中找到主节点
		masterNodes := []model.NodeConfig{}
		for _, node := range config.Cluster.Nodes {
			for _, role := range node.Roles {
				if role == enum.RoleMaster {
					masterNodes = append(masterNodes, node)
					break
				}
			}
		}
		if len(masterNodes) == 0 {
			utils.GetLogger().Println("No master nodes found in config, cannot perform ETCD restoration.")
			// Decide if this is an error or warning
			// 决定这是错误还是警告
		} else {
			// Step 2a: Transfer the ETCD snapshot from the backup location to a master node
			// 步骤 2a：将 ETCD 快照从备份位置传输到主节点
			utils.GetLogger().Printf("Transferring ETCD snapshot from backup location '%s' to master node %s...", backupLocation, masterNodes[0].Address)
			// TODO: Implement transfer logic
			snapshotPathOnNode := "/tmp/etcd-snapshot-restore.db" // Temporary path on the node // 节点上的临时路径
			utils.GetLogger().Println("Placeholder: ETCD snapshot transferred to master node.")

			// Step 2b: Stop ETCD and API server, restore ETCD, start services
			// 步骤 2b：停止 ETCD 和 API 服务器，恢复 ETCD，启动服务
			masterNodeToRestore := masterNodes[0]
			etcdDataDirOnNode := "/var/lib/etcd" // Standard ETCD data dir // 标准 ETCD 数据目录
			etcdctlCmd := "sudo etcdctl ..."     // Needs correct endpoints and certs from the *restored* state
			// Needs careful implementation to handle stopping/starting services and potentially rebuilding the cluster state.
			// 需要小心实现以处理停止/启动服务，并可能重建集群状态。
			// This is a high-risk operation.
			// 这是一个高风险操作。
			err := storage.RestoreETCD(ctx, &masterNodeToRestore, etcdctlCmd, snapshotPathOnNode, etcdDataDirOnNode) // Use storage.RestoreETCD
			if err != nil {
				return errors.NewWithCause(errors.ErrTypeReliability, "ETCD restoration failed", err)
			}
			utils.GetLogger().Println("ETCD restored on master node.")
		}
	} else {
		utils.GetLogger().Println("ETCD backup was not enabled in the backup's config. Skipping ETCD restoration.")
	}

	// Step 3: Verify Host Cluster health after restoration
	// 步骤 3：验证 Host Cluster 健康
	utils.GetLogger().Println("Verifying Host Cluster health after restoration...")
	// Use the provided hostK8sClient to check component statuses, node readiness etc.
	// 使用提供的 hostK8sClient 检查组件状态、节点就绪等。
	// This might require re-initializing the hostK8sClient if the kubeconfig was restored.
	// 如果 kubeconfig 已恢复，这可能需要重新初始化 hostK8sClient。
	// TODO: Implement K8s health checks using hostK8sClient
	utils.GetLogger().Println("Placeholder: Host Cluster health verified.")

	// Step 4: Reconcile vclusters and applications (if needed)
	// 步骤 4：协调 vcluster 和应用程序（如果需要）
	// If ETCD restoration brought back the K8s state, vcluster and application deployments should theoretically reconcile themselves.
	// However, manual steps might be needed if there were external dependencies or changes.
	// 如果 ETCD 恢复带回了 K8s 状态，vcluster 和应用程序部署理论上应该能够自我协调。
	// 但是，如果存在外部依赖项或更改，可能需要手动步骤。
	utils.GetLogger().Println("Placeholder: Reconciling vclusters and applications.")

	utils.GetLogger().Println("Platform restoration completed successfully (placeholder).")
	return nil
}

// TODO: Implement helper functions for specific operations (drainNode, removeNodeFromCluster, reimageNode, addNode, etcd backup/restore)
// TODO: 实现特定操作的辅助函数（drainNode, removeNodeFromCluster, reimageNode, addNode, etcd backup/restore）
// TODO: Implement helper functions to get Host K8s client (potentially using the restored kubeconfig)
// TODO: 实现获取 Host K8s 客户端的辅助函数（可能使用恢复的 kubeconfig）
// func getHostK8sClient(config *model.PlatformConfig) (kubernetes.Interface, error) { ... }
// TODO: Implement helper functions to compare node configurations for upgrade needs
// TODO: 实现比较节点配置以判断是否需要升级的辅助函数
// func nodesNeedUpgrade(currentNode, newNode *model.NodeConfig) bool { ... }
// func findNodeConfig(config *model.PlatformConfig, address string) *model.NodeConfig { ... }
// TODO: Implement helper functions for transferring backup files (local, SSH/SFTP, S3, NFS)
// TODO: 实现传输备份文件的辅助函数（本地、SSH/SFTP、S3、NFS）

func formatNodeAddresses(nodes []model.NodeConfig) string {
	addresses := make([]string, len(nodes))
	for i, node := range nodes {
		addresses[i] = node.Address
	}
	return strings.Join(addresses, ", ")
}
