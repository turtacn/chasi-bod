// Package types defines common data structures used across chasi-bod.
// 包 types 定义了 chasi-bod 中使用的通用数据结构。
package types

import (
	"github.com/turtacn/chasi-bod/common/types/enum"
)

// Node represents a single node in the cluster.
// Node 表示集群中的单个节点。
type Node struct {
	Name            string            `json:"name"`            // Node name / 节点名称
	Address         string            `json:"address"`         // Node IP address or hostname / 节点 IP 地址或主机名
	Roles           []enum.NodeRole   `json:"roles"`           // Node roles (e.g., master, worker) / 节点角色（例如，主节点、工作节点）
	Labels          map[string]string `json:"labels"`          // Kubernetes node labels / Kubernetes 节点标签
	Annotations     map[string]string `json:"annotations"`     // Kubernetes node annotations / Kubernetes 节点注解
	Taints          []string          `json:"taints"`          // Kubernetes node taints / Kubernetes 节点污点
	OperatingSystem string            `json:"operatingSystem"` // Operating System type and version / 操作系统类型和版本
	Architecture    string            `json:"architecture"`    // Node architecture (e.g., amd64) / 节点架构（例如，amd64）
}

// Cluster represents a collection of nodes forming a Kubernetes cluster.
// Cluster 表示构成 Kubernetes 集群的节点集合。
type Cluster struct {
	Name              string        `json:"name"`              // Cluster name / 集群名称
	APIServer         string        `json:"apiServer"`         // API server endpoint / API 服务器端点
	Nodes             []Node        `json:"nodes"`             // List of nodes in the cluster / 集群中的节点列表
	KubernetesVersion string        `json:"kubernetesVersion"` // Kubernetes version / Kubernetes 版本
	ContainerRuntime  string        `json:"containerRuntime"`  // Container runtime used / 使用的容器运行时
	NetworkConfig     NetworkConfig `json:"networkConfig"`     // Cluster network configuration / 集群网络配置
	StorageConfig     StorageConfig `json:"storageConfig"`     // Cluster storage configuration / 集群存储配置
}

// NetworkConfig represents the network configuration for a cluster or vcluster.
// NetworkConfig 表示集群或 vcluster 的网络配置。
type NetworkConfig struct {
	Plugin       string            `json:"plugin"`       // CNI plugin name / CNI 插件名称
	PodCIDR      string            `json:"podCIDR"`      // Pod network CIDR / Pod 网络 CIDR
	ServiceCIDR  string            `json:"serviceCIDR"`  // Service network CIDR / Service 网络 CIDR
	DNSServiceIP string            `json:"dnsServiceIP"` // DNS Service IP / DNS Service IP
	Interfaces   []InterfaceConfig `json:"interfaces"`   // Network interface configurations / 网络接口配置
	// Add more network specific configurations as needed
	// 根据需要添加更多网络特定配置
}

// InterfaceConfig represents the configuration of a network interface.
// InterfaceConfig 表示网络接口的配置。
type InterfaceConfig struct {
	Name    string   `json:"name"`    // Interface name (e.g., eth0) / 接口名称（例如，eth0）
	IPAddrs []string `json:"ipAddrs"` // IP addresses (e.g., ["192.168.1.10/24"]) / IP 地址（例如，["192.168.1.10/24"]）
	Gateway string   `json:"gateway"` // Default gateway / 默认网关
	// Add more interface specific configurations
	// 添加更多接口特定配置
}

// StorageConfig represents the storage configuration for a cluster or vcluster.
// StorageConfig 表示集群或 vcluster 的存储配置。
type StorageConfig struct {
	DefaultStorageClass string               `json:"defaultStorageClass"` // Default StorageClass name / 默认 StorageClass 名称
	StorageClasses      []StorageClassConfig `json:"storageClasses"`      // List of StorageClass configurations / StorageClass 配置列表
	PVConfigs           []PVConfig           `json:"pvConfigs"`           // List of predefined PV configurations (if any) / 预定义的 PV 配置列表（如果有）
	// Add more storage specific configurations as needed
	// 根据需要添加更多存储特定配置
}

