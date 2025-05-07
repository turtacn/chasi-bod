// Package phases defines the individual steps involved in deploying the chasi-bod platform.
// 包 phases 定义了部署 chasi-bod 平台的各个步骤。
package phases

import (
	"context"
	"fmt"
	"time"

	//"time"

	"github.com/turtacn/chasi-bod/common/errors"
	"github.com/turtacn/chasi-bod/common/utils" // Assuming logger is here // 假设日志记录器在这里
	"github.com/turtacn/chasi-bod/pkg/config/model"
	//vcluster_client "github.com/turtacn/chasi-bod/pkg/vcluster/client" // Alias to avoid naming conflict // 别名以避免命名冲突
	vcluster_mgr "github.com/turtacn/chasi-bod/pkg/vcluster" // Alias for vcluster manager // vcluster manager 的别名
	// Assuming you have vcluster management logic and K8s client
	// 假设您有 vcluster 管理逻辑和 K8s 客户端
	"k8s.io/client-go/kubernetes"
)

// VClusterDeployPhase defines the interface for deploying virtual Kubernetes clusters.
// VClusterDeployPhase 定义了部署虚拟 Kubernetes 集群的接口。
// This phase assumes the Host Kubernetes cluster is already installed and ready.
// 此阶段假设 Host Kubernetes 集群已安装并就绪。
type VClusterDeployPhase interface {
	// Run executes the vcluster deployment phase.
	// Run 执行 vcluster 部署阶段。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// config: The overall platform configuration, including vcluster definitions. / 整体平台配置，包括 vcluster 定义。
	// hostK8sClient: A Kubernetes client connected to the Host Cluster. / 连接到 Host 集群的 Kubernetes 客户端。
	// Returns an error if the phase fails.
	// 如果阶段失败则返回错误。
	Run(ctx context.Context, config *model.PlatformConfig, hostK8sClient interface{} /* kubernetes.Interface */) error
}

// NewVClusterDeployPhase creates a new VClusterDeployPhase instance.
// NewVClusterDeployPhase 创建一个新的 VClusterDeployPhase 实例。
// Returns a VClusterDeployPhase implementation.
// 返回 VClusterDeployPhase 实现。
func NewVClusterDeployPhase() VClusterDeployPhase {
	return &defaultVClusterDeployPhase{}
}

// defaultVClusterDeployPhase is a default implementation of the VClusterDeployPhase.
// defaultVClusterDeployPhase 是 VClusterDeployPhase 的默认实现。
type defaultVClusterDeployPhase struct {
	// vclusterManager vcluster_mgr.Manager // Vcluster manager to perform operations // 执行操作的 Vcluster 管理器
}

