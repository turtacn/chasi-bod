// Package enum defines project-wide enumerated types.
// 包 enum 定义了项目范围内的枚举类型。
package enum

// NodeRole represents the role of a node in a Kubernetes cluster.
// NodeRole 表示节点在 Kubernetes 集群中的角色。
type NodeRole string

const (
	// RoleMaster indicates a control plane node.
	// RoleMaster 表示控制平面节点。
	RoleMaster NodeRole = "master"
	// RoleWorker indicates a worker node.
	// RoleWorker 表示工作节点。
	RoleWorker NodeRole = "worker"
	// RoleEdge indicates an edge node (less common in this context, but possible).
	// RoleEdge 表示边缘节点（在此上下文中不太常见，但可能）。
	RoleEdge NodeRole = "edge"
)

// ApplicationType represents the category of a business application based on resource usage.
// ApplicationType 表示基于资源使用情况的业务应用程序类别。
type ApplicationType string

const (
	// AppTypeComputeBound indicates a CPU-intensive application.
	// AppTypeComputeBound 表示 CPU 密集型应用程序。
	AppTypeComputeBound ApplicationType = "compute-bound"
	// AppTypeIOBound indicates an I/O-intensive application.
	// AppTypeIOBound 表示 I/O 密集型应用程序。
	AppTypeIOBound ApplicationType = "io-bound"
	// AppTypeMemoryBound indicates a memory-intensive application.
	// AppTypeMemoryBound 表示内存密集型应用程序。
	AppTypeMemoryBound ApplicationType = "memory-bound"
	// AppTypeNetworkBound indicates a network-intensive application.
	// AppTypeNetworkBound 表示网络密集型应用程序。
	AppTypeNetworkBound ApplicationType = "network-bound"
	// AppTypeGeneral indicates a general-purpose application.
	// AppTypeGeneral 表示通用应用程序。
	AppTypeGeneral ApplicationType = "general"
	// Add combinations or more granular types as needed
	// 根据需要添加组合或更细粒度的类型
)

// BuilderOutputFormat represents the output format of the platform image builder.
// BuilderOutputFormat 表示平台镜像构建器的输出格式。
type BuilderOutputFormat string

const (
	// OutputFormatISO indicates an ISO image.
	// OutputFormatISO 表示 ISO 镜像。
	OutputFormatISO BuilderOutputFormat = "iso"
	// OutputFormatVMA indicates a VMA (VMware vCenter) image.
	// OutputFormatVMA 表示 VMA (VMware vCenter) 镜像。
	OutputFormatVMA BuilderOutputFormat = "vma"
	// OutputFormatOVA indicates an OVA (Open Virtual Appliance) image.
	// OutputFormatOVA 表示 OVA (开放虚拟设备) 镜像。
	OutputFormatOVA BuilderOutputFormat = "ova"
	// OutputFormatQCOW2 indicates a QCOW2 (QEMU Copy On Write) image.
	// OutputFormatQCOW2 表示 QCOW2 (QEMU 写时复制) 镜像。
	OutputFormatQCOW2 BuilderOutputFormat = "qcow2"
	// OutputFormatAMI indicates an AMI (Amazon Machine Image).
	// OutputFormatAMI BuilderOutputFormat = "ami" // Uncomment if supporting AWS
	// Add other formats as needed
	// 根据需要添加其他格式
)
