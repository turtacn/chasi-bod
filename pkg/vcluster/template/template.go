// Package template provides functionality for processing vcluster configuration templates.
// 包 template 提供了处理 vcluster 配置模板的功能。
// These templates are used to generate Kubernetes manifests for deploying vcluster instances.
// 这些模板用于生成用于部署 vcluster 实例的 Kubernetes manifests。
package template

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings" // Added for string operations // 添加用于字符串操作
	"text/template"

	"github.com/turtacn/chasi-bod/common/errors"
	"github.com/turtacn/chasi-bod/common/utils" // Assuming logger and file ops are here // 假设日志记录器和文件操作在这里
	"github.com/turtacn/chasi-bod/pkg/config/model"
	// Assuming constants package is available for accessing default template dir etc.
	// 假设 constants 包可用于访问默认模板目录等。
	"github.com/turtacn/chasi-bod/common/constants"
)

// DefaultTemplateDir is the default directory for vcluster templates.
// DefaultTemplateDir 是 vcluster 模板的默认目录。
// const DefaultTemplateDir = "configs/vcluster" // Moved to constants package // 已移动到 constants 包

// LoadTemplate reads and parses a named vcluster template file.
// LoadTemplate 读取并解析命名为 name 的 vcluster 模板文件。
// name: The base name of the template file (e.g., "basic.yaml"). / 模板文件的基本名称（例如，“basic.yaml”）。
// Returns the parsed template and an error if loading or parsing failed.
// 返回解析后的模板，以及加载或解析失败时的错误。
func LoadTemplate(name string) (*template.Template, error) {
	templatePath := filepath.Join(constants.DefaultTemplateDir, name) // Use constants.DefaultTemplateDir

	exists, err := utils.PathExists(templatePath)
	if err != nil {
		return nil, errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to check existence of template file %s", templatePath), err)
	}
	if !exists {
		return nil, errors.New(errors.ErrTypeNotFound, fmt.Sprintf("vcluster template file not found at %s", templatePath))
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, errors.NewWithCause(errors.ErrTypeConfig, fmt.Sprintf("failed to parse vcluster template file %s", templatePath), err)
	}

	return tmpl, nil
}

// ProcessTemplate applies the vcluster configuration to a template and returns the resulting YAML.
// ProcessTemplate 将 vcluster 配置应用于模板，并返回结果 YAML。
// tmpl: The parsed template. / 解析后的模板。
// config: The vcluster configuration to use as data for the template. / 用作模板数据的 vcluster 配置。
// Returns the processed YAML content as a byte slice and an error if execution failed.
// 返回处理后的 YAML 内容作为字节切片，以及执行失败时的错误。
func ProcessTemplate(tmpl *template.Template, config *model.VClusterConfig) ([]byte, error) {
	var buf bytes.Buffer
	// Execute the template with the vcluster configuration data
	// 使用 vcluster 配置数据执行模板
	err := tmpl.Execute(&buf, config)
	if err != nil {
		return nil, errors.NewWithCause(errors.ErrTypeConfig, fmt.Sprintf("failed to execute vcluster template '%s'", tmpl.Name()), err) // Include template name in error
	}

	return buf.Bytes(), nil
}

// LoadAndProcessTemplate combines loading and processing a template.
// LoadAndProcessTemplate 组合加载和处理模板。
// templateName: The base name of the template file. / 模板文件的基本名称。
// config: The vcluster configuration. / vcluster 配置。
// Returns the processed YAML content and an error.
// 返回处理后的 YAML 内容和错误。
func LoadAndProcessTemplate(templateName string, config *model.VClusterConfig) ([]byte, error) {
	tmpl, err := LoadTemplate(templateName)
	if err != nil {
		return nil, err // Errors are already wrapped by LoadTemplate
	}

	processedContent, err := ProcessTemplate(tmpl, config)
	if err != nil {
		return nil, err // Errors are already wrapped by ProcessTemplate
	}

	return processedContent, nil
}

// ListTemplates lists available vcluster templates in the default template directory.
// ListTemplates 列出默认模板目录中可用的 vcluster 模板。
// Returns a slice of template names and an error.
// 返回模板名称的切片和错误。
func ListTemplates() ([]string, error) {
	entries, err := os.ReadDir(constants.DefaultTemplateDir) // Use constants.DefaultTemplateDir
	if err != nil {
		// If the directory doesn't exist, return an empty list and no error
		// 如果目录不存在，返回空列表且没有错误
		if os.IsNotExist(err) {
			utils.GetLogger().Printf("Vcluster template directory not found at %s. Returning empty list.", constants.DefaultTemplateDir)
			return []string{}, nil
		}
		return nil, errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to read template directory %s", constants.DefaultTemplateDir), err)
	}

	var templateNames []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") { // Assuming templates end with .yaml
			templateNames = append(templateNames, entry.Name())
		}
	}
	utils.GetLogger().Printf("Found %d vcluster template files in %s", len(templateNames), constants.DefaultTemplateDir)
	return templateNames, nil
}
