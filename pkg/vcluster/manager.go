// Package vcluster provides functionality for managing virtual Kubernetes clusters using loft-sh/vcluster.
// 包 vcluster 提供了使用 loft-sh/vcluster 管理虚拟 Kubernetes 集群的功能。
package vcluster

import (
	"context"
	"fmt"
	"time"

	"github.com/turtacn/chasi-bod/common/errors"
	"github.com/turtacn/chasi-bod/common/utils" // Assuming logger is here // 假设日志记录器在这里
	"github.com/turtacn/chasi-bod/pkg/config/model"
	vcluster_client "github.com/turtacn/chasi-bod/pkg/vcluster/client" // Alias to avoid naming conflict // 别名以避免命名冲突

	corev1 "k8s.io/api/core/v1"                    // Added for Namespace and Secret types // 添加用于 Namespace 和 Secret 类型
	apierrors "k8s.io/apimachinery/pkg/api/errors" // Added for checking Kubernetes API errors // 添加用于检查 Kubernetes API 错误
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"  // Added for meta types // 添加用于元数据类型
	"k8s.io/apimachinery/pkg/util/wait"            // Added for polling // 添加用于轮询
	"k8s.io/client-go/kubernetes"
	// You might need vcluster-specific client or helpers later
	// 稍后您可能需要 vcluster 特定的客户端或辅助工具
	// "github.com/loft-sh/vcluster/pkg/helm" // If using vcluster's internal helm deployment
	// "k8s.io/client-go/rest" // Might be needed to build client config // 可能需要它来构建客户端配置
)

// Manager defines the interface for managing vcluster instances in the Host Cluster.
// Manager 定义了在 Host 集群中管理 vcluster 实例的接口。
type Manager interface {
	// Create creates a new vcluster instance based on the configuration.
	// Create 根据配置创建一个新的 vcluster 实例。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// config: The configuration for the vcluster to create. / 要创建的 vcluster 配置。
	// Returns an error if creation fails.
	// 如果创建失败则返回错误。
	Create(ctx context.Context, config *model.VClusterConfig) error

	// Delete deletes a vcluster instance by name.
	// Delete 根据名称删除一个 vcluster 实例。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// name: The name of the vcluster to delete. / 要删除的 vcluster 的名称。
	// Returns an error if deletion fails.
	// 如果删除失败则返回错误。
	Delete(ctx context.Context, name string) error

	// WaitForReady waits for a vcluster instance to become ready (e.g., API server accessible).
	// WaitForReady 等待 vcluster 实例就绪（例如，API 服务器可访问）。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// name: The name of the vcluster to wait for. / 要等待的 vcluster 的名称。
	// Returns an error if the vcluster does not become ready within the timeout.
	// 如果 vcluster 在超时时间内未就绪则返回错误。
	WaitForReady(ctx context.Context, name string) error

	// GetVClusterClient returns a Kubernetes client configured to interact with the specified virtual cluster.
	// GetVClusterClient 返回配置用于与指定虚拟集群交互的 Kubernetes 客户端。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// name: The name of the vcluster. / vcluster 的名称。
	// Returns a Kubernetes client and an error if retrieval/configuration fails.
	// 返回 Kubernetes 客户端，以及检索/配置失败时的错误。
	GetVClusterClient(ctx context.Context, name string) (kubernetes.Interface, error)

	// Add other management methods like Update, List, GetVClusterStatus etc.
	// 添加其他管理方法，例如 Update, List, GetVClusterStatus 等。
	// Update(ctx context.Context, config *model.VClusterConfig) error
	// List(ctx context.Context) ([]*model.VClusterConfig, error) // Maybe return status as well
	// GetVClusterStatus(ctx context.Context, name string) (*VClusterStatus, error) // Define VClusterStatus struct
}

// defaultManager is a default implementation of the VCluster Manager.
// defaultManager 是 VCluster Manager 的默认实现。
type defaultManager struct {
	hostK8sClient kubernetes.Interface // Client to interact with the Host Cluster / 用于与 Host 集群交互的客户端
}

