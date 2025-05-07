// Package validator provides functionality to validate chasi-bod configuration.
// 包 validator 提供了校验 chasi-bod 配置的功能。
package validator

import (
	"fmt"
	"net"
	//"time"

	"github.com/turtacn/chasi-bod/common/errors"
	"github.com/turtacn/chasi-bod/common/types"
	"github.com/turtacn/chasi-bod/common/types/enum"
	"github.com/turtacn/chasi-bod/common/utils" // Assuming utils are needed for IP/Hostname validation and logger
	"github.com/turtacn/chasi-bod/pkg/config/model"
)

// ValidateConfig validates the entire PlatformConfig structure.
// ValidateConfig 校验整个 PlatformConfig 结构体。
// config: The PlatformConfig struct to validate. / 要校验的 PlatformConfig 结构体。
// Returns an error if validation fails.
// 如果校验失败则返回错误。
func ValidateConfig(config *model.PlatformConfig) error {
	if config == nil {
		return errors.New(errors.ErrTypeValidation, "platform configuration is nil")
	}

	// Validate Metadata
	// 校验元数据
	if config.Metadata.Name == "" {
		return errors.New(errors.ErrTypeValidation, "metadata.name is required")
	}

	// Validate OutputConfig
	// 校验输出配置
	if err := validateOutputConfig(&config.Output); err != nil {
		return fmt.Errorf("invalid output configuration: %w", err)
	}

	// Validate ClusterConfig
	// 校验集群配置
	if err := validateClusterConfig(&config.Cluster); err != nil {
		return fmt.Errorf("invalid cluster configuration: %w", err)
	}

	// Validate SysctlConfig (basic check)
	// 校验 Sysctl 配置（基本检查）
	if err := validateSysctlConfig(config.SysctlConfig); err != nil {
		return fmt.Errorf("invalid global sysctl configuration: %w", err)
	}

	// Validate VClusterConfig map
	// 校验 VCluster 配置映射
	if len(config.VClusters) > 0 {
		for name, vclusterCfg := range config.VClusters {
			if err := validateVClusterConfig(name, &vclusterCfg, config.Cluster.KubernetesVersion); err != nil {
				return fmt.Errorf("invalid vcluster configuration '%s': %w", name, err)
			}
		}
	}

	// Validate ApplicationConfig map
	// 校验应用程序配置映射
	if len(config.Applications) > 0 {
		for name, appCfg := range config.Applications {
			// Check if the target vcluster exists
			// 检查目标 vcluster 是否存在
			if _, exists := config.VClusters[appCfg.VClusterName]; !exists {
				return errors.New(errors.ErrTypeValidation, fmt.Sprintf("application '%s' targets non-existent vcluster '%s'", name, appCfg.VClusterName))
			}
			if err := validateApplicationConfig(name, &appCfg); err != nil {
				return fmt.Errorf("invalid application configuration '%s': %w", name, err)
			}
		}
	}

	// Validate DFXConfig
	// 校验 DFX 配置
	if err := validateDFXConfig(&config.DFXConfig); err != nil {
		return fmt.Errorf("invalid dfx configuration: %w", err)
	}

	// TODO: Add cross-resource validation (e.g., check for port conflicts across vclusters/host)
	// TODO: 添加跨资源校验（例如，检查 vcluster/host 之间的端口冲突）

	return nil
}

// validateOutputConfig validates the OutputConfig.
// validateOutputConfig 校验 OutputConfig。
func validateOutputConfig(config *model.OutputConfig) error {
	if config.OutputDir == "" {
		return errors.New(errors.ErrTypeValidation, "output.outputDir is required")
	}
	// Basic format validation (check if it's a known enum value)
	// 基本格式校验（检查是否是已知枚举值）
	switch config.Format {
	case enum.OutputFormatISO, enum.OutputFormatQCOW2, enum.OutputFormatOVA, enum.OutputFormatVMA:
		// Valid formats
		// 有效格式
	default:
		return errors.New(errors.ErrTypeValidation, fmt.Sprintf("invalid output.format '%s'", config.Format))
	}
	if config.ImageName == "" {
		// Default image name could be derived if not provided, but making it required simplifies things for now.
		// 如果未提供，可以派生默认镜像名称，但这会简化当前的情况，因此将其设为必需。
		return errors.New(errors.ErrTypeValidation, "output.imageName is required")
	}
	return nil
}

