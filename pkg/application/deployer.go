// Package application provides functionality for deploying and managing business applications within virtual clusters.
// 包 application 提供了在虚拟集群中部署和管理业务应用程序的功能。
// It interacts with a specific virtual cluster to deploy the application manifests.
// 它与特定的虚拟集群交互以部署应用程序 manifests。
package application

import (
	"context"
	"fmt"
	"github.com/turtacn/chasi-bod/common/errors"
	"github.com/turtacn/chasi-bod/common/utils" // Assuming logger is here // 假设日志记录器在这里
	"github.com/turtacn/chasi-bod/pkg/config/model"
	vcluster_client "github.com/turtacn/chasi-bod/pkg/vcluster/client" // Alias to avoid naming conflict // 别名以避免命名冲突
	"k8s.io/client-go/kubernetes"
	// Assuming you have specific deployment method packages
	// 假设您有特定的部署方法包
	"github.com/turtacn/chasi-bod/pkg/application/helm"
	// "github.com/turtacn/chasi-bod/pkg/application/kustomize" // If implementing Kustomize deployment // 如果实现 Kustomize 部署
	// "github.com/turtacn/chasi-bod/pkg/application/manifest" // If implementing raw manifest application // 如果实现原始 manifest 应用
	// corev1 "k8s.io/api/core/v1" // Might be needed for config injection // 可能需要用于配置注入
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1" // Might be needed for config injection // 可能需要用于配置注入
	// "k8s.io/apimachinery/pkg/util/wait" // Needed for rollout status check // 需要用于 rollout 状态检查
)

// Deployer defines the interface for deploying and managing a business application.
// Deployer 定义了部署和管理业务应用程序的接口。
// It interacts with a specific virtual cluster to deploy the application manifests.
// 它与特定的虚拟集群交互以部署应用程序 manifests。
type Deployer interface {
	// Deploy deploys the application to the target vcluster.
	// Deploy 将应用程序部署到目标 vcluster。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// config: The configuration for the application deployment. / 应用程序部署的配置。
	// hostK8sClient: The client for the Host Kubernetes Cluster, potentially needed to get vcluster client.
	// hostK8sClient: 用于 Host Kubernetes 集群的客户端，可能需要它来获取 vcluster 客户端。
	// Returns an error if deployment fails.
	// 如果部署失败则返回错误。
	Deploy(ctx context.Context, config *model.ApplicationConfig, hostK8sClient kubernetes.Interface) error

	// Upgrade upgrades the deployed application.
	// Upgrade 升级已部署的应用程序。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// config: The configuration for the application upgrade. / 应用程序升级的配置。
	// hostK8sClient: The client for the Host Kubernetes Cluster. / 用于 Host Kubernetes 集群的客户端。
	// Returns an error if upgrade fails.
	// 如果升级失败则返回错误。
	Upgrade(ctx context.Context, config *model.ApplicationConfig, hostK8sClient kubernetes.Interface) error

	// Delete deletes the deployed application.
	// Delete 删除已部署的应用程序。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// appName: The name of the application. / 应用程序的名称。
	// vclusterName: The name of the vcluster where the app is deployed. / 部署应用程序的 vcluster 名称。
	// namespace: The namespace within the vcluster where the app is deployed. / 部署应用程序的 vcluster 内的命名空间。
	// hostK8sClient: The client for the Host Kubernetes Cluster. / 用于 Host Kubernetes 集群的客户端。
	// Returns an error if deletion fails.
	// 如果删除失败则返回错误。
	Delete(ctx context.Context, appName string, vclusterName string, namespace string, hostK8sClient kubernetes.Interface) error

	// GetStatus gets the status of the deployed application.
	// GetStatus 获取已部署应用程序的状态。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// appName: The name of the application. / 应用程序的名称。
	// vclusterName: The name of the vcluster where the app is deployed. / 部署应用程序的 vcluster 名称。
	// namespace: The namespace within the vcluster where the app is deployed. / 部署应用程序的 vcluster 内的命名空间。
	// hostK8sClient: The client for the Host Kubernetes Cluster. / 用于 Host Kubernetes 集群的客户端。
	// Returns the application status and an error.
	// 返回应用程序状态和错误。
	// GetStatus(ctx context.Context, appName string, vclusterName string, namespace string, hostK8sClient kubernetes.Interface) (*ApplicationStatus, error) // Define ApplicationStatus struct // 定义 ApplicationStatus 结构体
}

