// Package cli implements the command-line interface for chasi-bod.
// 包 cli 实现了 chasi-bod 的命令行界面。
package cli

import (
	"context"
	"fmt"
	//"os"
	"strings" // Added for string joining // 添加用于字符串拼接
	//"time" // Added for timeouts // 添加用于超时

	"github.com/spf13/cobra"                        // Using Cobra for CLI structure / 使用 Cobra 构建 CLI 结构
	"github.com/turtacn/chasi-bod/common/constants" // Assuming constants are here // 假设常量在这里
	"github.com/turtacn/chasi-bod/common/errors"    // Assuming custom errors are here // 假设自定义错误在这里
	"github.com/turtacn/chasi-bod/common/utils"     // Assuming logger is here // 假设日志记录器在这里
	"github.com/turtacn/chasi-bod/pkg/application"  // Import application package // 导入应用程序包
	//"github.com/turtacn/chasi-bod/pkg/builder" // Assuming builder package exists // 假设 builder 包存在
	"github.com/turtacn/chasi-bod/pkg/config/loader"                   // Assuming config loader exists // 假设配置加载器存在
	"github.com/turtacn/chasi-bod/pkg/config/model"                    // Import config model // 导入配置模型
	"github.com/turtacn/chasi-bod/pkg/config/validator"                // Assuming config validator exists // 假设配置校验器存在
	"github.com/turtacn/chasi-bod/pkg/deployer"                        // Assuming deployer package exists // 假设 deployer 包存在
	"github.com/turtacn/chasi-bod/pkg/dfx/healthz"                     // Import healthz package // 导入 healthz 包
	"github.com/turtacn/chasi-bod/pkg/lifecycle"                       // Assuming lifecycle package exists // 假设生命周期包存在
	vcluster_mgr "github.com/turtacn/chasi-bod/pkg/vcluster"           // Alias for vcluster manager // vcluster manager 别名
	vcluster_client "github.com/turtacn/chasi-bod/pkg/vcluster/client" // Alias for vcluster client // vcluster 客户端别名
	// Placeholder for Kubernetes client-go, needed for some commands
	// Kubernetes client-go 的占位符，某些命令需要它
	"k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/tools/clientcmd" // Needed to load kubeconfig // 需要它来加载 kubeconfig
	// "k8s.io/client-go/rest" // Needed to build client config // 需要它来构建客户端配置
)

// RootCmd represents the base command when called without any subcommands.
// RootCmd 表示在不带任何子命令调用时的基本命令。
var RootCmd = &cobra.Command{
	Use:   constants.DefaultProjectName,
	Short: "Chasi-Bod: Manage fused business systems on shared Kubernetes with vCluster isolation.",
	Long: `Chasi-Bod is a platform for building, sharing, and running fused business systems
on a shared Kubernetes cluster using vCluster for isolation.

It simplifies the lifecycle management of both the underlying platform and the
complex applications deployed within it.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize logger here if not done in main, or add more config
		// If logging flags are added, parse them here and re-initialize logger if needed
		// 如果 main 中未初始化日志记录器，则在此处初始化，或添加更多配置
		// 如果添加了日志记录标志，在此处解析它们并根据需要重新初始化日志记录器

		// Ensure context is available for commands
		// 确保上下文对命令可用
		if cmd.Context() == nil {
			cmd.SetContext(context.Background())
		}
		return nil
	},
}

// configFilePath is the path to the chasi-bod platform configuration file.
// configFilePath 是 chasi-bod 平台配置文件的路径。
var configFilePath string

// init initializes the CLI commands and flags.
// init 初始化 CLI 命令和标志。
func init() {
	// Add global flags here
	// 在此处添加全局标志
	RootCmd.PersistentFlags().StringVarP(&configFilePath, "config", "c", constants.DefaultConfigPath, "Path to the chasi-bod platform configuration file")

	// Add subcommands
	// 添加子命令
	RootCmd.AddCommand(buildCmd)
	RootCmd.AddCommand(deployCmd)
	RootCmd.AddCommand(upgradeCmd)
	RootCmd.AddCommand(scaleCmd)
	RootCmd.AddCommand(backupCmd)
	RootCmd.AddCommand(restoreCmd)
	RootCmd.AddCommand(healthzCmd)     // Add healthz command // 添加 healthz 命令
	RootCmd.AddCommand(vclusterCmd)    // Base command for vcluster operations // vcluster 操作的基本命令
	RootCmd.AddCommand(applicationCmd) // Base command for application operations // 应用程序操作的基本命令
}

// Execute adds all child commands to the root command and sets flags appropriately.
// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once.
// 这由 main.main() 调用。它只需要发生一次。
func Execute() error {
	return RootCmd.Execute()
}

// --- Command Handlers ---

// buildCmd represents the build command.
// buildCmd 表示 build 命令。
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the chasi-bod platform image",
	Long:  `Builds the reproducible chasi-bod platform image based on the configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(cmd.Context(), constants.DefaultTimeout*2) // Example timeout for build // 示例构建超时时间
		defer cancel()

		// Load configuration
		// 加载配置
		config, err := loader.LoadConfig(configFilePath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Validate configuration
		// 校验配置
		if err := validator.ValidateConfig(config); err != nil {
			return fmt.Errorf("config validation failed: %w", err)
		}

		// Create a new builder orchestrator
		// 创建一个新的 builder 协调器
		bldr, err := builder.NewBuilder() // Assuming NewBuilder exists // 假设 NewBuilder 存在
		if err != nil {
			return fmt.Errorf("failed to create builder: %w", err)
		}

		// Run the build process
		// 运行构建过程
		utils.GetLogger().Println("Starting platform image build...")
		outputPath, err := bldr.Build(ctx, config) // Assuming Build method exists // 假设 Build 方法存在
		if err != nil {
			return fmt.Errorf("platform image build failed: %w", err)
		}

		utils.GetLogger().Printf("Platform image built successfully at: %s", outputPath)
		return nil
	},
}

