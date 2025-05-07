// Package constants defines project-wide constant values.
// 包 constants 定义了项目范围内的常量值。
package constants

import (
//"time"
)

// DefaultProjectName defines the default name of the project.
// DefaultProjectName 定义了项目的默认名称。
const DefaultProjectName = "chasi-bod"

// DefaultConfigPath defines the default path for configuration files.
// DefaultConfigPath 定义了配置文件的默认路径。
const DefaultConfigPath = "/etc/chasi-bod/config.yaml"

// DefaultDataDir defines the default data directory for chasi-bod.
// DefaultDataDir 定义了 chasi-bod 的默认数据目录。
const DefaultDataDir = "/var/lib/chasi-bod"

// DefaultLogDir defines the default log directory.
// DefaultLogDir 定义了默认的日志目录。
const DefaultLogDir = "/var/log/chasi-bod"

// DefaultBuildOutputDir defines the default directory for build outputs.
// DefaultBuildOutputDir 定义了构建输出的默认目录。
const DefaultBuildOutputDir = "./output"

// DefaultTimeout specifies a default timeout duration for operations.
// DefaultTimeout 指定了操作的默认超时时间。
const DefaultTimeout = 5 * time.Minute

// APIVersion defines the API version used by chasi-bod.
// APIVersion 定义了 chasi-bod 使用的 API 版本。
const APIVersion = "v1alpha1"

// DefaultVClusterNamespacePrefix defines the default prefix for vcluster namespaces in the host cluster.
// DefaultVClusterNamespacePrefix 定义了 host 集群中 vcluster 命名空间的默认前缀。
const DefaultVClusterNamespacePrefix = "vcluster-"

// DefaultAPIPort defines the default port for the chasi-bod API server (if implemented).
// DefaultAPIPort 定义了 chasi-bod API 服务器的默认端口（如果实现的话）。
const DefaultAPIPort = 8080

// DefaultKubeAPIServerPort defines the default port for the Kubernetes API server.
// DefaultKubeAPIServerPort 定义了 Kubernetes API 服务器的默认端口。
const DefaultKubeAPIServerPort = 6443

// DefaultSSHPort defines the default SSH port.
// DefaultSSHPort 定义了默认的 SSH 端口。
const DefaultSSHPort = 22

// DefaultContainerRuntimeEndpoint defines the default endpoint for the container runtime.
// DefaultContainerRuntimeEndpoint 定义了容器运行时的默认端点。
const DefaultContainerRuntimeEndpoint = "unix:///var/run/containerd/containerd.sock"

// ExitCodeSuccess represents a successful exit code.
// ExitCodeSuccess 表示成功的退出码。
const ExitCodeSuccess = 0

// ExitCodeFailure represents a generic failure exit code.
// ExitCodeFailure 表示通用的失败退出码。
const ExitCodeFailure = 1