// NewManager creates a new VCluster Manager.
// NewManager 创建一个新的 VCluster Manager。
// hostK8sClient: A Kubernetes client connected to the Host Cluster. / 连接到 Host 集群的 Kubernetes 客户端。
// Returns a VCluster Manager implementation.
// 返回 VCluster Manager 实现。
func NewManager(hostK8sClient kubernetes.Interface) Manager {
	return &defaultManager{
		hostK8sClient: hostK8sClient,
	}
}

// Create creates a new vcluster instance.
// Create 创建一个新的 vcluster 实例。
func (m *defaultManager) Create(ctx context.Context, config *model.VClusterConfig) error {
	utils.GetLogger().Printf("Creating vcluster '%s' in host namespace '%s'", config.Name, config.Namespace)

	// Step 1: Ensure the host namespace exists
	// 步骤 1：确保 host 命名空间存在
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: config.Namespace,
		},
	}
	_, err := m.hostK8sClient.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return errors.NewWithCause(errors.ErrTypeVCluster, fmt.Sprintf("failed to create host namespace %s for vcluster %s", config.Namespace, config.Name), err)
	}
	utils.GetLogger().Printf("Host namespace '%s' for vcluster '%s' ensured.", config.Namespace, config.Name)

	// Step 2: Deploy the vcluster using Helm or vcluster manifests
	// 步骤 2：使用 Helm 或 vcluster manifests 部署 vcluster
	// loft-sh/vcluster is typically deployed via Helm or the vcluster CLI
	// This implementation should encapsulate that logic, often by applying K8s manifests.
	// loft-sh/vcluster 通常通过 Helm 或 vcluster CLI 部署
	// 此实现应封装该逻辑，通常通过应用 K8s manifests 来实现。
	utils.GetLogger().Printf("Deploying vcluster '%s' via manifests...", config.Name)

	// Need to load and process vcluster deployment manifests.
	// Requires finding the correct manifest based on config.KubernetesVersion etc.
	// 需要加载和处理 vcluster 部署 manifests。
	// 需要根据 config.KubernetesVersion 等找到正确的 manifest。
	// The manifests are usually K8s YAML files defining Deployment, Service, Secret etc.
	// manifest 通常是定义 Deployment, Service, Secret 等的 K8s YAML 文件。
	// They might be embedded in the chasi-bod binary or located in a known path.
	// 它们可能嵌入在 chasi-bod 二进制文件中或位于已知路径。

	// Example: Assuming manifests are loaded from a template or embedded resource
	// 示例：假设 manifests 从模板或嵌入资源加载
	// vclusterManifests, err := vcluster_template.LoadAndProcessVClusterManifests(config) // Needs implementation
	// if err != nil { return fmt.Errorf("failed to load vcluster manifests: %w", err) }

	// TODO: Implement logic to apply vcluster manifests to the Host Cluster in the specified namespace
	// TODO: 实现将 vcluster manifests 应用到 Host 集群指定命名空间的逻辑
	// This involves parsing the manifest YAML and using the hostK8sClient to create/update resources.
	// 这涉及解析 manifest YAML 并使用 hostK8sClient 创建/更新资源。
	// Can use client-go's dynamic client or a library like "k8s.io/kubectl/pkg/cmd/apply" (complex).
	// 可以使用 client-go 的 dynamic client 或像 "k8s.io/kubectl/pkg/cmd/apply" 这样的库（复杂）。
	// A simpler approach is to parse manually and call appropriate Create/Update methods on typed clients (corev1, appsv1).
	// 更简单的方法是手动解析并调用类型化客户端（corev1, appsv1）上的相应 Create/Update 方法。

	utils.GetLogger().Printf("Placeholder: Applying vcluster '%s' K8s manifests.", config.Name)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Second): // Simulate manifest application time
		utils.GetLogger().Printf("Placeholder: Vcluster '%s' K8s manifests applied.", config.Name)
	}

	// Step 3: Wait for the vcluster's API server pod/deployment to be ready in the host namespace
	// 步骤 3：等待 vcluster 的 API 服务器 Pod/deployment 在 host 命名空间中运行并就绪
	utils.GetLogger().Printf("Waiting for vcluster '%s' deployment to be ready in host namespace '%s'...", config.Name, config.Namespace)
	// You can wait for the main vcluster deployment (or statefulset) to have ready replicas.
	// The deployment name is usually the vcluster name.
	// 您可以等待主要的 vcluster deployment（或 statefulset）具有就绪副本。
	// deployment 名称通常与 vcluster 名称相同。
	vclusterDeploymentName := config.Name // Assuming deployment name matches vcluster name

	err = wait.PollUntilContextTimeout(ctx, 5*time.Second, 5*time.Minute, func(ctx context.Context) (bool, error) {
		deployment, getErr := m.hostK8sClient.AppsV1().Deployments(config.Namespace).Get(ctx, vclusterDeploymentName, metav1.GetOptions{})
		if getErr != nil {
			if apierrors.IsNotFound(getErr) {
				utils.GetLogger().Printf("Debug: Vcluster deployment '%s' not found yet in namespace %s, retrying...", vclusterDeploymentName, config.Namespace)
				return false, nil // Deployment not created yet
			}
			utils.GetLogger().Printf("Error getting vcluster deployment '%s' in namespace %s: %v", vclusterDeploymentName, config.Namespace, getErr)
			return false, getErr // Return error to stop polling
		}

		if deployment.Status.ReadyReplicas > 0 && deployment.Status.ReadyReplicas >= deployment.Status.Replicas {
			utils.GetLogger().Printf("VCluster deployment '%s' is ready (%d/%d replicas).", vclusterDeploymentName, deployment.Status.ReadyReplicas, deployment.Status.Replicas)
			return true, nil // Deployment is ready
		}

		utils.GetLogger().Printf("VCluster deployment '%s' not ready yet (%d/%d replicas) in namespace %s, retrying...",
			vclusterDeploymentName, deployment.Status.ReadyReplicas, deployment.Status.Replicas, config.Namespace)
		return false, nil // Deployment exists but not ready
	})

	if err != nil {
		return errors.NewWithCause(errors.ErrTypeVCluster, fmt.Sprintf("vcluster deployment '%s' timed out or failed to become ready in host namespace '%s'", vclusterDeploymentName, config.Namespace), err)
	}

	utils.GetLogger().Printf("VCluster '%s' created and ready in host namespace '%s'.", config.Name, config.Namespace)
	return nil
}

