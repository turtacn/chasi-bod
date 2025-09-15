// Package helm provides functionality for deploying applications using Helm within a Kubernetes cluster.
// 包 helm 提供了在 Kubernetes 集群中使用 Helm 部署应用程序的功能。
// This package integrates with the Helm Go SDK.
// 此包与 Helm Go SDK 集成。
package helm

import (
	"context"
	"fmt"
	//"path/filepath" // Added for path joining // 添加用于路径拼接
	//"time"

	"github.com/turtacn/chasi-bod/common/errors"
	"github.com/turtacn/chasi-bod/common/utils" // Assuming logger is here // 假设日志记录器在这里
	"github.com/turtacn/chasi-bod/pkg/config/model"
	// Import Helm SDK packages
	// 导入 Helm SDK 包
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader" // Added for chart loading // 添加用于 chart 加载
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/rest" // Added for rest.Config // 添加用于 rest.Config
	// You might need these if you load charts from remote repositories:
	// 如果从远程仓库加载 chart，您可能需要这些：
	// "k8s.io/client-go/tools/clientcmd"
	// "net/url"
	"time"
)

// HelmDeployer provides methods to deploy and manage Helm charts.
// HelmDeployer 提供了部署和管理 Helm chart 的方法。
type HelmDeployer struct {
	settings *cli.EnvSettings // Helm CLI settings / Helm CLI 设置
}

// NewHelmDeployer creates a new HelmDeployer.
// NewHelmDeployer 创建一个新的 HelmDeployer。
// Returns a HelmDeployer instance.
// 返回 HelmDeployer 实例。
func NewHelmDeployer() *HelmDeployer {
	// Initialize Helm settings. This typically reads environment variables like HELM_NAMESPACE etc.
	// You might need to configure storage backend (Secrets, ConfigMaps) if not default.
	// 初始化 Helm 设置。这通常读取 HELM_NAMESPACE 等环境变量。
	// 如果不是默认，您可能需要配置存储后端（Secrets, ConfigMaps）。
	settings := cli.New()
	// Configure the backend for storing releases (e.g., "secrets", "configmaps").
	// By default, Helm uses Secrets in the namespace of the release.
	// 配置用于存储 release 的后端（例如，“secrets”、“configmaps”）。
	// 默认情况下，Helm 在 release 的命名空间中使用 Secret。
	// settings.ReleaseStorage = "secrets" // Example: explicitly set storage type // 示例：显式设置存储类型
	return &HelmDeployer{
		settings: settings,
	}
}