// defaultDeployer is a default implementation of the Application Deployer.
// defaultDeployer 是应用程序 Deployer 的默认实现。
type defaultDeployer struct {
	// Dependencies for different deployment methods
	// 不同部署方法的依赖项
	helmDeployer *helm.HelmDeployer // Assuming HelmDeployer struct exists and is used // 假设 HelmDeployer 结构体存在并被使用
	// kustomizeDeployer kustomize.Deployer // If implementing Kustomize deployment // 如果实现 Kustomize 部署
	// manifestDeployer manifest.Deployer // If implementing raw manifest application // 如果实现原始 manifest 应用
}

// NewDeployer creates a new Application Deployer.
// NewDeployer 创建一个新的应用程序 Deployer。
// Returns an Application Deployer implementation.
// 返回应用程序 Deployer 实现。
func NewDeployer() Deployer {
	return &defaultDeployer{
		helmDeployer: helm.NewHelmDeployer(), // Assuming NewHelmDeployer exists // 假设 NewHelmDeployer 存在
		// Initialize other deployers here
		// 在这里初始化其他 deployer
		// kustomizeDeployer: kustomize.NewDeployer(),
		// manifestDeployer: manifest.NewDeployer(),
	}
}

// Deploy deploys the application.
// Deploy 部署应用程序。
func (d *defaultDeployer) Deploy(ctx context.Context, config *model.ApplicationConfig, hostK8sClient kubernetes.Interface) error {
	utils.GetLogger().Printf("Deploying application '%s' to vcluster '%s' in namespace '%s'",
		config.Name, config.VClusterName, config.Namespace)

	// Step 1: Get the client for the target vcluster
	// 步骤 1：获取目标 vcluster 的客户端
	utils.GetLogger().Printf("Getting client for vcluster '%s'...", config.VClusterName)
	vclusterClient, err := vcluster_client.GetVClusterClient(ctx, config.VClusterName, hostK8sClient)
	if err != nil {
		return errors.NewWithCause(errors.ErrTypeVCluster, fmt.Sprintf("failed to get client for target vcluster '%s'", config.VClusterName), err)
	}
	utils.GetLogger().Printf("Successfully obtained client for vcluster '%s'.", config.VClusterName)

	// Step 2: Determine the deployment method and execute it
	// 步骤 2：确定部署方法并执行
	// Only one deployment method should be specified in the config (as per validation)
	// 配置中应只指定一种部署方法（根据校验规则）
	if config.HelmChart != nil {
		utils.GetLogger().Printf("Deploying application '%s' using Helm chart '%s'...", config.Name, config.HelmChart.Chart)
		// Helm deployment needs the vcluster client and the HelmChartConfig
		// Helm 部署需要 vcluster 客户端和 HelmChartConfig
		_, err := d.helmDeployer.DeployChart(ctx, vclusterClient, config.HelmChart, config.Namespace) // Assuming DeployChart method exists and returns release
		if err != nil {
			return errors.NewWithCause(errors.ErrTypeApplication, fmt.Sprintf("failed to deploy application '%s' using Helm", config.Name), err)
		}
		utils.GetLogger().Printf("Application '%s' deployed successfully using Helm.", config.Name)

	} else if config.Kustomize != nil {
		utils.GetLogger().Printf("Placeholder: Deploying application '%s' using Kustomize path '%s'...", config.Name, config.Kustomize.Path)
		// TODO: Implement Kustomize deployment logic using vclusterClient
		// The Kustomize deployer would likely build the manifests and then apply them using the K8s client.
		// Kustomize deployer 可能会构建 manifests，然后使用 K8s 客户端应用它们。
		// err := d.kustomizeDeployer.Deploy(ctx, vclusterClient, config.Kustomize, config.Namespace)
		return errors.New(errors.ErrTypeNotImplemented, fmt.Sprintf("Kustomize deployment for application '%s' not implemented yet", config.Name))

	} else if len(config.Manifests) > 0 {
		utils.GetLogger().Printf("Placeholder: Deploying application '%s' using raw manifests...", config.Name)
		// TODO: Implement raw manifest application logic using vclusterClient
		// The manifest deployer would read the manifest files and apply them using the K8s client.
		// manifest deployer 会读取 manifest 文件并使用 K8s 客户端应用它们。
		// err := d.manifestDeployer.Deploy(ctx, vclusterClient, config.Manifests, config.Namespace)
		return errors.New(errors.ErrTypeNotImplemented, fmt.Sprintf("raw manifest deployment for application '%s' not implemented yet", config.Name))

	} else {
		// This case should ideally be caught by configuration validation
		// 这个情况理想情况下应该被配置校验捕获
		return errors.New(errors.ErrTypeInternal, fmt.Sprintf("no deployment method specified for application '%s'", config.Name))
	}

	// Step 3: Optional - Inject configuration (e.g., update ConfigMaps/Secrets in vcluster)
	// 步骤 3：可选 - 注入配置（例如，更新 vcluster 中的 ConfigMaps/Secrets）
	if len(config.ConfigInjection) > 0 {
		utils.GetLogger().Printf("Injecting configuration for application '%s'...", config.Name)
		// This might involve:
		// - Reading sensitive data from a central secret store (e.g., Vault)
		// - Creating or updating ConfigMaps/Secrets in the target vcluster namespace using vclusterClient
		// - Restarting pods if needed to pick up new config
		// 这可能涉及：
		// - 从中央 Secret 存储（例如，Vault）读取敏感数据
		// - 在目标 vcluster 命名空间中使用 vclusterClient 创建或更新 ConfigMaps/Secrets
		// - 如果需要，重启 pods 以获取新配置
		// TODO: Implement config injection logic using vclusterClient
		utils.GetLogger().Printf("Placeholder: Configuration injected for application '%s'.", config.Name)
	}

	// Step 4: Optional - Wait for application rollout to complete
	// 步骤 4：可选 - 等待应用程序 rollout 完成
	utils.GetLogger().Printf("Waiting for application '%s' rollout to complete...", config.Name)
	// This involves checking the status of Deployments, StatefulSets etc. in the vcluster using vclusterClient
	// 这涉及使用 vclusterClient 检查 vcluster 中 Deployments, StatefulSets 等的状态
	// TODO: Implement rollout status check using vclusterClient
	// Example:
	// err = waitForDeploymentRollout(ctx, vclusterClient, config.Namespace, config.Name, 10*time.Minute) // Assuming a helper exists
	// if err != nil { return errors.NewWithCause(...) }
	utils.GetLogger().Printf("Placeholder: Application '%s' rollout completed.", config.Name)

	utils.GetLogger().Printf("Application Deployer completed successfully for '%s'.", config.Name)
	return nil
}