// deployCmd represents the deploy command.
// deployCmd 表示 deploy 命令。
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the chasi-bod platform to target nodes",
	Long:  `Deploys the built chasi-bod platform image to the specified target nodes and initializes the Host Kubernetes cluster and vclusters.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(cmd.Context(), constants.DefaultTimeout*5) // Example timeout for deploy // 示例部署超时时间
		defer cancel()

		// Load configuration
		// 加载配置
		config, err := loader.LoadConfig(configFilePath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Validate configuration
		// 校验配置
		if err := validator.ValidateConfig(config); err != nil {
			return fmt.Errorf("config validation failed: %w", err)
		}

		// Create a new deployer orchestrator
		// 创建一个新的 deployer 协调器
		dplr, err := deployer.NewDeployer() // Assuming NewDeployer exists // 假设 NewDeployer 存在
		if err != nil {
			return fmt.Errorf("failed to create deployer: %w", err)
		}

		// Run the deployment process
		// 运行部署过程
		utils.GetLogger().Println("Starting platform deployment...")
		if err := dplr.Deploy(ctx, config); err != nil {
			return fmt.Errorf("platform deployment failed: %w", err)
		}

		utils.GetLogger().Println("Platform deployment completed successfully.")
		return nil
	},
}

// upgradeCmd represents the upgrade command.
// upgradeCmd 表示 upgrade 命令。
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the chasi-bod platform",
	Long:  `Upgrades the running chasi-bod platform to a new version specified by the configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(cmd.Context(), constants.DefaultTimeout*10) // Example timeout for upgrade // 示例升级超时时间
		defer cancel()

		// Load current and new configurations
		// 加载当前和新配置
		// This might require a flag for the current config path if different from the new one
		// 这可能需要一个标志来指定当前配置文件的路径，如果与新配置不同
		// For simplicity, assume the current config is stored in a default location derived from the name.
		// 为简单起见，假设当前配置存储在从名称派生的默认位置。
		currentConfigPath := strings.Replace(configFilePath, ".yaml", ".current.yaml", 1) // Example derivation // 示例派生
		currentConfig, err := loader.LoadConfig(currentConfigPath)
		if err != nil {
			// If current config not found, it might be the first deployment or an error. Decide behavior.
			// 如果未找到当前配置，可能是首次部署或错误。决定行为。
			utils.GetLogger().Printf("Warning: Current config not found at %s. Assuming first deployment or error in getting current state.", currentConfigPath)
			// For upgrade, not finding the current config is likely an error.
			// 对于升级，未找到当前配置可能是错误。
			return fmt.Errorf("failed to load current config from %s: %w", currentConfigPath, err)
		}

		newConfig, err := loader.LoadConfig(configFilePath)
		if err != nil {
			return fmt.Errorf("failed to load new config from %s: %w", configFilePath, err)
		}

		// Validate new configuration
		// 校验新配置
		if err := validator.ValidateConfig(newConfig); err != nil {
			return fmt.Errorf("new config validation failed: %w", err)
		}

		// Create necessary managers
		// 创建必要的管理器
		dplr, err := deployer.NewDeployer()
		if err != nil {
			return fmt.Errorf("failed to create deployer: %w", err)
		}
		bldr, err := builder.NewBuilder() // Assuming Builder orchestrator exists
		if err != nil {
			return fmt.Errorf("failed to create builder: %w", err)
		}
		// Need Host K8s client for vcluster manager and lifecycle manager
		// 需要 Host K8s 客户端用于 vcluster 管理器和生命周期管理器
		// Requires loading kubeconfig, which should be available after initial deploy/previous upgrade.
		// 需要加载 kubeconfig，它应该在初始部署/之前的升级后可用。
		hostK8sClient, err := getHostK8sClient() // Needs helper function to load kubeconfig
		if err != nil {
			return fmt.Errorf("failed to get host K8s client: %w", err)
		}
		vclusterMgr := vcluster_mgr.NewManager(hostK8sClient) // Assuming NewManager takes K8s client

		// Create a new lifecycle manager
		// 创建一个新的生命周期管理器
		lifecycleMgr := lifecycle.NewManager(dplr, bldr, vclusterMgr) // Pass all dependencies

		// Run the upgrade process
		// 运行升级过程
		utils.GetLogger().Println("Starting platform upgrade...")
		if err := lifecycleMgr.UpgradePlatform(ctx, currentConfig, newConfig); err != nil {
			return fmt.Errorf("platform upgrade failed: %w", err)
		}

		utils.GetLogger().Println("Platform upgrade completed successfully.")
		// TODO: Optionally, save the newConfig as the current config after successful upgrade
		// TODO: 可选地，在升级成功后将新配置保存为当前配置
		// err = loader.SaveConfig(newConfig, currentConfigPath)
		// if err != nil { utils.GetLogger().Printf("Warning: Failed to save new config as current config: %v", err) }

		return nil
	},
}