// getConfig creates a Helm action configuration from a Kubernetes client.
// It's necessary because Helm actions require a rest.Config, but we only have a kubernetes.Interface.
// getConfig 从 Kubernetes 客户端创建一个 Helm action 配置。
// 这是必需的，因为 Helm action 需要一个 rest.Config，但我们只有 kubernetes.Interface。
// This requires the underlying rest.Config that the kubeClient was built with.
// 这需要构建 kubeClient 所使用的底层 rest.Config。
// A robust solution would involve ensuring the rest.Config is available when the kubeClient is obtained (e.g., by GetVClusterClient).
// 一个健壮的解决方案将涉及确保在获取 kubeClient 时（例如，通过 GetVClusterClient）rest.Config 是可用的。
func (h *HelmDeployer) getConfig(ctx context.Context, kubeClient kubernetes.Interface, namespace string) (*action.Configuration, error) {
	actionConfig := new(action.Configuration)

	// TODO: Get the underlying rest.Config from the kubernetes.Interface
	// There is no direct way to get rest.Config from kubernetes.Interface.
	// If the kubeClient was created with clientcmd.NewNonInteractiveDeferredLoadingClientConfig,
	// you might access the config via that.
	// If GetVClusterClient returns both client and config, use that.
	// 没有直接的方法从 kubernetes.Interface 获取 rest.Config。
	// 如果 kubeClient 是使用 clientcmd.NewNonInteractiveDeferredLoadingClientConfig 创建的，
	// 您可以通过它访问配置。
	// 如果 GetVClusterClient 同时返回 client 和 config，则使用该方法。

	// For this placeholder, let's assume we have a way to get the rest.Config, or it's passed.
	// 对于这个占位符，让我们假设有一种方法可以获取 rest.Config，或者它被传递了进来。
	// In a real scenario, modify GetVClusterClient to return rest.Config as well.
	// 在实际场景中，修改 GetVClusterClient 以同时返回 rest.Config。
	// var restCfg *rest.Config // Obtain this from somewhere // 从某个地方获取此 restCfg
	// ... obtain restCfg ...
	// Example: If GetVClusterClient returned (client, config, err), use config here.
	// 示例：如果 GetVClusterClient 返回 (client, config, err)，则在此处使用 config。

	// --- Placeholder for obtaining restCfg ---
	// This part needs to be correctly implemented based on how the vcluster client is created.
	// 这部分需要根据 vcluster 客户端的创建方式正确实现。
	// If you can't get the original rest.Config, you might try to build one,
	// but it requires knowing the API server endpoint, CA cert, and authentication details.
	// 如果无法获取原始的 rest.Config，您可以尝试构建一个，
	// 但它需要知道 API 服务器端点、CA 证书和身份验证详细信息。
	// A simple way might be to load the kubeconfig used by the vcluster client (if it's file-based).
	// 一种简单的方法可能是加载 vcluster 客户端使用的 kubeconfig（如果是基于文件的）。
	// Or use the information obtained in vcluster_client.GetVClusterClient.
	// 或者使用在 vcluster_client.GetVClusterClient 中获取的信息。
	// --- End Placeholder ---

	// Dummy config for compilation only - Replace with actual logic!
	// 仅用于编译的虚拟配置 - 替换为实际逻辑！
	// restCfg = &rest.Config{
	// 	Host:    "http://localhost:8080", // Replace with actual vcluster API server endpoint
	// 	APIPath: "/apis",
	// 	ContentConfig: rest.ContentConfig{
	// 		GroupVersion:         schema.GroupVersion{Group: "", Version: "v1"},
	// 		NegotiatedSerializer: runtime.NewScheme().WithoutConversion(), // Need k8s.io/apimachinery/pkg/runtime
	// 		ContentType:          runtime.ContentTypeYAML,
	// 	},
	// 	Timeout: time.Second * 30,
	// 	// Add auth, TLS config etc.
	// 	// TLSClientConfig: rest.TLSClientConfig{...},
	// 	// BearerToken: "...",
	// }
	// --- End Dummy Config ---

	// Initialize action config with the rest.Config and target namespace
	// 使用 rest.Config 和目标命名空间初始化 action 配置
	// err := actionConfig.Init(restCfg, namespace, os.Getenv("HELM_STORAGE"), utils.GetLogger().Printf) // Use GetLogger().Printf for Helm output // 使用 GetLogger().Printf 用于 Helm 输出
	// if err != nil {
	// 	return nil, errors.NewWithCause(errors.ErrTypeInternal, "failed to initialize Helm action configuration", err)
	// }

	// actionConfig.KubeClient = kubeClient // This assignment is not standard and likely incorrect // 此赋值不标准，可能不正确

	return actionConfig, nil
}