// Delete deletes a vcluster instance.
// Delete 删除一个 vcluster 实例。
func (m *defaultManager) Delete(ctx context.Context, name string) error {
	utils.GetLogger().Printf("Deleting vcluster '%s'...", name)

	// Determine the vcluster's host namespace
	// 确定 vcluster 的 host 命名空间
	// This should ideally come from a stored state or the config if available,
	// but convention (e.g., vcluster-<name>) is common.
	// 这理想情况下应来自存储的状态或配置（如果可用），
	// 但约定（例如，vcluster-<name>）是常见的。
	namespace := fmt.Sprintf("vcluster-%s", name) // TODO: Get namespace reliably from config or stored state

	// Step 1: Delete the namespace in the host cluster
	// 步骤 1：删除 host 集群中的命名空间
	utils.GetLogger().Printf("Deleting host namespace '%s' for vcluster '%s'...", namespace, name)
	err := m.hostK8sClient.CoreV1().Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return errors.NewWithCause(errors.ErrTypeVCluster, fmt.Sprintf("failed to delete host namespace %s for vcluster %s", namespace, name), err)
	}
	utils.GetLogger().Printf("Host namespace '%s' for vcluster '%s' deletion requested.", namespace, name)

	// Step 2: Wait for the namespace to be fully deleted
	// 步骤 2：等待命名空间完全删除
	utils.GetLogger().Printf("Waiting for host namespace '%s' to be deleted...", namespace)
	err = wait.PollUntilContextTimeout(ctx, 5*time.Second, 5*time.Minute, func(ctx context.Context) (bool, error) {
		_, getErr := m.hostK8sClient.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
		if getErr != nil && apierrors.IsNotFound(getErr) {
			return true, nil // Namespace is gone
		}
		if getErr != nil {
			utils.GetLogger().Printf("Error checking host namespace status: %v", getErr)
			return false, getErr // Other error, stop polling
		}
		utils.GetLogger().Printf("Host namespace '%s' still exists...", namespace)
		return false, nil // Namespace still exists
	})
	if err != nil {
		return errors.NewWithCause(errors.ErrTypeVCluster, fmt.Sprintf("host namespace '%s' for vcluster '%s' did not fully delete within timeout", namespace, name), err)
	}
	utils.GetLogger().Printf("Host namespace '%s' deleted.", namespace)

	utils.GetLogger().Printf("VCluster '%s' deleted successfully.", name)
	return nil
}