// validateClusterConfig validates the ClusterConfig.
// validateClusterConfig 校验 ClusterConfig。
func validateClusterConfig(config *model.ClusterConfig) error {
	if config.Name == "" {
		return errors.New(errors.ErrTypeValidation, "cluster.name is required")
	}
	if config.KubernetesVersion == "" {
		return errors.New(errors.ErrTypeValidation, "cluster.kubernetesVersion is required")
	}
	// TODO: Add stricter version format validation
	// TODO: 添加更严格的版本格式校验

	// Validate NetworkConfig
	// 校验网络配置
	if err := validateNetworkConfig(&config.Network); err != nil {
		return fmt.Errorf("invalid cluster.network configuration: %w", err)
	}

	// Validate StorageConfig
	// 校验存储配置
	if err := validateStorageConfig(&config.Storage); err != nil {
		return fmt.Errorf("invalid cluster.storage configuration: %w", err)
	}

	// Validate NodeConfig list
	// 校验节点配置列表
	if len(config.Nodes) == 0 {
		return errors.New(errors.ErrTypeValidation, "cluster.nodes list is required and cannot be empty")
	}
	hasMaster := false
	for _, nodeCfg := range config.Nodes {
		if err := validateNodeConfig(&nodeCfg); err != nil {
			return fmt.Errorf("invalid node configuration for address '%s': %w", nodeCfg.Address, err)
		}
		for _, role := range nodeCfg.Roles {
			if role == enum.RoleMaster {
				hasMaster = true
			}
		}
	}
	if !hasMaster {
		return errors.New(errors.ErrTypeValidation, "cluster.nodes must contain at least one node with role 'master'")
	}

	// Validate BaseOSConfig
	// 校验基础操作系统配置
	if err := validateBaseOSConfig(&config.BaseOS); err != nil {
		return fmt.Errorf("invalid cluster.baseOS configuration: %w", err)
	}

	return nil
}

// validateBaseOSConfig validates the BaseOSConfig.
// validateBaseOSConfig 校验基础操作系统配置。
func validateBaseOSConfig(config *model.BaseOSConfig) error {
	if config.Image == "" {
		return errors.New(errors.ErrTypeValidation, "cluster.baseOS.image is required")
	}
	// TODO: Add image format/syntax validation
	// TODO: 添加镜像格式/语法校验

	// Validate SysctlConfig (basic check)
	// 校验 Sysctl 配置（基本检查）
	if err := validateSysctlConfig(config.SysctlConfig); err != nil {
		return fmt.Errorf("invalid cluster.baseOS.sysctl configuration: %w", err)
	}

	// Validate Files
	// 校验文件配置
	for _, fileCfg := range config.Files {
		if fileCfg.Source == "" || fileCfg.Dest == "" {
			return errors.New(errors.ErrTypeValidation, "cluster.baseOS.files requires source and dest paths")
		}
		// TODO: Validate file permissions format
		// TODO: 校验文件权限格式
	}

	// Validate Users
	// 校验用户配置
	for _, userCfg := range config.Users {
		if userCfg.Name == "" {
			return errors.New(errors.ErrTypeValidation, "cluster.baseOS.users requires user name")
		}
		// TODO: Add password/ssh key requirements based on usage
		// TODO: 根据使用情况添加密码/ssh 密钥要求
	}

	return nil
}