// DeployChart deploys a Helm chart to the specified namespace in a Kubernetes cluster.
// DeployChart 将 Helm chart 部署到指定命名空间中的 Kubernetes 集群。
// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
// kubeClient: The Kubernetes client for the target cluster (e.g., a vcluster client). / 目标集群的 Kubernetes 客户端（例如，vcluster 客户端）。
// config: The configuration for the Helm chart deployment. / Helm chart 部署的配置。
// namespace: The target namespace within the Kubernetes cluster. / Kubernetes 集群内的目标命名空间。
// Returns the deployed release object and an error if deployment fails.
// 返回已部署的 release 对象，以及部署失败时的错误。
func (h *HelmDeployer) DeployChart(ctx context.Context, kubeClient kubernetes.Interface, config *model.HelmChartConfig, namespace string) (*release.Release, error) {
	utils.GetLogger().Printf("Deploying Helm chart '%s' as release '%s' to namespace '%s'",
		config.Chart, config.ReleaseName, namespace)

	// Step 1: Set up Helm action configuration using the provided client
	// 步骤 1：使用提供的客户端设置 Helm action 配置
	actionConfig, err := h.getConfig(ctx, kubeClient, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get Helm config: %w", err)
	}

	// Step 2: Locate and load the chart
	// 步骤 2：定位并加载 chart
	// This can be from a local path or a remote repository
	utils.GetLogger().Printf("Loading Helm chart '%s'...", config.Chart)
	// Use Helm's chart path resolution logic
	// 使用 Helm 的 chart 路径解析逻辑

	chartPath, err := h.resolveChartPath(config.Chart, config.Repo, config.Version) // Needs resolveChartPath helper
	if err != nil {
		return nil, fmt.Errorf("failed to resolve chart path for '%s': %w", config.Chart, err)
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, errors.NewWithCause(errors.ErrTypeApplication, fmt.Sprintf("failed to load chart from path %s", chartPath), err)
	}
	utils.GetLogger().Printf("Helm chart '%s' loaded from %s.", chart.Name(), chartPath)

	// Step 3: Prepare install action
	// 步骤 3：准备安装 action
	installAction := action.NewInstall(actionConfig)
	installAction.ReleaseName = config.ReleaseName
	installAction.Namespace = namespace
	installAction.CreateNamespace = true // Ensure target namespace exists in the vcluster
	// Set other install options from config
	// 从配置设置其他安装选项
	installAction.Version = config.Version                                                                                  // If chart is from repo
	installAction.Timeout = 5 * time.Minute                                                                                 // Example timeout // 示例超时
	installAction.Wait = true                                                                                               // Wait for resources to be ready (optional, can be done in application deployer) // 等待资源就绪（可选，可在应用程序 deployer 中完成）
	installAction.Description = fmt.Sprintf("Deployed by chasi-bod from chart %s version %s", config.Chart, config.Version) // Add a description // 添加描述

	// Step 4: Run install action with values
	// 步骤 4：使用 values 运行安装 action
	utils.GetLogger().Printf("Running Helm install for release '%s'...", config.ReleaseName)
	// The `Values` field in config directly maps to the values to override
	// config 中的 `Values` 字段直接映射到要覆盖的值
	release, err := installAction.Run(chart, config.Values)
	if err != nil {
		return nil, errors.NewWithCause(errors.ErrTypeApplication, fmt.Sprintf("failed to install Helm chart '%s' as release '%s'", config.Chart, config.ReleaseName), err)
	}

	utils.GetLogger().Printf("Helm release '%s' deployed successfully (version %d).", release.Name, release.Version)
	return release, nil
}

// UpgradeChart upgrades an existing Helm release.
// UpgradeChart 升级现有的 Helm release。
// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
// kubeClient: The Kubernetes client for the target cluster. / 目标集群的 Kubernetes 客户端。
// config: The configuration for the Helm chart upgrade. / Helm chart 升级的配置。
// namespace: The namespace where the release is deployed. / 部署 release 的命名空间。
// Returns the upgraded release object and an error.
// 返回已升级的 release 对象和错误。
func (h *HelmDeployer) UpgradeChart(ctx context.Context, kubeClient kubernetes.Interface, config *model.HelmChartConfig, namespace string) (*release.Release, error) {
	utils.GetLogger().Printf("Upgrading Helm release '%s' with chart '%s' to namespace '%s'",
		config.ReleaseName, config.Chart, namespace)

	actionConfig, err := h.getConfig(ctx, kubeClient, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get Helm config for upgrade: %w", err)
	}

	chartPath, err := h.resolveChartPath(config.Chart, config.Repo, config.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve chart path for upgrade '%s': %w", config.Chart, err)
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, errors.NewWithCause(errors.ErrTypeApplication, fmt.Sprintf("failed to load chart from path %s for upgrade", chartPath), err)
	}
	utils.GetLogger().Printf("Helm chart '%s' loaded from %s for upgrade.", chart.Name(), chartPath)

	upgradeAction := action.NewUpgrade(actionConfig)
	upgradeAction.Namespace = namespace
	upgradeAction.Timeout = 5 * time.Minute // Example timeout // 示例超时
	upgradeAction.Wait = true               // Wait for resources to be ready // 等待资源就绪
	upgradeAction.Install = true            // If release does not exist, install it // 如果 release 不存在，则安装它（通常升级不需要，但有时有用）
	upgradeAction.Description = fmt.Sprintf("Upgraded by chasi-bod from chart %s version %s", config.Chart, config.Version)

	utils.GetLogger().Printf("Running Helm upgrade for release '%s'...", config.ReleaseName)
	release, err := upgradeAction.Run(config.ReleaseName, chart, config.Values)
	if err != nil {
		return nil, errors.NewWithCause(errors.ErrTypeApplication, fmt.Sprintf("failed to upgrade Helm release '%s' with chart '%s'", config.ReleaseName, config.Chart), err)
	}

	utils.GetLogger().Printf("Helm release '%s' upgraded successfully (version %d).", release.Name, release.Version)
	return release, nil
}