// scaleCmd represents the scale command.
// scaleCmd 表示 scale 命令。
var scaleCmd = &cobra.Command{
	Use:   "scale",
	Short: "Scale the chasi-bod Host Cluster",
	Long:  `Scales the underlying Host Kubernetes cluster by adding or removing nodes based on the configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(cmd.Context(), constants.DefaultTimeout*5) // Example timeout for scale // 示例扩缩容超时时间
		defer cancel()

		// Load configuration (represents the desired state after scaling)
		// 加载配置（表示扩缩容后的期望状态）
		newConfig, err := loader.LoadConfig(configFilePath)
		if err != nil {
			return fmt.Errorf("failed to load new config from %s: %w", configFilePath, err)
		}

		// Validate configuration
		// 校验配置
		if err := validator.ValidateConfig(newConfig); err != nil {
			return fmt.Errorf("new config validation failed: %w", err)
		}

		// TODO: Load current cluster state to determine nodesToAdd/nodesToRemove
		// This is complex, might involve querying the cluster directly or reading state file.
		// For simplicity, let's assume the user provides nodesToAdd/nodesToRemove via flags or a separate file.
		// TODO: 加载当前集群状态以确定要添加/移除的节点
		// 这很复杂，可能涉及直接查询集群或读取状态文件。
		// 为简单起见，我们假设用户通过标志或单独的文件提供要添加/移除的节点。
		var nodesToAdd []model.NodeConfig    // Placeholder - populate from flags or config comparison // 占位符 - 从标志或配置比较中填充
		var nodesToRemove []model.NodeConfig // Placeholder - populate from flags or config comparison // 占位符 - 从标志或配置比较中填充

		utils.GetLogger().Println("Placeholder: Determining nodes to add/remove (needs current state logic or input).")

		// Create necessary managers (deployer, builder, vcluster manager)
		// 创建必要的管理器（deployer, builder, vcluster manager）
		dplr, err := deployer.NewDeployer()
		if err != nil {
			return fmt.Errorf("failed to create deployer: %w", err)
		}
		bldr, err := builder.NewBuilder() // Assuming Builder orchestrator exists
		if err != nil {
			return fmt.Errorf("failed to create builder: %w", err)
		}
		// Need Host K8s client
		// 需要 Host K8s 客户端
		hostK8sClient, err := getHostK8sClient() // Needs helper function to load kubeconfig
		if err != nil {
			return fmt.Errorf("failed to get host K8s client: %w", err)
		}
		vclusterMgr := vcluster_mgr.NewManager(hostK8sClient)

		// Create a new lifecycle manager
		// 创建一个新的生命周期管理器
		lifecycleMgr := lifecycle.NewManager(dplr, bldr, vclusterMgr)

		// Run the scaling process
		// 运行扩缩容过程
		utils.GetLogger().Println("Starting Host Cluster scaling...")
		if err := lifecycleMgr.ScaleHostCluster(ctx, newConfig, nodesToAdd, nodesToRemove, hostK8sClient); err != nil { // Pass hostK8sClient
			return fmt.Errorf("host cluster scaling failed: %w", err)
		}

		utils.GetLogger().Println("Host Cluster scaling completed successfully.")
		// TODO: Optionally, save the newConfig as the current config after successful scaling
		// TODO: 可选地，在扩缩容成功后将新配置保存为当前配置

		return nil
	},
}

// backupCmd represents the backup command.
// backupCmd 表示 backup 命令。
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup the chasi-bod platform state",
	Long:  `Backs up the state of the chasi-bod platform, including ETCD and configuration files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(cmd.Context(), constants.DefaultTimeout*1) // Example timeout for backup // 示例备份超时时间
		defer cancel()

		// Load configuration
		// 加载配置
		config, err := loader.LoadConfig(configFilePath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Validate configuration (especially DFX/Reliability section)
		// 校验配置（尤其是 DFX/Reliability 部分）
		if err := validator.ValidateConfig(config); err != nil {
			return fmt.Errorf("config validation failed: %w", err)
		}

		// Create necessary managers (lifecycle manager needs backup capability)
		// 创建必要的管理器（生命周期管理器需要备份能力）
		// Need Host K8s client for ETCD backup
		// 需要 Host K8s 客户端用于 ETCD 备份
		hostK8sClient, err := getHostK8sClient() // Needs helper function to load kubeconfig
		if err != nil {
			// If we can't get K8s client, we might still be able to backup configs. Decide behavior.
			// 如果无法获取 K8s 客户端，可能仍然可以备份配置。决定行为。
			utils.GetLogger().Printf("Warning: Failed to get host K8s client for ETCD backup: %v", err)
			// Proceed with config backup only if K8s client is not essential for it.
			// 如果 K8s 客户端对于配置备份不是必需的，则只进行配置备份。
			// For now, return error if K8s client is needed for any part of backup.
			// 现在，如果备份的任何部分需要 K8s 客户端，则返回错误。
			return fmt.Errorf("failed to get host K8s client for backup: %w", err)
		}

		lifecycleMgr := lifecycle.NewManager(nil, nil, nil) // Deployer, Builder, VClusterManager might not be needed for simple backup/restore // 简单备份/恢复可能不需要 Deployer, Builder, VClusterManager

		// Run the backup process
		// 运行备份过程
		utils.GetLogger().Println("Starting platform backup...")
		if err := lifecycleMgr.BackupPlatform(ctx, config, hostK8sClient); err != nil { // Pass hostK8sClient
			return fmt.Errorf("platform backup failed: %w", err)
		}

		utils.GetLogger().Println("Platform backup completed successfully.")
		return nil
	},
}

// restoreCmd represents the restore command.
// restoreCmd 表示 restore 命令。
var restoreCmd = &cobra.Command{
	Use:   "restore <backup-location>",
	Short: "Restore the chasi-bod platform state",
	Long:  `Restores the state of the chasi-bod platform from a backup. Requires careful handling.`,
	Args:  cobra.ExactArgs(1), // Requires backup location argument // 需要备份位置参数
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(cmd.Context(), constants.DefaultTimeout*2) // Example timeout for restore // 示例恢复超时时间
		defer cancel()

		backupLocation := args[0] // Get backup location from arguments // 从参数获取备份位置

		// Load configuration
		// The config file defines the *target* state after restoration.
		// 加载配置
		// 配置文件定义了恢复后的目标状态。
		config, err := loader.LoadConfig(configFilePath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Validate configuration
		// 校验配置
		if err := validator.ValidateConfig(config); err != nil {
			return fmt.Errorf("config validation failed: %w", err)
		}

		// Create necessary managers (lifecycle manager needs restore capability)
		// 创建必要的管理器（生命周期管理器需要恢复能力）
		// Need Host K8s client for verification after restore
		// 恢复后验证需要 Host K8s 客户端
		// Note: Getting the K8s client *before* ETCD restore might be tricky if ETCD is the source of truth.
		// Maybe load the kubeconfig from the backup? Or try to connect to the potentially down API server?
		// 注意：在 ETCD 恢复之前获取 K8s 客户端可能很棘手，如果 ETCD 是事实来源。
		// 也许从备份加载 kubeconfig？或者尝试连接可能已关闭的 API 服务器？
		// For now, assume we can get a client that might be unhealthy initially.
		// 现在，假设我们可以获取一个最初可能不健康的客户端。
		hostK8sClient, err := getHostK8sClient() // Needs helper function to load kubeconfig
		if err != nil {
			// Restoration might not require a healthy K8s client initially, but verification does.
			// 恢复最初可能不需要健康的 K8s 客户端，但验证需要。
			utils.GetLogger().Printf("Warning: Failed to get host K8s client before restoration: %v", err)
			// Continue, but later verification step might fail.
			// 继续，但后面的验证步骤可能会失败。
		}

		lifecycleMgr := lifecycle.NewManager(nil, nil, nil) // Deployer, Builder, VClusterManager might not be needed // 可能不需要 Deployer, Builder, VClusterManager

		// Run the restoration process
		// 运行恢复过程
		utils.GetLogger().Println("Starting platform restoration...")
		if err := lifecycleMgr.RestorePlatform(ctx, config, backupLocation, hostK8sClient); err != nil { // Pass hostK8sClient
			return fmt.Errorf("platform restoration failed: %w", err)
		}

		utils.GetLogger().Println("Platform restoration completed successfully.")
		return nil
	},
}