// Upgrade upgrades the deployed application.
// Upgrade 升级已部署的应用程序。
func (d *defaultDeployer) Upgrade(ctx context.Context, config *model.ApplicationConfig, hostK8sClient kubernetes.Interface) error {
	utils.GetLogger().Printf("Upgrading application '%s' to vcluster '%s' in namespace '%s'",
		config.Name, config.VClusterName, config.Namespace)

	// Step 1: Get the client for the target vcluster
	// 步骤 1：获取目标 vcluster 的客户端
	utils.GetLogger().Printf("Getting client for vcluster '%s'...", config.VClusterName)
	vclusterClient, err := vcluster_client.GetVClusterClient(ctx, config.VClusterName, hostK8sClient)
	if err != nil {
		return errors.NewWithCause(errors.ErrTypeVCluster, fmt.Sprintf("failed to get client for target vcluster '%s' for upgrade", config.VClusterName), err)
	}
	utils.GetLogger().Printf("Successfully obtained client for vcluster '%s'.", config.VClusterName)

	// Step 2: Determine the deployment method and execute upgrade
	// 步骤 2：确定部署方法并执行升级
	if config.HelmChart != nil {
		utils.GetLogger().Printf("Upgrading application '%s' using Helm chart '%s'...", config.Name, config.HelmChart.Chart)
		// Helm upgrade needs the vcluster client and the HelmChartConfig
		// Helm upgrade 需要 vcluster 客户端和 HelmChartConfig
		_, err := d.helmDeployer.UpgradeChart(ctx, vclusterClient, config.HelmChart, config.Namespace) // Assuming UpgradeChart method exists
		if err != nil {
			return errors.NewWithCause(errors.ErrTypeApplication, fmt.Sprintf("failed to upgrade application '%s' using Helm", config.Name), err)
		}
		utils.GetLogger().Printf("Application '%s' upgraded successfully using Helm.", config.Name)

	} else if config.Kustomize != nil {
		utils.GetLogger().Printf("Placeholder: Upgrading application '%s' using Kustomize path '%s'...", config.Name, config.Kustomize.Path)
		// TODO: Implement Kustomize upgrade logic
		// err := d.kustomizeDeployer.Upgrade(ctx, vclusterClient, config.Kustomize, config.Namespace)
		return errors.New(errors.ErrTypeNotImplemented, fmt.Sprintf("Kustomize upgrade for application '%s' not implemented yet", config.Name))

	} else if len(config.Manifests) > 0 {
		utils.GetLogger().Printf("Placeholder: Upgrading application '%s' using raw manifests...", config.Name)
		// TODO: Implement raw manifest upgrade logic (e.g., apply with --overwrite)
		// err := d.manifestDeployer.Upgrade(ctx, vclusterClient, config.Manifests, config.Namespace)
		return errors.New(errors.ErrTypeNotImplemented, fmt.Sprintf("raw manifest upgrade for application '%s' not implemented yet", config.Name))

	} else {
		// Should be caught by validation, but double-check
		// 应被校验捕获，但双重检查
		return errors.New(errors.ErrTypeInternal, fmt.Sprintf("no deployment method specified for application '%s' for upgrade", config.Name))
	}

	// Step 3: Optional - Inject configuration (if changed)
	// 步骤 3：可选 - 注入配置（如果更改）
	if len(config.ConfigInjection) > 0 {
		utils.GetLogger().Printf("Injecting configuration for application '%s' after upgrade...", config.Name)
		// TODO: Re-run config injection logic
		utils.GetLogger().Printf("Placeholder: Configuration injected for application '%s' after upgrade.", config.Name)
	}

	// Step 4: Optional - Wait for application rollout to complete
	// 步骤 4：可选 - 等待应用程序 rollout 完成
	utils.GetLogger().Printf("Waiting for application '%s' rollout to complete after upgrade...", config.Name)
	// TODO: Implement rollout status check
	utils.GetLogger().Printf("Placeholder: Application '%s' rollout completed after upgrade.", config.Name)

	utils.GetLogger().Printf("Application Deployer completed upgrade for '%s'.", config.Name)
	return nil
}