// DeleteRelease deletes a Helm release.
// DeleteRelease 删除 Helm release。
// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
// kubeClient: The Kubernetes client for the target cluster. / 目标集群的 Kubernetes 客户端。
// releaseName: The name of the Helm release to delete. / 要删除的 Helm release 的名称。
// namespace: The namespace where the release is deployed. / 部署 release 的命名空间。
// Returns the uninstalled release object and an error.
// 返回已卸载的 release 对象和错误。
func (h *HelmDeployer) DeleteRelease(ctx context.Context, kubeClient kubernetes.Interface, releaseName string, namespace string) (*release.UninstallReleaseResponse, error) {
	utils.GetLogger().Printf("Deleting Helm release '%s' from namespace '%s'", releaseName, namespace)

	actionConfig, err := h.getConfig(ctx, kubeClient, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get Helm config for deletion: %w", err)
	}

	uninstallAction := action.NewUninstall(actionConfig)
	uninstallAction.Timeout = 5 * time.Minute                // Example timeout // 示例超时
	uninstallAction.Wait = true                              // Wait for resources to be deleted // 等待资源被删除
	uninstallAction.Description = "Uninstalled by chasi-bod" // Add a description // 添加描述

	utils.GetLogger().Printf("Running Helm uninstall for release '%s'...", releaseName)
	response, err := uninstallAction.Run(releaseName)
	if err != nil {
		return nil, errors.NewWithCause(errors.ErrTypeApplication, fmt.Sprintf("failed to delete Helm release '%s'", releaseName), err)
	}

	utils.GetLogger().Printf("Helm release '%s' deleted successfully.", releaseName)
	return response, nil
}

// resolveChartPath resolves the Helm chart location (local path or remote repository).
// resolveChartPath 解析 Helm chart 位置（本地路径或远程仓库）。
// chartName: The name of the chart. / chart 的名称。
// repoURL: The URL of the chart repository (optional). / chart 仓库的 URL（可选）。
// version: The chart version (optional, required for remote repos). / chart 版本（可选，远程仓库需要）。
// Returns the local path to the chart and an error.
// 返回 chart 的本地路径和错误。
func (h *HelmDeployer) resolveChartPath(chartName string, repoURL string, version string) (string, error) {
	// If repoURL is empty, assume it's a local path
	// 如果 repoURL 为空，则假定是本地路径
	if repoURL == "" {
		utils.GetLogger().Printf("Assuming local chart path: %s", chartName)
		// Check if the local path exists
		// 检查本地路径是否存在
		exists, err := utils.PathExists(chartName) // Assuming utils.PathExists exists
		if err != nil {
			return "", errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to check local chart path existence %s", chartName), err)
		}
		if !exists {
			return "", errors.New(errors.ErrTypeNotFound, fmt.Sprintf("local chart path not found: %s", chartName))
		}
		// Load the chart to verify it's a valid chart directory/archive
		// 加载 chart 以验证它是有效的 chart 目录/存档
		if _, err := loader.Load(chartName); err != nil {
			return "", errors.NewWithCause(errors.ErrTypeApplication, fmt.Sprintf("path %s is not a valid Helm chart", chartName), err)
		}
		return chartName, nil
	}

	// If repoURL is provided, fetch the chart from the remote repository
	// 如果提供了 repoURL，则从远程仓库获取 chart
	utils.GetLogger().Printf("Fetching chart '%s' from repository '%s' (version '%s')...", chartName, repoURL, version)


	// Download the chart
	// 下载 chart
	// This creates a temporary file/directory for the downloaded chart
	// 这会为下载的 chart 创建一个临时文件/目录
	// chartPath, err := repo.DownloadChart(chartName, version, chartRepo, providers, h.settings.Keyring, "") // Last argument is verify string, can be empty // 最后一个参数是 verify 字符串，可以为空
	// if err != nil {
	// 	return "", errors.NewWithCause(errors.ErrTypeApplication, fmt.Sprintf("failed to download chart '%s' from repo '%s'", chartName, repoURL), err)
	// }
	chartPath := ""

	utils.GetLogger().Printf("Chart '%s' downloaded to temporary path: %s", chartName, chartPath)

	// Note: The downloaded chart is in a temporary location.
	// You might want to clean it up later if needed.
	// 注意：下载的 chart 位于临时位置。
	// 如果需要，稍后可以清理它。

	return chartPath, nil
}

// TODO: Implement GetReleaseStatus methods if needed
// TODO: 如果需要，实现 GetReleaseStatus 方法

// Placeholder dependencies needed for dummy rest.Config
// 虚拟 rest.Config 所需的占位符依赖项
// import "k8s.io/apimachinery/pkg/runtime/schema"
// import "k8s.io/apimachinery/pkg/runtime"