// StorageClassConfig represents the configuration of a StorageClass.
// StorageClassConfig 表示 StorageClass 的配置。
type StorageClassConfig struct {
	Name        string            `json:"name"`        // StorageClass name / StorageClass 名称
	Provisioner string            `json:"provisioner"` // Provisioner name / Provisioner 名称
	Parameters  map[string]string `json:"parameters"`  // Provisioner parameters / Provisioner 参数
	// Add more StorageClass specific configurations
	// 添加更多 StorageClass 特定配置
}

// PVConfig represents the configuration of a PersistentVolume.
// PVConfig 表示 PersistentVolume 的配置。
type PVConfig struct {
	Name                   string                 `json:"name"`                   // PV name / PV 名称
	Capacity               string                 `json:"capacity"`               // PV capacity (e.g., "10Gi") / PV 容量（例如，“10Gi”）
	AccessModes            []string               `json:"accessModes"`            // Access modes (e.g., ["ReadWriteOnce"]) / 访问模式（例如，["ReadWriteOnce"]）
	PersistentVolumeSource PersistentVolumeSource `json:"persistentVolumeSource"` // PV source configuration / PV 源配置
	// Add other PV source types as needed (NFS, iSCSI, etc.)
	// 根据需要添加其他 PV 源类型（NFS、iSCSI 等）
}

// PersistentVolumeSource represents the source configuration for a PersistentVolume.
// PersistentVolumeSource 表示 PersistentVolume 的源配置。
type PersistentVolumeSource struct {
	HostPath *HostPathVolumeSource `json:"hostPath,omitempty"` // HostPath configuration / HostPath 配置
	// Add other PV source types as needed (NFS, iSCSI, etc.)
	// 根据需要添加其他 PV 源类型（NFS、iSCSI 等）
}

// HostPathVolumeSource represents a HostPath volume source.
// HostPathVolumeSource 表示 HostPath 卷源。
type HostPathVolumeSource struct {
	Path string `json:"path"` // Host path / 主机路径
	Type string `json:"type"` // Host path type / 主机路径类型
}

// SysctlConfig represents kernel parameter configuration.
// SysctlConfig 表示内核参数配置。
type SysctlConfig map[string]string // Map of sysctl parameters and their desired values / sysctl 参数及其期望值的映射

// VClusterConfig represents the configuration for a virtual cluster.
// VClusterConfig 表示虚拟集群的配置。
type VClusterConfig struct {
	Name              string            `json:"name"`              // VCluster name / VCluster 名称
	Namespace         string            `json:"namespace"`         // Host Cluster namespace where vcluster runs / vcluster 在 host 集群中运行的命名空间
	Template          string            `json:"template"`          // Name of the vcluster template to use / 使用的 vcluster 模板名称
	KubernetesVersion string            `json:"kubernetesVersion"` // Kubernetes version for the vcluster / vcluster 的 Kubernetes 版本
	ServiceCIDR       string            `json:"serviceCIDR"`       // Service network CIDR for the vcluster / vcluster 的 Service 网络 CIDR
	PodCIDR           string            `json:"podCIDR"`           // Pod network CIDR for the vcluster / vcluster 的 Pod 网络 CIDR
	ResourceRequests  map[string]string `json:"resourceRequests"`  // Resource requests for vcluster control plane pods in Host Cluster / vcluster 控制面 Pod 在 Host 集群中的资源请求
	ResourceLimits    map[string]string `json:"resourceLimits"`    // Resource limits for vcluster control plane pods in Host Cluster / vcluster 控制面 Pod 在 Host 集群中的资源限制
	// Add more vcluster specific configurations based on loft-sh/vcluster options
	// 根据 loft-sh/vcluster 选项添加更多 vcluster 特定配置
}