// validateNodeConfig validates a single NodeConfig.
// validateNodeConfig 校验单个 NodeConfig。
func validateNodeConfig(config *model.NodeConfig) error {
	if config.Address == "" {
		return errors.New(errors.ErrTypeValidation, "node address is required")
	}
	// Validate address format (IP or hostname)
	// 校验地址格式（IP 或主机名）
	if !utils.IsValidIPAddress(config.Address) && !utils.IsValidHostname(config.Address) {
		return errors.New(errors.ErrTypeValidation, fmt.Sprintf("invalid node address format '%s'", config.Address))
	}

	if config.User == "" {
		return errors.New(errors.ErrTypeValidation, fmt.Sprintf("node %s: user is required", config.Address))
	}

	if config.Password != "" && config.PrivateKey != "" {
		return errors.New(errors.ErrTypeValidation, fmt.Sprintf("node %s: cannot specify both password and privateKey", config.Address))
	}
	if config.Password == "" && config.PrivateKey == "" {
		// Decide if one is required or if other auth methods are implicit
		// 决定是否必需其中之一，或者是否隐式支持其他身份验证方法
		// For now, let's assume one is required for SSH
		// 现在，我们假设 SSH 需要其中之一
		return errors.New(errors.ErrTypeValidation, fmt.Sprintf("node %s: either password or privateKey is required for SSH", config.Address))
	}

	if config.Port == 0 {
		config.Port = 22 // Default SSH port
	}
	if config.Port < 1 || config.Port > 65535 {
		return errors.New(errors.ErrTypeValidation, fmt.Sprintf("node %s: invalid port number %d", config.Address, config.Port))
	}

	if len(config.Roles) == 0 {
		return errors.New(errors.ErrTypeValidation, fmt.Sprintf("node %s: roles are required and cannot be empty", config.Address))
	}
	for _, role := range config.Roles {
		switch role {
		case enum.RoleMaster, enum.RoleWorker, enum.RoleEdge:
			// Valid role
			// 有效角色
		default:
			return errors.New(errors.ErrTypeValidation, fmt.Sprintf("node %s: invalid role '%s'", config.Address, role))
		}
	}

	// Validate Node-specific SysctlConfig (basic check)
	// 校验节点特定的 Sysctl 配置（基本检查）
	if err := validateSysctlConfig(config.SysctlConfig); err != nil {
		return fmt.Errorf("node %s: invalid sysctl configuration: %w", config.Address, err)
	}

	// Validate DiskConfigs
	// 校验磁盘配置
	for _, diskCfg := range config.DiskConfigs {
		if diskCfg.Device == "" {
			return errors.New(errors.ErrTypeValidation, fmt.Sprintf("node %s: disk config requires device name", config.Address))
		}
		if diskCfg.MountPoint != "" && diskCfg.Filesystem == "" {
			return errors.New(errors.ErrTypeValidation, fmt.Sprintf("node %s, disk %s: filesystem is required if mountPoint is specified", config.Address, diskCfg.Device))
		}
		// TODO: Add stricter validation for device names, filesystems, mount points
		// TODO: 对设备名称、文件系统、挂载点添加更严格的校验
	}

	return nil
}

// validateNetworkConfig validates the NetworkConfig.
// validateNetworkConfig 校验 NetworkConfig。
func validateNetworkConfig(config *types.NetworkConfig) error {
	if config.Plugin == "" {
		return errors.New(errors.ErrTypeValidation, "network.plugin is required")
	}
	// TODO: Add validation for known CNI plugins

	if config.PodCIDR != "" {
		if _, _, err := net.ParseCIDR(config.PodCIDR); err != nil {
			return errors.New(errors.ErrTypeValidation, fmt.Sprintf("invalid network.podCIDR format '%s': %v", config.PodCIDR, err))
		}
	}
	if config.ServiceCIDR != "" {
		if _, _, err := net.ParseCIDR(config.ServiceCIDR); err != nil {
			return errors.New(errors.ErrTypeValidation, fmt.Sprintf("invalid network.serviceCIDR format '%s': %v", config.ServiceCIDR, err))
		}
	}

	// TODO: Add validation for interfaces, IP addresses, gateways etc.

	return nil
}

