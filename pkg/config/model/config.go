// Package model defines the data structures for chasi-bod configuration files.
// 包 model 定义了 chasi-bod 配置文件的核心数据结构。
package model

import (
	"github.com/turtacn/chasi-bod/common/types"
	"github.com/turtacn/chasi-bod/common/types/enum"
)

// PlatformConfig represents the overall configuration for a chasi-bod platform instance.
// PlatformConfig 表示 chasi-bod 平台实例的总体配置。
type PlatformConfig struct {
	APIVersion   string                       `yaml:"apiVersion"`   // API version of the config schema / 配置模式的 API 版本
	Kind         string                       `yaml:"kind"`         // Kind of the config (e.g., "PlatformConfig") / 配置的类型（例如，“PlatformConfig”）
	Metadata     Metadata                     `yaml:"metadata"`     // Metadata for the configuration / 配置的元数据
	Cluster      ClusterConfig                `yaml:"cluster"`      // Host Cluster configuration / Host Cluster 配置
	VClusters    map[string]VClusterConfig    `yaml:"vclusters"`    // VCluster configurations by name / 按名称区分的 VCluster 配置
	Applications map[string]ApplicationConfig `yaml:"applications"` // Application configurations by name / 按名称区分的应用程序配置
	SysctlConfig types.SysctlConfig           `yaml:"sysctl"`       // Host OS sysctl configuration / Host OS sysctl 配置
	// Add other top-level configurations like DFX settings, etc.
	// 添加其他顶层配置，例如 DFX 设置等
	Output    OutputConfig `yaml:"output"` // Output configuration for the platform image / 平台镜像的输出配置
	DFXConfig DFXConfig    `yaml:"dfx"`    // DFX (Design for Excellence) configuration / DFX（卓越设计）配置
}

// Metadata contains metadata for the configuration.
// Metadata 包含配置的元数据。
type Metadata struct {
	Name        string            `yaml:"name"`        // Name of the platform configuration / 平台配置的名称
	Annotations map[string]string `yaml:"annotations"` // Optional annotations / 可选注解
	Labels      map[string]string `yaml:"labels"`      // 可选标签
}

// OutputConfig defines the configuration for generating the platform image.
// OutputConfig 定义了生成平台镜像的配置。
type OutputConfig struct {
	Format    enum.BuilderOutputFormat `yaml:"format"`    // Desired output format (e.g., iso, qcow2) / 期望的输出格式（例如，iso, qcow2）
	OutputDir string                   `yaml:"outputDir"` // Directory to save the output image / 保存输出镜像的目录
	// Add format-specific options as needed
	// 根据需要添加格式特定的选项
	ImageName string `yaml:"imageName"` // Name for the output image file / 输出镜像文件的名称
}

// ClusterConfig represents the configuration for the Host Cluster.
// ClusterConfig 表示 Host Cluster 的配置。
type ClusterConfig struct {
	Name              string              `yaml:"name"`              // Host Cluster name / Host Cluster 名称
	KubernetesVersion string              `yaml:"kubernetesVersion"` // Desired Kubernetes version for Host Cluster / Host Cluster 期望的 Kubernetes 版本
	ContainerRuntime  string              `yaml:"containerRuntime"`  // Desired Container runtime for Host Cluster / Host Cluster 期望的容器运行时
	Network           types.NetworkConfig `yaml:"network"`           // Network configuration / 网络配置
	Storage           types.StorageConfig `yaml:"storage"`           // Storage configuration / 存储配置
	Nodes             []NodeConfig        `yaml:"nodes"`             // Node configurations for deployment / 部署的节点配置
	// Add other host cluster specific configurations like apiserver cert sans etc.
	// 添加其他 Host Cluster 特定配置，例如 apiserver 证书 sans 等
	BaseOS BaseOSConfig `yaml:"baseOS"` // Base OS configuration for the image builder / 镜像构建器的基础操作系统配置
}