// Delete deletes the deployed application.
// Delete 删除已部署的应用程序。
func (d *defaultDeployer) Delete(ctx context.Context, appName string, vclusterName string, namespace string, hostK8sClient kubernetes.Interface) error {
	utils.GetLogger().Printf("Deleting application '%s' from vcluster '%s' in namespace '%s'",
		appName, vclusterName, namespace)

	// Step 1: Get the client for the target vcluster
	// 步骤 1：获取目标 vcluster 的客户端
	utils.GetLogger().Printf("Getting client for vcluster '%s'...", vclusterName)
	vclusterClient, err := vcluster_client.GetVClusterClient(ctx, vclusterName, hostK8sClient)
	if err != nil {
		return errors.NewWithCause(errors.ErrTypeVCluster, fmt.Sprintf("failed to get client for target vcluster '%s' for deletion", vclusterName), err)
	}
	utils.GetLogger().Printf("Successfully obtained client for vcluster '%s'.", vclusterName)

	// Step 2: Determine the deployment method used and execute deletion
	// 步骤 2：确定使用的部署方法并执行删除
	// This requires knowing HOW the application was deployed. This info might need to be stored in a state file or labels/annotations.
	// 这需要知道应用程序是如何部署的。此信息可能需要存储在状态文件或标签/注解中。
	// For now, let's assume we can infer the method (less robust).
	// 现在，让我们假设可以推断出方法（不够健壮）。

	// TODO: Infer deployment method (Helm, Kustomize, Manifests)
	// Based on stored state or labels/annotations on resources in the vcluster.
	// Example: Check for Helm release secrets in the vcluster namespace.
	// 根据 vcluster 中资源上的存储状态或标签/注解。
	// 示例：检查 vcluster 命名空间中的 Helm release secret。

	// Placeholder for inferred method
	inferredMethod := "helm" // Assume Helm for example

	if inferredMethod == "helm" {
		utils.GetLogger().Printf("Deleting application '%s' using Helm release name '%s'...", appName, appName) // Assuming app name is release name
		// Helm deletion needs the vcluster client and the release name
		// Helm 删除需要 vcluster 客户端和 release 名称
		err := d.helmDeployer.DeleteRelease(ctx, vclusterClient, appName, namespace) // Assuming DeleteRelease method exists
		if err != nil {
			return errors.NewWithCause(errors.ErrTypeApplication, fmt.Sprintf("failed to delete application '%s' using Helm", appName), err)
		}
		utils.GetLogger().Printf("Application '%s' deleted successfully using Helm.", appName)

	} else if inferredMethod == "kustomize" {
		utils.GetLogger().Printf("Placeholder: Deleting application '%s' using Kustomize...", appName)
		// TODO: Implement Kustomize deletion logic (e.g., kustomize prune or delete applied manifests)
		// err := d.kustomizeDeployer.Delete(ctx, vclusterClient, appName, namespace)
		return errors.New(errors.ErrTypeNotImplemented, fmt.Sprintf("Kustomize deletion for application '%s' not implemented yet", appName))

	} else if inferredMethod == "manifests" {
		utils.GetLogger().Printf("Placeholder: Deleting application '%s' using raw manifests...", appName)
		// TODO: Implement raw manifest deletion logic (e.g., delete applied manifests)
		// Requires storing information about the resources created by the manifests.
		// 需要存储有关由 manifests 创建的资源的信息。
		// err := d.manifestDeployer.Delete(ctx, vclusterClient, appName, namespace)
		return errors.New(errors.ErrTypeNotImplemented, fmt.Sprintf("raw manifest deletion for application '%s' not implemented yet", appName))

	} else {
		return errors.New(errors.ErrTypeInternal, fmt.Sprintf("could not determine deployment method for application '%s' in vcluster '%s'", appName, vclusterName))
	}

	utils.GetLogger().Printf("Application Deployer completed deletion for '%s'.", appName)
	return nil
}

// Placeholder struct for application status if needed
// 应用程序状态的占位符结构体（如果需要）
// type ApplicationStatus struct {
// 	Name    string `json:"name"`
// 	VCluster string `json:"vcluster"`
// 	Namespace string `json:"namespace"`
// 	Ready   bool   `json:"ready"`
// 	Message string `json:"message"`
// 	// Add more status details (e.g., replica count, rollout history)
// 	// 添加更多状态详情（例如，副本数、rollout 历史）
// }

// TODO: Implement helper functions like waitForDeploymentRollout, updateConfigMapInVCluster etc.
// TODO: 实现辅助函数，例如 waitForDeploymentRollout, updateConfigMapInVCluster 等。
// TODO: Implement functions to infer application deployment method
// TODO: 实现推断应用程序部署方法的函数
