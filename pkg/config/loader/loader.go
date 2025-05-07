// Package loader provides functionality to load chasi-bod configuration.
// 包 loader 提供了加载 chasi-bod 配置的功能。
package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/turtacn/chasi-bod/common/errors"
	"github.com/turtacn/chasi-bod/common/utils" // Assuming logger is here
	"github.com/turtacn/chasi-bod/pkg/config/model"
	"gopkg.in/yaml.v2" // Using yaml.v2 as a common choice
)

// LoadConfig reads the configuration from the specified file path and unmarshals it into a PlatformConfig struct.
// LoadConfig 从指定的文件路径读取配置，并将其反序列化到 PlatformConfig 结构体中。
// filePath: The path to the configuration file (e.g., YAML). / 配置文件的路径（例如，YAML）。
// Returns the loaded PlatformConfig struct and an error if loading or unmarshalling failed.
// 返回加载的 PlatformConfig 结构体，以及加载或反序列化失败时的错误。
func LoadConfig(filePath string) (*model.PlatformConfig, error) {
	// Read the file content
	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New(errors.ErrTypeNotFound, fmt.Sprintf("configuration file not found at %s", filePath))
		}
		return nil, errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to read configuration file %s", filePath), err)
	}

	var config model.PlatformConfig
	// Unmarshal the YAML content into the struct
	// 将 YAML 内容反序列化到结构体中
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, errors.NewWithCause(errors.ErrTypeConfig, fmt.Sprintf("failed to unmarshal configuration file %s", filePath), err)
	}

	utils.GetLogger().Printf("Successfully loaded configuration from %s", filePath)

	return &config, nil
}

// SaveConfig marshals the PlatformConfig struct into YAML format and writes it to the specified file path.
// SaveConfig 将 PlatformConfig 结构体序列化为 YAML 格式，并写入指定的文件路径。
// config: The PlatformConfig struct to save. / 要保存的 PlatformConfig 结构体。
// filePath: The path to save the configuration file. / 保存配置文件的路径。
// Returns an error if marshalling or writing failed.
// 序列化或写入失败时返回错误。
func SaveConfig(config *model.PlatformConfig, filePath string) error {
	// Marshal the struct into YAML
	// 将结构体序列化为 YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return errors.NewWithCause(errors.ErrTypeConfig, "failed to marshal configuration to YAML", err)
	}

	// Ensure the directory exists
	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to create directory %s for saving config", dir), err)
	}

	// Write the data to the file with appropriate permissions (e.g., 0600 for sensitive configs)
	// 将数据写入文件，并设置适当的权限（例如，对于敏感配置使用 0600）
	// Use 0644 for general config files
	// 对于通用配置文件使用 0644
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to write configuration file %s", filePath), err)
	}

	utils.GetLogger().Printf("Successfully saved configuration to %s", filePath)

	return nil
}