// BaseOSConfig represents the base OS configuration for the image builder.
// BaseOSConfig 表示镜像构建器的基础操作系统配置。
type BaseOSConfig struct {
	Image             string             `yaml:"image"`             // Base OS image (e.g., "ubuntu:22.04", "centos:stream8") / 基础操作系统镜像（例如，“ubuntu:22.04”、“centos:stream8”）
	Packages          []string           `yaml:"packages"`          // List of packages to install / 要安装的软件包列表
	KernelArgs        []string           `yaml:"kernelArgs"`        // Kernel boot arguments / 内核引导参数
	Files             []FileConfig       `yaml:"files"`             // List of files to copy into the image / 要复制到镜像中的文件列表
	Commands          []string           `yaml:"commands"`          // List of shell commands to run during build / 构建期间要运行的 shell 命令列表
	SysctlConfig      types.SysctlConfig `yaml:"sysctl"`            // Base OS sysctl configuration / 基础操作系统 sysctl 配置
	SSHAuthorizedKeys []string           `yaml:"sshAuthorizedKeys"` // SSH authorized keys to add / 要添加的 SSH 授权密钥
	Users             []UserConfig       `yaml:"users"`             // Users to create / 要创建的用户
}

// FileConfig represents a file to copy into the image during the build process.
// FileConfig 表示在构建过程中要复制到镜像中的文件。
type FileConfig struct {
	Source string `yaml:"source"` // Source path on the build machine / 构建机器上的源路径
	Dest   string `yaml:"dest"`   // Destination path in the target image filesystem / 目标镜像文件系统中的目标路径
	Mode   string `yaml:"mode"`   // File permissions (e.g., "0644", "0755") as octal string / 文件权限（例如，“0644”、“0755”），八进制字符串
}

// UserConfig represents a user to create in the image during the build process.
// UserConfig 表示在构建过程中要在镜像中创建的用户。
type UserConfig struct {
	Name     string   `yaml:"name"`     // Username / 用户名
	Password string   `yaml:"password"` // Password hash (use secure methods) / 密码哈希（使用安全方法）
	UID      *int     `yaml:"uid"`      // User ID (optional) / 用户 ID（可选）
	GID      *int     `yaml:"gid"`      // Group ID (optional) / 组 ID（可选）
	Groups   []string `yaml:"groups"`   // List of groups the user belongs to / 用户所属的组列表
	Shell    string   `yaml:"shell"`    // User's default shell / 用户的默认 shell
	HomeDir  string   `yaml:"homeDir"`  // User's home directory / 用户的家目录
	Sudo     bool     `yaml:"sudo"`     // Whether to grant sudo access / 是否授予 sudo 权限
}

// NodeConfig represents the configuration for a node during deployment.
// NodeConfig 表示部署期间节点的配置。
type NodeConfig struct {
	Address      string             `yaml:"address"`              // Node IP address or hostname / 节点 IP 地址或主机名
	User         string             `yaml:"user"`                 // SSH user for deployment / 用于部署的 SSH 用户
	Password     string             `yaml:"password,omitempty"`   // SSH password (discouraged, use PrivateKey) / SSH 密码（不推荐，使用 PrivateKey）
	PrivateKey   string             `yaml:"privateKey,omitempty"` // SSH private key path / SSH 私钥路径
	Port         int                `yaml:"port"`                 // SSH port / SSH 端口
	Roles        []enum.NodeRole    `yaml:"roles"`                // Node roles / 节点角色
	Labels       map[string]string  `yaml:"labels"`               // Desired Kubernetes node labels / 期望的 Kubernetes 节点标签
	Annotations  map[string]string  `yaml:"annotations"`          // Desired Kubernetes node annotations / 期望的 Kubernetes 节点注解
	Taints       []string           `yaml:"taints"`               // Desired Kubernetes node taints / 期望的 Kubernetes 节点污点
	SysctlConfig types.SysctlConfig `yaml:"sysctl"`               // Node-specific sysctl configuration / 节点特定的 sysctl 配置
	// Add more node specific configurations like disks, mount points etc.
	// 添加更多节点特定配置，例如磁盘、挂载点等
	DiskConfigs []DiskConfig `yaml:"diskConfigs"` // Disk configurations for partitioning and mounting / 磁盘配置用于分区和挂载
}