// validateStorageConfig validates the StorageConfig.
// validateStorageConfig 校验 StorageConfig。
func validateStorageConfig(config *types.StorageConfig) error {
	// DefaultStorageClass is optional
	// 默认 StorageClass 是可选的

	// Validate StorageClasses
	// 校验 StorageClass
	for _, scCfg := range config.StorageClasses {
		if scCfg.Name == "" || scCfg.Provisioner == "" {
			return errors.New(errors.ErrTypeValidation, "storage.storageClasses requires name and provisioner")
		}
		// TODO: Add validation for provisioner format
		// TODO: 校验 provisioner 格式
	}

	// Validate PVConfigs
	// 校验 PV 配置
	for _, pvCfg := range config.PVConfigs {
		if pvCfg.Name == "" || pvCfg.Capacity == "" {
			return errors.New(errors.ErrTypeValidation, "storage.pvConfigs requires name and capacity")
		}
		// TODO: Validate capacity format (e.g., "10Gi")
		// TODO: 校验容量格式（例如，“10Gi”）
		if len(pvCfg.AccessModes) == 0 {
			return errors.New(errors.ErrTypeValidation, fmt.Sprintf("storage.pvConfigs '%s' requires at least one accessMode", pvCfg.Name))
		}
		// TODO: Validate access mode values (e.g., "ReadWriteOnce")
		// TODO: 校验访问模式值（例如，“ReadWriteOnce”）
		// TODO: Validate PersistentVolumeSource configuration based on type
		// TODO: 根据类型校验 PersistentVolumeSource 配置
	}

	return nil
}

// validateSysctlConfig performs a basic validation on a SysctlConfig map.
// validateSysctlConfig 对 SysctlConfig 映射执行基本校验。
func validateSysctlConfig(config types.SysctlConfig) error {
	// Basic check: keys should not be empty
	// 基本检查：键不能为空
	for param := range config {
		if param == "" {
			return errors.New(errors.ErrTypeValidation, "sysctl parameter name cannot be empty")
		}
		// TODO: Add regex or format validation for parameter names if needed
		// TODO: 如果需要，添加参数名称的正则表达式或格式校验
	}
	// Note: More robust validation might require checking if the parameter exists on a target OS,
	// which is hard to do during config validation time.
	// 注意：更强大的校验可能需要检查参数是否存在于目标操作系统上，这在配置校验时很难做到。
	return nil
}

// validateVClusterConfig validates a single VClusterConfig.
// validateVClusterConfig 校验单个 VClusterConfig。
// name: The name of the vcluster being validated. / 正在校验的 vcluster 的名称。
// hostK8sVersion: The version of the host Kubernetes cluster. / Host Kubernetes 集群的版本。
func validateVClusterConfig(name string, config *model.VClusterConfig, hostK8sVersion string) error {
	if config.Name == "" {
		// Should not happen if validating map keys, but good check
		// 应该不会发生如果校验映射键，但这是一个好的检查
		return errors.New(errors.ErrTypeValidation, "vcluster name is required")
	}
	if config.Name != name {
		return errors.New(errors.ErrTypeValidation, fmt.Sprintf("vcluster name in config ('%s') must match map key ('%s')", config.Name, name))
	}
	if config.Namespace == "" {
		return errors.New(errors.ErrTypeValidation, fmt.Sprintf("vcluster '%s': namespace is required", name))
	}
	// TODO: Validate namespace format

	if config.KubernetesVersion == "" {
		// If not specified, vcluster might default to host K8s version - decide on behavior
		// 如果未指定，vcluster 可能会默认使用 host K8s 版本 - 根据行为决定
		// Let's make it required for clarity for now.
		// 现在为了清晰起见，我们将其设为必需。
		return errors.New(errors.ErrTypeValidation, fmt.Sprintf("vcluster '%s': kubernetesVersion is required", name))
	}
	// TODO: Validate vcluster KubernetesVersion compatibility with HostK8sVersion
	// TODO: 校验 vcluster KubernetesVersion 与 HostK8sVersion 的兼容性

	// Validate CIDRs
	// 校验 CIDR
	if config.ServiceCIDR != "" {
		if _, _, err := net.ParseCIDR(config.ServiceCIDR); err != nil {
			return errors.New(errors.ErrTypeValidation, fmt.Sprintf("vcluster '%s': invalid serviceCIDR format '%s': %v", name, config.ServiceCIDR, err))
		}
	}
	if config.PodCIDR != "" {
		if _, _, err := net.ParseCIDR(config.PodCIDR); err != nil {
			return errors.New(errors.ErrTypeValidation, fmt.Sprintf("vcluster '%s': invalid podCIDR format '%s': %v", name, config.PodCIDR, err))
		}
	}
	// TODO: Add checks to ensure vcluster CIDRs do not conflict with Host CIDRs

	// TODO: Validate resource requests/limits format

	// TODO: Validate SyncConfig

	// Validate optional network/storage configs if present
	// 校验可选的网络/存储配置（如果存在）
	if config.Network != nil {
		if err := validateNetworkConfig(config.Network); err != nil {
			return fmt.Errorf("vcluster '%s': invalid network configuration: %w", name, err)
		}
	}
	if config.Storage != nil {
		if err := validateStorageConfig(config.Storage); err != nil {
			return fmt.Errorf("vcluster '%s': invalid storage configuration: %w", name, err)
		}
	}

	return nil
}