// healthzCmd represents the healthz command.
// healthzCmd 表示 healthz 命令。
var healthzCmd = &cobra.Command{
	Use:   "healthz",
	Short: "Check platform health",
	Long:  `Checks the health status of chasi-bod platform components and deployed infrastructure.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(cmd.Context(), constants.DefaultTimeout) // Example timeout for health check // 示例健康检查超时时间
		defer cancel()

		// Note: The health check server is typically started as part of a long-running chasi-bod process (e.g., a controller).
		// This command is for on-demand checks.
		// 注意：健康检查服务器通常作为长期运行的 chasi-bod 进程（例如，控制器）的一部分启动。
		// 此命令用于按需检查。

		// Create a healthz manager
		// 创建一个 healthz 管理器
		healthMgr := healthz.NewManager() // Assuming NewManager exists // 假设 NewManager 存在

		// Run cluster checks (requires Host K8s client and config)
		// 运行集群检查（需要 Host K8s 客户端和配置）
		// Load configuration (needed to identify vclusters etc. for checks)
		// 加载配置（需要用于识别 vcluster 等以进行检查）
		config, err := loader.LoadConfig(configFilePath)
		if err != nil {
			utils.GetLogger().Printf("Warning: Failed to load config for cluster health checks: %v. Some checks may be skipped.", err)
			// Decide if config loading failure should be a hard error for healthz command.
			// 决定配置加载失败是否应成为 healthz 命令的硬错误。
			// For now, continue but log warning.
			// 现在，继续但记录警告。
		}

		hostK8sClient, err := getHostK8sClient() // Needs helper function to load kubeconfig
		if err != nil {
			utils.GetLogger().Printf("Warning: Failed to get host K8s client for cluster health checks: %v. K8s-dependent checks will fail.", err)
			// Continue, but K8s checks will likely fail.
			// 继续，但 K8s 检查可能会失败。
		}

		// Run the cluster checks
		// 运行集群检查
		utils.GetLogger().Println("Running Host Cluster and vcluster health checks...")
		err = healthMgr.RunClusterChecks(ctx, hostK8sClient, config) // Pass client and config // 传递客户端和配置
		if err != nil {
			// RunClusterChecks returns an error if critical checks fail.
			// RunClusterChecks 在关键检查失败时返回错误。
			utils.GetLogger().Printf("Cluster health checks reported failures: %v", err)
			// Note: This command currently only runs cluster checks, not component checks.
			// If you want to run component checks (like logger status), call healthMgr.GetOverallStatus(ctx)
			// 注意：此命令目前仅运行集群检查，不运行组件检查。
			// 如果您想运行组件检查（例如，日志记录器状态），请调用 healthMgr.GetOverallStatus(ctx)
			return fmt.Errorf("cluster health checks failed: %w", err) // Return error to indicate failure // 返回错误表示失败
		}

		utils.GetLogger().Println("Cluster health checks completed successfully.")
		return nil // Return nil if cluster checks pass // 如果集群检查通过则返回 nil
	},
}

// vclusterCmd represents the base command for vcluster operations.
// vclusterCmd 表示 vcluster 操作的基本命令。
var vclusterCmd = &cobra.Command{
	Use:   "vcluster",
	Short: "Manage virtual Kubernetes clusters",
	Long:  `Provides subcommands to manage virtual Kubernetes clusters (vclusters).`,
}

// Add vcluster subcommands (e.g., create, delete, list, connect)
// 添加 vcluster 子命令（例如，create, delete, list, connect）
func init() {
	vclusterCmd.AddCommand(vclusterCreateCmd)
	vclusterCmd.AddCommand(vclusterDeleteCmd)
	vclusterCmd.AddCommand(vclusterListCmd)
	vclusterCmd.AddCommand(vclusterConnectCmd) // Example: generate kubeconfig or port-forward // 示例：生成 kubeconfig 或端口转发
}

var vclusterCreateCmd = &cobra.Command{
	Use:   "create <vcluster-name>",
	Short: "Create a new vcluster",
	Args:  cobra.ExactArgs(1), // Requires vcluster name argument // 需要 vcluster 名称参数
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(cmd.Context(), constants.DefaultTimeout*3) // Example timeout for create // 示例创建超时时间
		defer cancel()

		vclusterName := args[0]

		// Load main configuration
		// 加载主配置
		config, err := loader.LoadConfig(configFilePath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Find the specific vcluster configuration by name
		// 按名称查找特定的 vcluster 配置
		vclusterCfg, exists := config.VClusters[vclusterName]
		if !exists {
			return errors.New(errors.ErrTypeNotFound, fmt.Sprintf("vcluster configuration '%s' not found in config file", vclusterName))
		}

		// Get Host K8s client
		// 获取 Host K8s 客户端
		hostK8sClient, err := getHostK8sClient() // Needs helper function to load kubeconfig
		if err != nil {
			return fmt.Errorf("failed to get host K8s client: %w", err)
		}

		// Create vcluster manager
		// 创建 vcluster 管理器
		vclusterMgr := vcluster_mgr.NewManager(hostK8sClient) // Assuming NewManager takes K8s client

		// Use vcluster manager to create the vcluster
		// 使用 vcluster 管理器创建 vcluster
		utils.GetLogger().Printf("Starting creation of vcluster '%s'", vclusterName)
		if err := vclusterMgr.Create(ctx, &vclusterCfg); err != nil {
			return fmt.Errorf("failed to create vcluster '%s': %w", vclusterName, err)
		}
		utils.GetLogger().Printf("Vcluster '%s' creation initiated.", vclusterName)

		// Optionally, wait for vcluster to be ready after creation
		// 可选地，创建后等待 vcluster 就绪
		waitFlag, _ := cmd.Flags().GetBool("wait") // Assuming a --wait flag exists // 假设 --wait 标志存在
		if waitFlag {
			utils.GetLogger().Printf("Waiting for vcluster '%s' to be ready...", vclusterName)
			if err := vclusterMgr.WaitForReady(ctx, vclusterName); err != nil {
				return fmt.Errorf("vcluster '%s' did not become ready: %w", vclusterName, err)
			}
			utils.GetLogger().Printf("Vcluster '%s' is ready.", vclusterName)
		}

		utils.GetLogger().Printf("Vcluster '%s' command completed.", vclusterName)
		return nil
	},
}

var vclusterDeleteCmd = &cobra.Command{
	Use:   "delete <vcluster-name>",
	Short: "Delete a vcluster",
	Args:  cobra.ExactArgs(1), // Requires vcluster name argument // 需要 vcluster 名称参数
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(cmd.Context(), constants.DefaultTimeout*2) // Example timeout for delete // 示例删除超时时间
		defer cancel()

		vclusterName := args[0]

		// Get Host K8s client
		// 获取 Host K8s 客户端
		hostK8sClient, err := getHostK8sClient() // Needs helper function to load kubeconfig
		if err != nil {
			return fmt.Errorf("failed to get host K8s client: %w", err)
		}

		// Create vcluster manager
		// 创建 vcluster 管理器
		vclusterMgr := vcluster_mgr.NewManager(hostK8sClient)

		// Use vcluster manager to delete the vcluster
		// 使用 vcluster 管理器删除 vcluster
		utils.GetLogger().Printf("Starting deletion of vcluster '%s'", vclusterName)
		if err := vclusterMgr.Delete(ctx, vclusterName); err != nil {
			return fmt.Errorf("failed to delete vcluster '%s': %w", vclusterName, err)
		}
		utils.GetLogger().Printf("Vcluster '%s' deletion initiated.", vclusterName)

		// Optionally, wait for vcluster to be deleted
		// 可选地，等待 vcluster 删除
		waitFlag, _ := cmd.Flags().GetBool("wait") // Assuming a --wait flag exists // 假设 --wait 标志存在
		if waitFlag {
			utils.GetLogger().Printf("Waiting for vcluster '%s' to be deleted...", vclusterName)
			// Waiting for deletion is typically waiting for the host namespace to be gone
			// 等待删除通常是等待 host 命名空间消失
			// Use a helper from vcluster_mgr or implement here
			// 使用 vcluster_mgr 中的辅助函数或在此处实现
			utils.GetLogger().Printf("Placeholder: Waiting for vcluster '%s' namespace to be deleted.", vclusterName)
			// err = vclusterMgr.WaitForDeletion(ctx, vclusterName) // Needs implementation
			// if err != nil { return fmt.Errorf("vcluster '%s' did not fully delete: %w", vclusterName, err) }
			// utils.GetLogger().Printf("Vcluster '%s' is deleted.", vclusterName)
		}

		utils.GetLogger().Printf("Vcluster '%s' command completed.", vclusterName)
		return nil
	},
}

var vclusterListCmd = &cobra.Command{
	Use:   "list",
	Short: "List vclusters",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(cmd.Context(), constants.DefaultTimeout) // Example timeout for list // 示例列表超时时间
		defer cancel()

		// Get Host K8s client
		// 获取 Host K8s 客户端
		hostK8sClient, err := getHostK8sClient() // Needs helper function to load kubeconfig
		if err != nil {
			return fmt.Errorf("failed to get host K8s client: %w", err)
		}

		// Create vcluster manager
		// 创建 vcluster 管理器
		vclusterMgr := vcluster_mgr.NewManager(hostK8sClient)

		// Use vcluster manager to list vclusters
		// 使用 vcluster 管理器列出 vcluster
		// The List method might return a list of vcluster configs or statuses.
		// List 方法可能返回 vcluster 配置列表或状态列表。
		// For now, assume it returns something printable.
		// 现在，假设它返回可打印的内容。
		// vclusters, err := vclusterMgr.List(ctx) // Needs implementation
		// if err != nil { return fmt.Errorf("failed to list vclusters: %w", err) }

		// TODO: Format and print list of vclusters
		// TODO: 格式化并打印 vcluster 列表
		utils.GetLogger().Println("Placeholder: Listing vclusters (logic not implemented).")
		utils.GetLogger().Println("NAME\tHOST NAMESPACE\tK8S VERSION\tSTATUS") // Example header // 示例标题
		// for _, v := range vclusters { // Iterate and print details // 迭代并打印详细信息
		// 	fmt.Printf("%s\t%s\t%s\t%s\n", v.Name, v.Namespace, v.KubernetesVersion, "Unknown") // Example printing // 示例打印
		// }

		return errors.New(errors.ErrTypeNotImplemented, "vcluster list command not fully implemented")
	},
}

var vclusterConnectCmd = &cobra.Command{
	Use:   "connect <vcluster-name>",
	Short: "Connect to a vcluster (e.g., generate kubeconfig)",
	Args:  cobra.ExactArgs(1), // Requires vcluster name argument // 需要 vcluster 名称参数
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(cmd.Context(), constants.DefaultTimeout) // Example timeout for connect // 示例连接超时时间
		defer cancel()

		vclusterName := args[0]

		// Get Host K8s client
		// 获取 Host K8s 客户端
		hostK8sClient, err := getHostK8sClient() // Needs helper function to load kubeconfig
		if err != nil {
			return fmt.Errorf("failed to get host K8s client: %w", err)
		}

		// Use vcluster.client package to get vcluster kubeconfig or set up port-forwarding
		// 使用 vcluster.client 包获取 vcluster kubeconfig 或设置端口转发
		// For simplicity, let's implement generating kubeconfig here.
		// 为简单起见，我们在这里实现生成 kubeconfig。
		utils.GetLogger().Printf("Generating kubeconfig for vcluster '%s'...", vclusterName)
		kubeconfigContent, err := vcluster_client.GetVClusterKubeConfig(ctx, vclusterName, hostK8sClient) // Assuming this function exists // 假设此函数存在
		if err != nil {
			return fmt.Errorf("failed to get kubeconfig for vcluster '%s': %w", vclusterName, err)
		}

		// Output the kubeconfig content to stdout
		// 将 kubeconfig 内容输出到 stdout
		fmt.Println(string(kubeconfigContent))

		utils.GetLogger().Printf("Kubeconfig for vcluster '%s' generated and printed to stdout.", vclusterName)
		return nil
	},
}

// applicationCmd represents the base command for application operations.
// applicationCmd 表示应用程序操作的基本命令。
var applicationCmd = &cobra.Command{
	Use:   "app",
	Short: "Manage business applications",
	Long:  `Provides subcommands to deploy and manage business applications within vclusters.`,
}

// Add application subcommands (e.g., deploy, upgrade, delete, status)
// 添加应用程序子命令（例如，deploy, upgrade, delete, status）
func init() {
	applicationCmd.AddCommand(applicationDeployCmd)
	applicationCmd.AddCommand(applicationUpgradeCmd)
	applicationCmd.AddCommand(applicationDeleteCmd)
	applicationCmd.AddCommand(applicationStatusCmd)

	// Add flags that are common to multiple app commands
	// 为多个 app 命令添加通用标志
	applicationDeployCmd.Flags().Bool("wait", true, "Wait for the application rollout to complete")
	applicationUpgradeCmd.Flags().Bool("wait", true, "Wait for the application rollout to complete")
	applicationDeleteCmd.Flags().Bool("wait", true, "Wait for the application resources to be deleted")

	applicationDeleteCmd.Flags().String("vcluster", "", "Target vcluster name")
	applicationDeleteCmd.MarkFlagRequired("vcluster")
	applicationDeleteCmd.Flags().StringP("namespace", "n", "default", "Target namespace within the vcluster") // Add namespace flag // 添加命名空间标志

	applicationStatusCmd.Flags().String("vcluster", "", "Target vcluster name")
	applicationStatusCmd.MarkFlagRequired("vcluster")
	applicationStatusCmd.Flags().StringP("namespace", "n", "default", "Target namespace within the vcluster") // Add namespace flag // 添加命名空间标志
}

var applicationDeployCmd = &cobra.Command{
	Use:   "deploy <app-config-file>",
	Short: "Deploy a business application",
	Long:  `Deploys a business application to a target vcluster based on the application configuration file.`,
	Args:  cobra.ExactArgs(1), // Requires application config file path // 需要应用程序配置文件路径
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(cmd.Context(), constants.DefaultTimeout*3) // Example timeout for app deploy // 示例应用程序部署超时时间
		defer cancel()

		appConfigFile := args[0] // Get application config file path // 获取应用程序配置文件路径

		// Load and validate application configuration (standalone validation)
		// 加载和校验应用程序配置（独立校验）
		appConfig, err := application.LoadApplicationConfig(appConfigFile) // Assuming LoadApplicationConfig exists // 假设 LoadApplicationConfig 存在
		if err != nil {
			return fmt.Errorf("failed to load application config from %s: %w", appConfigFile, err)
		}

		// Load main platform configuration for full validation and Host K8s client
		// 加载主平台配置以进行完整校验和获取 Host K8s 客户端
		platformConfig, err := loader.LoadConfig(configFilePath)
		if err != nil {
			return fmt.Errorf("failed to load main platform config from %s for app deploy: %w", configFilePath, err)
		}

		// Perform full validation (checks against platform config)
		// 执行完整校验（对照平台配置检查）
		// Need a dedicated validation function for application config within platform config context
		// 需要一个在平台配置上下文内校验应用程序配置的专用函数
		// validator.ValidateApplicationConfigWithPlatform(appConfig, platformConfig) // Needs implementation // 需要实现
		utils.GetLogger().Println("Placeholder: Performing full application config validation against platform config.")
		// Example validation: check if target vcluster exists in platform config
		// 示例校验：检查目标 vcluster 是否存在于平台配置中
		if _, exists := platformConfig.VClusters[appConfig.VClusterName]; !exists {
			return errors.New(errors.ErrTypeValidation, fmt.Sprintf("target vcluster '%s' for application '%s' not found in main platform config", appConfig.VClusterName, appConfig.Name))
		}

		// Get Host K8s client
		// 获取 Host K8s 客户端
		hostK8sClient, err := getHostK8sClient() // Needs helper function to load kubeconfig
		if err != nil {
			return fmt.Errorf("failed to get host K8s client for app deploy: %w", err)
		}

		// Create application deployer
		// 创建应用程序 deployer
		appDeployer := application.NewDeployer() // Assuming NewDeployer exists // 假设 NewDeployer 存在

		// Run the deployment process
		// 运行部署过程
		utils.GetLogger().Printf("Starting application deployment for '%s' to vcluster '%s'", appConfig.Name, appConfig.VClusterName)
		// Pass the Host K8s client which is needed by the app deployer to get the vcluster client
		// 传递 Host K8s 客户端，app deployer 需要它来获取 vcluster 客户端
		if err := appDeployer.Deploy(ctx, appConfig, hostK8sClient); err != nil {
			return fmt.Errorf("application deployment failed for '%s': %w", appConfig.Name, err)
		}

		utils.GetLogger().Printf("Application '%s' deployed successfully.", appConfig.Name)

		// TODO: Optional - Wait for application rollout if --wait flag is set
		// TODO: 可选 - 如果设置了 --wait 标志，则等待应用程序 rollout
		// waitFlag, _ := cmd.Flags().GetBool("wait")
		// if waitFlag {
		// 	utils.GetLogger().Printf("Waiting for application '%s' rollout to complete...", appConfig.Name)
		// 	// This requires getting the vcluster client and checking resource statuses (Deployments, StatefulSets)
		// 	// vclusterClient, err := vcluster_client.GetVClusterClient(ctx, appConfig.VClusterName, hostK8sClient)
		// 	// err = application.WaitForRollout(ctx, vclusterClient, appConfig) // Needs helper
		// 	// if err != nil { return fmt.Errorf("application '%s' rollout failed: %w", appConfig.Name, err) }
		// 	utils.GetLogger().Printf("Placeholder: Waiting for rollout for '%s'.", appConfig.Name)
		// }

		return nil
	},
}

var applicationUpgradeCmd = &cobra.Command{
	Use:   "upgrade <app-config-file>",
	Short: "Upgrade a business application",
	Long:  `Upgrades a business application deployed in a vcluster based on the application configuration file.`,
	Args:  cobra.ExactArgs(1), // Requires application config file path // 需要应用程序配置文件路径
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(cmd.Context(), constants.DefaultTimeout*3) // Example timeout for app upgrade // 示例应用程序升级超时时间
		defer cancel()

		appConfigFile := args[0]

		// Load and validate application configuration
		// 加载和校验应用程序配置
		appConfig, err := application.LoadApplicationConfig(appConfigFile) // Assuming LoadApplicationConfig exists
		if err != nil {
			return fmt.Errorf("failed to load application config from %s: %w", appConfigFile, err)
		}

		// Load main platform configuration for full validation and Host K8s client
		// 加载主平台配置以进行完整校验和获取 Host K8s 客户端
		platformConfig, err := loader.LoadConfig(configFilePath)
		if err != nil {
			return fmt.Errorf("failed to load main platform config from %s for app upgrade: %w", configFilePath, err)
		}
		// Perform full validation
		// 执行完整校验
		// validator.ValidateApplicationConfigWithPlatform(appConfig, platformConfig) // Needs implementation

		// Get Host K8s client
		// 获取 Host K8s 客户端
		hostK8sClient, err := getHostK8sClient()
		if err != nil {
			return fmt.Errorf("failed to get host K8s client for app upgrade: %w", err)
		}

		// Create application deployer
		// 创建应用程序 deployer
		appDeployer := application.NewDeployer()

		// Run the upgrade process
		// 运行升级过程
		utils.GetLogger().Printf("Starting application upgrade for '%s' in vcluster '%s'", appConfig.Name, appConfig.VClusterName)
		if err := appDeployer.Upgrade(ctx, appConfig, hostK8sClient); err != nil {
			return fmt.Errorf("application upgrade failed for '%s': %w", appConfig.Name, err)
		}

		utils.GetLogger().Printf("Application '%s' upgraded successfully.", appConfig.Name)

		// TODO: Optional - Wait for application rollout if --wait flag is set
		// TODO: 可选 - 如果设置了 --wait 标志，则等待应用程序 rollout
		// waitFlag, _ := cmd.Flags().GetBool("wait")
		// if waitFlag {
		// 	utils.GetLogger().Printf("Waiting for application '%s' rollout to complete after upgrade...", appConfig.Name)
		// 	// err = application.WaitForRollout(ctx, vclusterClient, appConfig) // Needs helper
		// 	// if err != nil { return fmt.Errorf("application '%s' rollout failed after upgrade: %w", appConfig.Name, err) }
		// 	utils.GetLogger().Printf("Placeholder: Waiting for rollout for '%s' after upgrade.", appConfig.Name)
		// }

		return nil
	},
}

var applicationDeleteCmd = &cobra.Command{
	Use:   "delete <app-name>",
	Short: "Delete a business application",
	Long:  `Deletes a business application from a vcluster.`,
	Args:  cobra.ExactArgs(1), // App name // 应用程序名称
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(cmd.Context(), constants.DefaultTimeout*2) // Example timeout for delete // 示例删除超时时间
		defer cancel()

		appName := args[0]
		// Get flags
		vclusterName, _ := cmd.Flags().GetString("vcluster")
		namespace, _ := cmd.Flags().GetString("namespace")

		// Get Host K8s client
		// 获取 Host K8s 客户端
		hostK8sClient, err := getHostK8sClient() // Needs helper function to load kubeconfig
		if err != nil {
			return fmt.Errorf("failed to get host K8s client for app delete: %w", err)
		}

		// Create application deployer
		// 创建应用程序 deployer
		appDeployer := application.NewDeployer()

		// Run the deletion process
		// 运行删除过程
		utils.GetLogger().Printf("Starting application deletion for '%s' from vcluster '%s' in namespace '%s'", appName, vclusterName, namespace)
		if err := appDeployer.Delete(ctx, appName, vclusterName, namespace, hostK8sClient); err != nil { // Pass namespace
			return fmt.Errorf("application deletion failed for '%s': %w", appName, err)
		}

		utils.GetLogger().Printf("Application '%s' deletion initiated.", appName)

		// TODO: Optional - Wait for deletion if --wait flag is set
		// TODO: 可选 - 如果设置了 --wait 标志，则等待删除
		// waitFlag, _ := cmd.Flags().GetBool("wait")
		// if waitFlag {
		// 	utils.GetLogger().Printf("Waiting for application '%s' resources to be deleted...", appName)
		// 	// Requires getting vcluster client and checking if resources are gone
		// 	// vclusterClient, err := vcluster_client.GetVClusterClient(ctx, vclusterName, hostK8sClient)
		// 	// err = application.WaitForDeletion(ctx, vclusterClient, appName, namespace) // Needs helper
		// 	// if err != nil { return fmt.Errorf("application '%s' resources not fully deleted: %w", appName, err) }
		// 	utils.GetLogger().Printf("Placeholder: Waiting for deletion for '%s'.", appName)
		// }

		return nil
	},
}

var applicationStatusCmd = &cobra.Command{
	Use:   "status <app-name>",
	Short: "Get status of a business application",
	Long:  `Gets the deployment status of a business application in a vcluster.`,
	Args:  cobra.ExactArgs(1), // App name // 应用程序名称
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(cmd.Context(), constants.DefaultTimeout) // Example timeout for status // 示例状态超时时间
		defer cancel()

		appName := args[0]
		// Get flags
		vclusterName, _ := cmd.Flags().GetString("vcluster")
		namespace, _ := cmd.Flags().GetString("namespace")

		// Get Host K8s client
		// 获取 Host K8s 客户端
		hostK8sClient, err := getHostK8sClient() // Needs helper function to load kubeconfig
		if err != nil {
			return fmt.Errorf("failed to get host K8s client for app status: %w", err)
		}

		// Create application deployer
		// 创建应用程序 deployer
		appDeployer := application.NewDeployer()

		// Use appDeployer to get status
		// 使用 appDeployer 获取状态
		// status, err := appDeployer.GetStatus(ctx, appName, vclusterName, namespace, hostK8sClient) // Needs implementation in deployer // 需要在 deployer 中实现
		// if err != nil {
		// 	return fmt.Errorf("failed to get status for application '%s': %w", appName, err)
		// }

		// TODO: Format and print status
		// TODO: 格式化并打印状态
		utils.GetLogger().Printf("Placeholder: Getting status for application '%s' in vcluster '%s' (namespace '%s').", appName, vclusterName, namespace)
		// fmt.Printf("Application: %s\n", status.Name)
		// fmt.Printf("VCluster: %s\n", status.VCluster)
		// fmt.Printf("Namespace: %s\n", status.Namespace)
		// fmt.Printf("Ready: %t\n", status.Ready)
		// fmt.Printf("Message: %s\n", status.Message)

		return errors.New(errors.ErrTypeNotImplemented, fmt.Sprintf("application status command not fully implemented for '%s' in vcluster '%s'", appName, vclusterName))
	},
}

// getHostK8sClient is a helper function to obtain a Kubernetes client for the Host Cluster.
// getHostK8sClient 是一个辅助函数，用于获取 Host Cluster 的 Kubernetes 客户端。
// It needs to load the kubeconfig file, which should be available after platform deployment.
// 它需要加载 kubeconfig 文件，该文件在平台部署后应该可用。
// TODO: Implement robust kubeconfig loading logic (e.g., from standard locations, environment variable, or a specific file path).
// TODO: 实现健壮的 kubeconfig 加载逻辑（例如，从标准位置、环境变量或特定文件路径）。
func getHostK8sClient() (kubernetes.Interface, error) {
	utils.GetLogger().Println("Attempting to get Host Kubernetes client...")
	// Example: Load kubeconfig from a default path or environment variable
	// 示例：从默认路径或环境变量加载 kubeconfig
	// kubeconfigPath := os.Getenv("KUBECONFIG")
	// if kubeconfigPath == "" {
	// 	kubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config") // Standard default // 标准默认路径
	// }
	//
	// // Or load from a specific path managed by chasi-bod after deploy
	// // 或者从 chasi-bod 部署后管理的特定路径加载
	// // config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath) // Use clientcmd
	// // if err != nil { return nil, errors.NewWithCause(errors.ErrTypeSystem, fmt.Sprintf("failed to load kubeconfig from %s", kubeconfigPath), err) }
	// // client, err := kubernetes.NewForConfig(config)
	// // if err != nil { return nil, errors.NewWithCause(errors.ErrTypeSystem, "failed to create Host K8s client", err) }
	//
	// // Placeholder for successful client creation
	// // 成功创建客户端的占位符
	// utils.GetLogger().Println("Placeholder: Successfully obtained Host Kubernetes client.")
	// return nil, errors.New(errors.ErrTypeNotImplemented, "getting Host Kubernetes client not implemented yet")

	// For compilation purposes, return a dummy client if client-go is imported.
	// For real use, replace this with actual kubeconfig loading.
	// 出于编译目的，如果导入了 client-go，则返回一个虚拟客户端。
	// 对于实际使用，将其替换为实际的 kubeconfig 加载。
	// Need to import k8s.io/client-go/rest and k8s.io/apimachinery/pkg/runtime
	// 需要导入 k8s.io/client-go/rest 和 k8s.io/apimachinery/pkg/runtime
	dummyConfig := &rest.Config{Host: "http://localhost"} // Minimal dummy config // 最小虚拟配置
	dummyClient, _ := kubernetes.NewForConfig(dummyConfig)
	utils.GetLogger().Println("Warning: Using dummy Host Kubernetes client. Implement getHostK8sClient.")
	return dummyClient, errors.New(errors.ErrTypeNotImplemented, "getting Host Kubernetes client not implemented yet")
}

// TODO: Implement helper functions for application status checks and waiting for rollout/deletion
// TODO: 实现用于应用程序状态检查和等待 rollout/删除的辅助函数
// func WaitForRollout(ctx context.Context, vclusterClient kubernetes.Interface, appConfig *model.ApplicationConfig) error { ... }
// func WaitForDeletion(ctx context.Context, vclusterClient kubernetes.Interface, appName string, namespace string) error { ... }