// DiskConfig represents configuration for a disk on a node.
// DiskConfig 表示节点上磁盘的配置。
type DiskConfig struct {
	Device string `yaml:"device"` // Disk device name (e.g., "/dev/sda") / 磁盘设备名称（例如，“/dev/sda”）
	// Add partitioning, filesystem, and mount point configurations
	// 添加分区、文件系统和挂载点配置
	Filesystem string `yaml:"filesystem"` // Filesystem type (e.g., "ext4", "xfs") / 文件系统类型（例如，“ext4”、“xfs”）
	MountPoint string `yaml:"mountPoint"` // Mount point (e.g., "/var/lib/chasi-bod") / 挂载点（例如，“/var/lib/chasi-bod”）
	Format     bool   `yaml:"format"`     // Whether to format the disk / 是否格式化磁盘
}

// VClusterConfig represents the configuration for a virtual cluster.
// VClusterConfig 表示虚拟集群的配置。
type VClusterConfig struct {
	Name              string            `yaml:"name"`              // VCluster name / VCluster 名称
	Namespace         string            `yaml:"namespace"`         // Host Cluster namespace where vcluster runs / vcluster 在 host 集群中运行的命名空间
	Template          string            `yaml:"template"`          // Name of the vcluster template to use (optional, for common configs) / 使用的 vcluster 模板名称（可选，用于通用配置）
	KubernetesVersion string            `yaml:"kubernetesVersion"` // Kubernetes version for the vcluster / vcluster 的 Kubernetes 版本
	ServiceCIDR       string            `yaml:"serviceCIDR"`       // Service network CIDR for the vcluster / vcluster 的 Service 网络 CIDR
	PodCIDR           string            `yaml:"podCIDR"`           // Pod network CIDR for the vcluster / vcluster 的 Pod 网络 CIDR
	ResourceRequests  map[string]string `yaml:"resourceRequests"`  // Resource requests for vcluster control plane pods in Host Cluster / vcluster 控制面 Pod 在 Host 集群中的资源请求
	ResourceLimits    map[string]string `yaml:"resourceLimits"`    // Resource limits for vcluster control plane pods in Host Cluster / vcluster 控制面 Pod 在 Host 集群中的资源限制
	// Add more vcluster specific configurations based on loft-sh/vcluster options
	// 根据 loft-sh/vcluster 选项添加更多 vcluster 特定配置
	Sync SyncConfig `yaml:"sync"` // vcluster syncer configuration / vcluster syncer 配置
	// Add vcluster specific network/storage configs if they override host ones
	// 如果 vcluster 特定的网络/存储配置覆盖了 host 配置，则添加这些配置
	Network *types.NetworkConfig `yaml:"network,omitempty"` // Optional vcluster network config / 可选的 vcluster 网络配置
	Storage *types.StorageConfig `yaml:"storage,omitempty"` // 可选的 vcluster 存储配置

}

// SyncConfig represents the vcluster syncer configuration.
// SyncConfig 表示 vcluster syncer 配置。
type SyncConfig struct {
	EnabledResources []string `yaml:"enabledResources"` // List of resource types to sync / 要同步的资源类型列表
	// Add resource-specific sync options as needed
	// 根据需要添加资源特定的同步选项
	PodPodLifecycle           string                          `yaml:"podPodLifecycle"`           // Pod lifecycle sync strategy / Pod 生命周期同步策略
	PersistentVolumeClaimSync PersistentVolumeClaimSyncConfig `yaml:"persistentVolumeClaimSync"` // PVC sync config / PVC 同步配置
	// Add other syncer configurations
	// 添加其他 syncer 配置
}

// PersistentVolumeClaimSyncConfig represents PVC sync configuration.
// PersistentVolumeClaimSyncConfig 表示 PVC 同步配置。
type PersistentVolumeClaimSyncConfig struct {
	Enabled bool `yaml:"enabled"` // Enable PVC sync / 启用 PVC 同步
	// Add PVC specific sync options
	// 添加 PVC 特定同步选项
}