// validateApplicationConfig validates a single ApplicationConfig.
// validateApplicationConfig 校验单个 ApplicationConfig。
// name: The name of the application being validated. / 正在校验的应用程序的名称。
func validateApplicationConfig(name string, config *model.ApplicationConfig) error {
	if config.Name == "" {
		// Should not happen if validating map keys, but good check
		// 应该不会发生如果校验映射键，但这是一个好的检查
		return errors.New(errors.ErrTypeValidation, "application name is required")
	}
	if config.Name != name {
		return errors.New(errors.ErrTypeValidation, fmt.Sprintf("application name in config ('%s') must match map key ('%s')", config.Name, name))
	}
	if config.VClusterName == "" {
		return errors.New(errors.ErrTypeValidation, fmt.Sprintf("application '%s': vclusterName is required", name))
	}
	if config.Namespace == "" {
		config.Namespace = "default" // Default namespace if not specified
	}
	// TODO: Validate namespace format

	// Validate application type
	// 校验应用程序类型
	switch config.Type {
	case enum.AppTypeComputeBound, enum.AppTypeIOBound, enum.AppTypeMemoryBound, enum.AppTypeNetworkBound, enum.AppTypeGeneral:
		// Valid type
		// 有效类型
	default:
		return errors.New(errors.ErrTypeValidation, fmt.Sprintf("application '%s': invalid type '%s'", name, config.Type))
	}

	// Ensure at least one way to specify manifests is provided
	// 确保至少提供了一种指定 manifests 的方式
	if config.HelmChart == nil && config.Kustomize == nil && len(config.Manifests) == 0 {
		return errors.New(errors.ErrTypeValidation, fmt.Sprintf("application '%s': one of helmChart, kustomize, or manifests must be specified", name))
	}
	// Ensure only one way is specified (optional, depends on desired flexibility)
	// 确保只指定了一种方式（可选，取决于期望的灵活性）
	manifestSources := 0
	if config.HelmChart != nil {
		manifestSources++
		// Validate Helm chart config
		// 校验 Helm chart 配置
		if config.HelmChart.Chart == "" || config.HelmChart.ReleaseName == "" {
			return errors.New(errors.ErrTypeValidation, fmt.Sprintf("application '%s': helmChart requires chart and releaseName", name))
		}
		// TODO: Validate repo/chart format
		// TODO: 校验 repo/chart 格式
	}
	if config.Kustomize != nil {
		manifestSources++
		// Validate Kustomize config
		// 校验 Kustomize 配置
		if config.Kustomize.Path == "" {
			return errors.New(errors.ErrTypeValidation, fmt.Sprintf("application '%s': kustomize requires path", name))
		}
		// TODO: Validate path existence/format
		// TODO: 校验路径存在性/格式
	}
	if len(config.Manifests) > 0 {
		manifestSources++
		// Validate manifest paths/URLs
		// 校验 manifest 路径/URL
		for _, manifestPath := range config.Manifests {
			if manifestPath == "" {
				return errors.New(errors.ErrTypeValidation, fmt.Sprintf("application '%s': manifest path cannot be empty", name))
			}
			// TODO: Validate path/URL format or existence
			// TODO: 校验路径/URL 格式或存在性
		}
	}
	if manifestSources > 1 {
		return errors.New(errors.ErrTypeValidation, fmt.Sprintf("application '%s': only one of helmChart, kustomize, or manifests can be specified", name))
	}

	// TODO: Validate ConfigInjection format or references

	return nil
}