// Run executes the vcluster deployment phase.
// Run 执行 vcluster 部署阶段。
func (p *defaultVClusterDeployPhase) Run(ctx context.Context, config *model.PlatformConfig, hostK8sClient interface{} /* kubernetes.Interface */) error {
	utils.GetLogger().Println("Running VClusterDeployPhase...")

	// Ensure hostK8sClient is a valid Kubernetes client interface
	// 确保 hostK8sClient 是一个有效的 Kubernetes 客户端接口
	k8sClient, ok := hostK8sClient.(kubernetes.Interface)
	if !ok {
		return errors.New(errors.ErrTypeInternal, "invalid Host Kubernetes client provided for vcluster deployment")
	}

	// Initialize vclusterManager with the hostK8sClient
	// 使用 hostK8sClient 初始化 vclusterManager
	vclusterManager := vcluster_mgr.NewManager(k8sClient) // Assuming NewManager takes a K8s client
	// p.vclusterManager = vclusterManager // If storing in struct

	if len(config.VClusters) == 0 {
		utils.GetLogger().Println("No vclusters defined in configuration, skipping vcluster deployment.")
		return nil
	}

	// Step 1: Deploy each vcluster defined in the configuration
	// 步骤 1：部署配置中定义的每个 vcluster
	utils.GetLogger().Printf("Deploying %d vclusters...", len(config.VClusters))
	for name, vclusterCfg := range config.VClusters {
		vclusterCtx, cancel := context.WithTimeout(ctx, 10*time.Minute) // Context for vcluster deployment
		defer cancel()

		utils.GetLogger().Printf("Deploying vcluster '%s' in host namespace '%s'...", name, vclusterCfg.Namespace)
		// Use vclusterManager to create the vcluster
		// 使用 vclusterManager 创建 vcluster
		err := vclusterManager.Create(vclusterCtx, &vclusterCfg) // Assuming Create method exists
		if err != nil {
			return errors.NewWithCause(errors.ErrTypeVCluster, fmt.Sprintf("failed to deploy vcluster '%s'", name), err)
		}
		utils.GetLogger().Printf("VCluster '%s' deployment initiated.", name)

		// Step 2: Wait for the vcluster to be ready (API server accessible)
		// 步骤 2：等待 vcluster 就绪（API 服务器可访问）
		utils.GetLogger().Printf("Waiting for vcluster '%s' to be ready...", name)
		// Use vclusterManager to wait for vcluster API server to be accessible
		// 使用 vclusterManager 等待 vcluster API 服务器可访问
		err = vclusterManager.WaitForReady(vclusterCtx, name) // Assuming WaitForReady method exists
		if err != nil {
			return errors.NewWithCause(errors.ErrTypeVCluster, fmt.Sprintf("vcluster '%s' did not become ready", name), err)
		}
		utils.GetLogger().Printf("VCluster '%s' is ready.", name)

		// TODO: Optionally, apply base configurations/templates inside the vcluster
		// This might involve getting a client for the *virtual* cluster and applying manifests
		// 这可能涉及获取虚拟集群的客户端并应用 manifests
		utils.GetLogger().Printf("Applying base configurations inside vcluster '%s'...", name)
		// vclusterClient, err := vclusterManager.GetVClusterClient(vclusterCtx, name) // Assuming GetVClusterClient exists
		// if err != nil {
		// 	utils.GetLogger().Printf("Warning: Failed to get client for vcluster '%s' to apply base configs: %v", name, err)
		// 	// Decide if this is a fatal error or just a warning
		// } else {
		// 	err = applyBaseConfigsInVCluster(vclusterCtx, vclusterClient, &vclusterCfg) // Needs implementation
		// 	if err != nil {
		// 		utils.GetLogger().Printf("Warning: Failed to apply base configs inside vcluster '%s': %v", name, err)
		// 		// Decide if this is a fatal error or just a warning
		// 	} else {
		// 		utils.GetLogger().Printf("Placeholder: Base configurations applied inside vcluster '%s'.", name)
		// 	}
		// }
		utils.GetLogger().Printf("Placeholder: Base configurations application skipped for vcluster '%s'.", name)

	}

	utils.GetLogger().Println("VClusterDeployPhase completed successfully.")
	return nil
}

// TODO: Implement a function to apply base configurations inside a virtual cluster
// TODO: 实现一个在虚拟集群内部应用基础配置的函数
// This would likely involve using the vclusterClient to create/update resources like ConfigMaps, Deployments, etc.
// 这很可能涉及使用 vclusterClient 创建/更新 ConfigMaps, Deployments 等资源。
// func applyBaseConfigsInVCluster(ctx context.Context, vclusterClient kubernetes.Interface, vclusterCfg *model.VClusterConfig) error {
// 	utils.GetLogger().Printf("Applying base configs in vcluster %s...", vclusterCfg.Name)
// 	// Example: Apply predefined base manifests (RBAC, default StorageClass if different from host etc.)
// 	// Load manifests from templates or embedded files
// 	// Apply manifests using vclusterClient.Apply() or Create()/Update()
// 	return errors.New(errors.ErrTypeNotImplemented, "applying base configs inside vcluster not implemented")
// }