// ApplicationConfig represents the configuration for a business application deployment.
// ApplicationConfig 表示业务应用部署的配置。
type ApplicationConfig struct {
	Name         string               `yaml:"name"`         // Application name / 应用程序名称
	VClusterName string               `yaml:"vclusterName"` // Name of the target vcluster / 目标 vcluster 名称
	Namespace    string               `yaml:"namespace"`    // Namespace within the vcluster to deploy to / 部署到 vcluster 内的命名空间
	Type         enum.ApplicationType `yaml:"type"`         // Application type (e.g., CPU-bound, IO-bound) / 应用程序类型（例如，CPU 密集型、IO 密集型）
	// Define how the application manifests are specified
	// 定义如何指定应用程序 manifests
	HelmChart *HelmChartConfig `yaml:"helmChart,omitempty"` // Helm chart configuration / Helm chart 配置
	Kustomize *KustomizeConfig `yaml:"kustomize,omitempty"` // Kustomize configuration / Kustomize 配置
	Manifests []string         `yaml:"manifests,omitempty"` // List of raw manifest file paths or URLs / 原始 manifest 文件路径或 URL 列表
	// Add more application specific configurations like resource requests/limits defaults, config injection etc.
	// 添加更多应用程序特定配置，例如资源请求/限制默认值、配置注入等
	ConfigInjection map[string]string `yaml:"configInjection,omitempty"` // Map of config keys to inject (e.g., database connection string from a central store) / 要注入的配置键映射（例如，从中央存储获取的数据库连接字符串）
}

// HelmChartConfig represents configuration for deploying a Helm chart.
// HelmChartConfig 表示部署 Helm chart 的配置。
type HelmChartConfig struct {
	Chart       string                 `yaml:"chart"`       // Chart name or path / Chart 名称或路径
	Repo        string                 `yaml:"repo"`        // Chart repository URL (if not local) / Chart 仓库 URL（如果不是本地）
	Version     string                 `yaml:"version"`     // Chart version / Chart 版本
	Values      map[string]interface{} `yaml:"values"`      // Values to override / 要覆盖的值
	ReleaseName string                 `yaml:"releaseName"` // Helm release name / Helm release 名称
	// Add other helm specific options
	// 添加其他 helm 特定选项
}

// KustomizeConfig represents configuration for deploying with Kustomize.
// KustomizeConfig 表示使用 Kustomize 部署的配置。
type KustomizeConfig struct {
	Path string `yaml:"path"` // Path to the kustomization directory / kustomization 目录的路径
	// Add other kustomize specific options
	// 添加其他 kustomize 特定选项
}

// DFXConfig represents configurations related to DFX.
// DFXConfig 表示与 DFX 相关的配置。
type DFXConfig struct {
	Logging     LoggingConfig     `yaml:"logging"`     // Logging configuration / 日志配置
	Metrics     MetricsConfig     `yaml:"metrics"`     // Metrics configuration / 指标配置
	Tracing     TracingConfig     `yaml:"tracing"`     // Tracing configuration / 追踪配置
	Healthz     HealthzConfig     `yaml:"healthz"`     // Health check configuration / 健康检查配置
	Reliability ReliabilityConfig `yaml:"reliability"` // Reliability configuration / 可靠性配置
	// Add more DFX aspects
	// 添加更多 DFX 方面
}

// LoggingConfig represents logging configurations.
// LoggingConfig 表示日志配置。
type LoggingConfig struct {
	Enabled bool   `yaml:"enabled"` // Enable logging collection / 启用日志收集
	Agent   string `yaml:"agent"`   // Logging agent type (e.g., "fluentd", "fluentbit") / 日志代理类型（例如，“fluentd”、“fluentbit”）
	// Add agent specific configs, output destinations etc.
	// 添加代理特定配置、输出目的地等
	Output OutputDestinationConfig `yaml:"output"` // Logging output destination / 日志输出目的地
}

// OutputDestinationConfig represents configuration for a log output destination.
// OutputDestinationConfig 表示日志输出目的地的配置。
type OutputDestinationConfig struct {
	Type     string `yaml:"type"`     // Destination type (e.g., "elasticsearch", "kafka", "stdout") / 目的地类型（例如，“elasticsearch”、“kafka”、“stdout”）
	Endpoint string `yaml:"endpoint"` // Endpoint address / 端点地址
	// Add type-specific configurations
	// 添加类型特定配置
	Options map[string]string `yaml:"options"` // Additional options / 额外选项
}