// WaitForReady waits for a vcluster instance to become ready (API server accessible).
// WaitForReady 等待 vcluster 实例就绪（API 服务器可访问）。
func (m *defaultManager) WaitForReady(ctx context.Context, name string) error {
	utils.GetLogger().Printf("Waiting for vcluster '%s' to be accessible...", name)

	// This requires obtaining a client for the *virtual* cluster and checking its API server's healthz endpoint.
	// The `vcluster_client.GetVClusterClient` function handles the connection logic (e.g., port-forwarding or service lookup).
	// 这需要获取虚拟集群的客户端，并检查其 API 服务器的 healthz 端点。
	// `vcluster_client.GetVClusterClient` 函数处理连接逻辑（例如，端口转发或服务查找）。
	err := wait.PollUntilContextTimeout(ctx, 5*time.Second, 5*time.Minute, func(ctx context.Context) (bool, error) {
		vClient, getClientErr := vcluster_client.GetVClusterClient(ctx, name, m.hostK8sClient) // Use the helper function
		if getClientErr != nil {
			// If we can't even get a client, the vcluster is likely not ready or accessible yet
			// 如果我们甚至无法获取客户端，vcluster 可能尚未就绪或不可访问
			utils.GetLogger().Printf("Debug: Failed to get client for vcluster '%s', retrying: %v", name, getClientErr)
			// Do not return error here unless it's a non-recoverable configuration error
			return false, nil
		}

		// Check a basic API endpoint like /healthz or listing namespaces in the virtual cluster
		// 检查虚拟集群中的基本 API 端点，例如 /healthz 或列出命名空间
		// Listing namespaces is a simple way to test if the API server is responding and auth is working.
		// 列出命名空间是测试 API 服务器是否响应和身份验证是否正常工作的简单方法。
		_, healthzErr := vClient.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if healthzErr == nil {
			utils.GetLogger().Printf("VCluster '%s' API server is accessible and healthy.", name)
			return true, nil // Success!
		}

		utils.GetLogger().Printf("Debug: VCluster '%s' API server check failed, retrying: %v", name, healthzErr)
		return false, nil // Not healthy yet
	})

	if err != nil {
		return errors.NewWithCause(errors.ErrTypeVCluster, fmt.Sprintf("vcluster '%s' API server did not become accessible within timeout", name), err)
	}

	utils.GetLogger().Printf("VCluster '%s' is accessible and ready.", name)
	return nil
}

// GetVClusterClient returns a Kubernetes client for the specified virtual cluster.
// It delegates the actual client creation logic to the vcluster_client package.
// GetVClusterClient 返回指定虚拟集群的 Kubernetes 客户端。
// 它将实际的客户端创建逻辑委托给 vcluster_client 包。
func (m *defaultManager) GetVClusterClient(ctx context.Context, name string) (kubernetes.Interface, error) {
	// Delegate to the client package's GetVClusterClient function
	// 委托给 client 包的 GetVClusterClient 函数
	return vcluster_client.GetVClusterClient(ctx, name, m.hostK8sClient) // Assuming GetVClusterClient exists in pkg/vcluster/client
}

// Placeholder for VClusterStatus struct if needed
// VClusterStatus 结构体的占位符（如果需要）
// type VClusterStatus struct {
// 	Name      string `json:"name"`
// 	Namespace string `json:"namespace"` // Host namespace
// 	Ready     bool   `json:"ready"`
// 	Message   string `json:"message"`
// 	// Add more status fields from vcluster
// 	// 添加更多来自 vcluster 的状态字段
// }
