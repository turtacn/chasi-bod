// Package descriptor provides functionality for processing chasi-bod application configuration descriptors.
// 包 descriptor 提供了处理 chasi-bod 应用程序配置描述符的功能。
// It is responsible for loading, validating, and potentially pre-processing application-specific configurations.
// 它负责加载、校验以及可能预处理应用程序特定配置。
package descriptor

import (
	"fmt"
	"os"
	//"path/filepath" // Added for path joining // 添加用于路径拼接

	"github.com/turtacn/chasi-bod/common/errors"
	"github.com/turtacn/chasi-bod/common/utils" // Assuming logger is here // 假设日志记录器在这里
	"github.com/turtacn/chasi-bod/pkg/config/model"
	config_validator "github.com/turtacn/chasi-bod/pkg/config/validator" // Import config validator // 导入配置校验器
	"gopkg.in/yaml.v2"                                                   // Using yaml.v2 for parsing // 使用 yaml.v2 进行解析
)

// ApplicationDeploymentPlan represents the structured plan derived from an application configuration.
// ApplicationDeploymentPlan 表示从应用程序配置派生的结构化计划。
// It contains the parsed and validated information needed by the application deployer.
// 它包含应用程序 deployer 所需的已解析和已校验的信息。
type ApplicationDeploymentPlan struct {
	Config *model.ApplicationConfig // The original configuration / 原始配置
	// Add any pre-processed or derived data here
	// 在这里添加任何预处理或派生数据
	// e.g., Parsed Helm values, Kustomize build output, list of manifests to apply
	// 例如，已解析的 Helm values，Kustomize 构建输出，要应用的 manifests 列表
	ProcessedManifests interface{} // Placeholder for processed deployment artifacts / 处理后的部署 artifact 的占位符
}

// LoadApplicationConfig loads and validates an application configuration file.
// LoadApplicationConfig 加载并校验应用程序配置文件。
// filePath: The path to the application configuration file (e.g., YAML). / 应用程序配置文件的路径（例如，YAML）。
// Returns the loaded and validated ApplicationConfig and an error.
// 返回已加载和已校验的 ApplicationConfig 和错误。
// Note: This performs *standalone* validation. Full validation requires the context of the main PlatformConfig.
// 注意：此函数执行独立校验。完整校验需要主 PlatformConfig 的上下文。
func LoadApplicationConfig(filePath string) (*model.ApplicationConfig, error) {
	utils.GetLogger().Printf("Loading application configuration from %s", filePath)

	// Read the file content
	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New(errors.ErrTypeNotFound, fmt.Sprintf("application configuration file not found at %s", filePath))
		}
		return nil, errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to read application configuration file %s", filePath), err)
	}

	var appConfig model.ApplicationConfig
	// Unmarshal the YAML content into the struct
	// 将 YAML 内容反序列化到结构体中
	err = yaml.Unmarshal(data, &appConfig)
	if err != nil {
		return nil, errors.NewWithCause(errors.ErrTypeConfig, fmt.Sprintf("failed to unmarshal application configuration file %s", filePath), err)
	}

	// Perform standalone validation on the loaded application config
	// 对已加载的应用程序配置执行独立校验
	// This checks the structure and basic field validity within the application config itself.
	// It does NOT check against the main platform config (e.g., if vclusterName exists).
	// 这检查应用程序配置本身的结构和基本字段有效性。
	// 它不检查与主平台配置的对照（例如，vclusterName 是否存在）。
	if err := config_validator.ValidateApplicationConfig("standalone", &appConfig); err != nil { // Using a placeholder name for context
		return nil, fmt.Errorf("application configuration standalone validation failed for file %s: %w", filePath, err) // Wrap the validation error
	}

	utils.GetLogger().Printf("Successfully loaded and validated application configuration from %s (standalone).", filePath)

	return &appConfig, nil
}

// ProcessApplicationConfig takes a validated application configuration and prepares a deployment plan.
// ProcessApplicationConfig 接收一个已校验的应用程序配置，并准备一个部署计划。
// This might involve resolving chart dependencies, building kustomize overlays, etc.
// This function assumes the input config is already validated within the context of the PlatformConfig.
// 这可能涉及解析 chart 依赖项，构建 kustomize overlay 等。
// 此函数假定输入配置已在 PlatformConfig 的上下文中进行校验。
// config: The validated application configuration. / 已校验的应用程序配置。
// Returns an ApplicationDeploymentPlan and an error.
// 返回 ApplicationDeploymentPlan 和错误。
func ProcessApplicationConfig(config *model.ApplicationConfig) (*ApplicationDeploymentPlan, error) {
	utils.GetLogger().Printf("Processing application configuration for '%s'...", config.Name)

	// This function would contain logic to:
	// - For Helm: Resolve chart source (local or remote), potentially download, maybe template?
	// - For Kustomize: Build the kustomization.yaml to produce final manifests
	// - For raw manifests: Just read the files

	// The output `ProcessedManifests` would depend on the deployment method.
	// It could be:
	// - A temporary path to a downloaded Helm chart.
	// - A byte slice of the kustomize build output.
	// - A list of file paths for raw manifests.

	// Currently, it's a placeholder that just wraps the config.
	// For Helm, the `helm` package handles chart loading, so this function might not need to do that here.
	// 对于 Helm，`helm` 包处理 chart 加载，因此此函数可能不需要在此处执行该操作。
	// It might be more about preparing values or other inputs for the deployer.
	// 这可能更多是为 deployer 准备值或其他输入。

	plan := &ApplicationDeploymentPlan{
		Config: config,
		// ProcessedManifests: depends on the deployment method and what the deployer expects.
		// ProcessedManifests: 取决于部署方法和 deployer 期望的内容。
	}

	// Example: If using Kustomize, this is where you'd run the kustomize build command.
	// 示例：如果使用 Kustomize，您将在此处运行 kustomize build 命令。
	// if config.Kustomize != nil {
	// 	kustomizeOutput, err := runKustomizeBuild(config.Kustomize.Path) // Need helper
	// 	if err != nil { return nil, fmt.Errorf("failed to build kustomize manifests: %w", err) }
	// 	plan.ProcessedManifests = kustomizeOutput // Store the built manifests
	// } else if len(config.Manifests) > 0 {
	// 	// For raw manifests, maybe just store the list of paths
	// 	plan.ProcessedManifests = config.Manifests
	// }

	utils.GetLogger().Printf("Application configuration processed for '%s'.", config.Name)
	return plan
}

// TODO: Implement helper functions for running kustomize build, resolving manifest paths etc.
// TODO: 实现用于运行 kustomize build、解析 manifest 路径等的辅助函数