// MetricsConfig represents metrics configurations.
// MetricsConfig 表示指标配置。
type MetricsConfig struct {
	Enabled       bool                 `yaml:"enabled"`       // Enable metrics collection / 启用指标收集
	Agent         string               `yaml:"agent"`         // Metrics agent type (e.g., "prometheus-agent", "agent-based") / 指标代理类型（例如，“prometheus-agent”、“基于代理的”）
	ScrapeConfigs []MetricScrapeConfig `yaml:"scrapeConfigs"` // Prometheus scrape configurations / Prometheus scrape 配置
	// Add agent specific configs, storage destinations etc.
	// 添加代理特定配置、存储目的地等
}

// MetricScrapeConfig represents a Prometheus scrape configuration.
// MetricScrapeConfig 表示 Prometheus scrape 配置。
type MetricScrapeConfig struct {
	JobName string `yaml:"jobName"` // Scrape job name / Scrape 作业名称
	// Add scrape target configuration
	// 添加 scrape 目标配置
	StaticConfigs []map[string][]string `yaml:"static_configs"` // Static targets / 静态目标
	// Add other scrape config options
	// 添加其他 scrape 配置选项
}

// TracingConfig represents tracing configurations.
// TracingConfig 表示追踪配置。
type TracingConfig struct {
	Enabled  bool   `yaml:"enabled"`  // Enable distributed tracing / 启用分布式追踪
	Agent    string `yaml:"agent"`    // Tracing agent type (e.g., "jaeger", "zipkin") / 追踪代理类型（例如，“jaeger”、“zipkin”）
	Endpoint string `yaml:"endpoint"` // Tracing collector endpoint / 追踪收集器端点
	// Add agent specific configs, sampling rates etc.
	// 添加代理特定配置、采样率等
}

// HealthzConfig represents health check configurations.
// HealthzConfig 表示健康检查配置。
type HealthzConfig struct {
	Enabled bool `yaml:"enabled"` // Enable platform-level health checks / 启用平台级健康检查
	// Add specific probes or checks
	// 添加特定的探针或检查
	Interval string `yaml:"interval"` // Check interval (e.g., "30s") / 检查间隔（例如，“30s”）
	Timeout  string `yaml:"timeout"`  // Check timeout (e.g., "5s") / 检查超时（例如，“5s”）
}

// ReliabilityConfig represents reliability configurations.
// ReliabilityConfig 表示可靠性配置。
type ReliabilityConfig struct {
	ConfigBackup ConfigBackupConfig `yaml:"configBackup"` // Configuration backup settings / 配置备份设置
	ETCDBackup   ETCDBackupConfig   `yaml:"etcdBackup"`   // ETCD backup settings for Host Cluster / Host Cluster 的 ETCD 备份设置
	// Add other reliability features
	// 添加其他可靠性功能
}

// ConfigBackupConfig represents configuration backup settings.
// ConfigBackupConfig 表示配置备份设置。
type ConfigBackupConfig struct {
	Enabled  bool   `yaml:"enabled"`  // Enable configuration backup / 启用配置备份
	Schedule string `yaml:"schedule"` // Backup schedule (e.g., "@daily" or cron string) / 备份计划（例如，“@daily”或 cron 字符串）
	Location string `yaml:"location"` // Backup storage location (e.g., S3 bucket, NFS path) / 备份存储位置（例如，S3 bucket，NFS 路径）
}

// ETCDBackupConfig represents ETCD backup settings.
// ETCDBackupConfig 表示 ETCD 备份设置。
type ETCDBackupConfig struct {
	Enabled  bool   `yaml:"enabled"`  // Enable ETCD backup / 启用 ETCD 备份
	Schedule string `yaml:"schedule"` // Backup schedule / 备份计划
	Location string `yaml:"location"` // Backup storage location / 备份存储位置
	// Add ETCD specific backup options
	// 添加 ETCD 特定备份选项
}