// validateDFXConfig validates the DFXConfig.
// validateDFXConfig 校验 DFXConfig。
func validateDFXConfig(config *model.DFXConfig) error {
	if config == nil {
		// DFX config is optional at the top level, but if present, its sub-configs should be checked
		// DFX 配置在顶层是可选的，但如果存在，应检查其子配置
		return nil
	}

	// Validate LoggingConfig
	// 校验日志配置
	if config.Logging.Enabled {
		if config.Logging.Agent == "" {
			return errors.New(errors.ErrTypeValidation, "dfx.logging.agent is required when logging is enabled")
		}
		// TODO: Add validation for known logging agents
		// Validate output destination if logging is enabled
		// 如果启用了日志记录，校验输出目的地
		if config.Logging.Output.Type == "" {
			return errors.New(errors.ErrTypeValidation, "dfx.logging.output.type is required when logging is enabled")
		}
		// TODO: Add validation for known output types and endpoint format
	}

	// Validate MetricsConfig
	// 校验指标配置
	if config.Metrics.Enabled {
		if config.Metrics.Agent == "" {
			return errors.New(errors.ErrTypeValidation, "dfx.metrics.agent is required when metrics is enabled")
		}
		// TODO: Add validation for known metrics agents
		// ScrapeConfigs are optional even if metrics is enabled
	}

	// Validate TracingConfig
	// 校验追踪配置
	if config.Tracing.Enabled {
		if config.Tracing.Agent == "" {
			return errors.New(errors.ErrTypeValidation, "dfx.tracing.agent is required when tracing is enabled")
		}
		if config.Tracing.Endpoint == "" {
			return errors.New(errors.ErrTypeValidation, "dfx.tracing.endpoint is required when tracing is enabled")
		}
		// TODO: Add validation for known tracing agents and endpoint format
	}

	// Validate HealthzConfig
	// 校验健康检查配置
	if config.Healthz.Enabled {
		// Interval and Timeout format validation
		// 间隔和超时格式校验
		if config.Healthz.Interval != "" {
			if _, err := time.ParseDuration(config.Healthz.Interval); err != nil {
				return errors.New(errors.ErrTypeValidation, fmt.Sprintf("invalid dfx.healthz.interval format '%s': %v", config.Healthz.Interval, err))
			}
		}
		if config.Healthz.Timeout != "" {
			if _, err := time.ParseDuration(config.Healthz.Timeout); err != nil {
				return errors.New(errors.ErrTypeValidation, fmt.Sprintf("invalid dfx.healthz.timeout format '%s': %v", config.Healthz.Timeout, err))
			}
		}
	}

	// Validate ReliabilityConfig
	// 校验可靠性配置
	if config.Reliability.ConfigBackup.Enabled {
		if config.Reliability.ConfigBackup.Location == "" {
			return errors.New(errors.ErrTypeValidation, "dfx.reliability.configBackup.location is required when config backup is enabled")
		}
		// TODO: Add validation for schedule format and location type
	}
	if config.Reliability.ETCDBackup.Enabled {
		if config.Reliability.ETCDBackup.Location == "" {
			return errors.New(errors.ErrTypeValidation, "dfx.reliability.etcdBackup.location is required when etcd backup is enabled")
		}
		// TODO: Add validation for schedule format and location type
	}

	return nil
}