// ApplicationConfig represents the configuration for a business application deployment.
// ApplicationConfig 表示业务应用部署的配置。
type ApplicationConfig struct {
	Name         string                 `json:"name"`         // Application name / 应用程序名称
	VClusterName string                 `json:"vclusterName"` // Name of the target vcluster / 目标 vcluster 名称
	Namespace    string                 `json:"namespace"`    // Namespace within the vcluster to deploy to / 部署到 vcluster 内的命名空间
	Type         enum.ApplicationType   `json:"type"`         // Application type (e.g., CPU-bound, IO-bound) / 应用程序类型（例如，CPU 密集型、IO 密集型）
	Manifests    string                 `json:"manifests"`    // Path or content of K8s manifests (Helm Chart, Kustomize, etc.) / K8s manifests 的路径或内容（Helm Chart、Kustomize 等）
	Values       map[string]interface{} `json:"values"`       // Values to override in manifests (e.g., Helm values) / 用于覆盖 manifests 中的值（例如，Helm 值）
	// Add more application specific configurations
	// 添加更多应用程序特定配置
}

// PlatformConfig represents the overall configuration for a chasi-bod platform instance.
// PlatformConfig 表示 chasi-bod 平台实例的总体配置。
type PlatformConfig struct {
	Cluster        ClusterConfig                `json:"cluster"`      // Host Cluster configuration / Host Cluster 配置
	VClusterConfig map[string]VClusterConfig    `json:"vclusters"`    // VCluster configurations by name / 按名称区分的 VCluster 配置
	Applications   map[string]ApplicationConfig `json:"applications"` // Application configurations by name / 按名称区分的应用程序配置
	SysctlConfig   SysctlConfig                 `json:"sysctl"`       // Host OS sysctl configuration / Host OS sysctl 配置
	// Add other top-level configurations like DFX settings, etc.
	// 添加其他顶层配置，例如 DFX 设置等
}

// ClusterConfig represents the configuration for the Host Cluster.
// ClusterConfig 表示 Host Cluster 的配置。
type ClusterConfig struct {
	Name              string        `json:"name"`              // Host Cluster name / Host Cluster 名称
	KubernetesVersion string        `json:"kubernetesVersion"` // Desired Kubernetes version for Host Cluster / Host Cluster 期望的 Kubernetes 版本
	ContainerRuntime  string        `json:"containerRuntime"`  // Desired Container runtime for Host Cluster / Host Cluster 期望的容器运行时
	Network           NetworkConfig `json:"network"`           // Network configuration / 网络配置
	Storage           StorageConfig `json:"storage"`           // Storage configuration / 存储配置
	Nodes             []NodeConfig  `json:"nodes"`             // Node configurations / 节点配置
	// Add other host cluster specific configurations
	// 添加其他 Host Cluster 特定配置
}

// NodeConfig represents the configuration for a node during deployment.
// NodeConfig 表示部署期间节点的配置。
type NodeConfig struct {
	Address      string            `json:"address"`              // Node IP address or hostname / 节点 IP 地址或主机名
	User         string            `json:"user"`                 // SSH user for deployment / 用于部署的 SSH 用户
	Password     string            `json:"password,omitempty"`   // SSH password (discouraged) / SSH 密码（不推荐）
	PrivateKey   string            `json:"privateKey,omitempty"` // SSH private key path / SSH 私钥路径
	Port         int               `json:"port"`                 // SSH port / SSH 端口
	Roles        []enum.NodeRole   `json:"roles"`                // Node roles / 节点角色
	Labels       map[string]string `json:"labels"`               // Desired Kubernetes node labels / 期望的 Kubernetes 节点标签
	Annotations  map[string]string `json:"annotations"`          // Desired Kubernetes node annotations / 期望的 Kubernetes 节点注解
	Taints       []string          `json:"taints"`               // Desired Kubernetes node taints / 期望的 Kubernetes 节点污点
	SysctlConfig SysctlConfig      `json:"sysctl"`               // Node-specific sysctl configuration / 节点特定的 sysctl 配置
	// Add more node specific configurations like disks, etc.
	// 添加更多节点特定配置，例如磁盘等
}
