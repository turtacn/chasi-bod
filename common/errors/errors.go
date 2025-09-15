// Package errors provides centralized error types and handling utilities.
// 包 errors 提供了集中的错误类型和处理工具。
package errors

import (
	"fmt"
)

// ErrorType represents a specific category of error.
// ErrorType 表示特定类别的错误。
type ErrorType string

const (
	// ErrTypeConfig indicates an error related to configuration.
	// ErrTypeConfig 表示与配置相关的错误。
	ErrTypeConfig ErrorType = "config"
	// ErrTypeValidation indicates a validation error.
	// ErrTypeValidation 表示校验错误。
	ErrTypeValidation ErrorType = "validation"
	// ErrTypeIO indicates an input/output error.
	// ErrTypeIO 表示输入/输出错误。
	ErrTypeIO ErrorType = "io"
	// ErrTypeNetwork indicates a network error.
	// ErrTypeNetwork 表示网络错误。
	ErrTypeNetwork ErrorType = "network"
	// ErrTypeSystem indicates an underlying system error (e.g., OS command failure).
	// ErrTypeSystem 表示底层系统错误（例如，操作系统命令失败）。
	ErrTypeSystem ErrorType = "system"
	// ErrTypeNotFound indicates a resource was not found.
	// ErrTypeNotFound 表示资源未找到。
	ErrTypeNotFound ErrorType = "not_found"
	// ErrTypeAlreadyExists indicates a resource already exists.
	// ErrTypeAlreadyExists 表示资源已存在。
	ErrTypeAlreadyExists ErrorType = "already_exists"
	// ErrTypeTimeout indicates an operation timed out.
	// ErrTypeTimeout 表示操作超时。
	ErrTypeTimeout ErrorType = "timeout"
	// ErrTypeNotImplemented indicates a feature is not implemented yet.
	// ErrTypeNotImplemented 表示功能尚未实现。
	ErrTypeNotImplemented ErrorType = "not_implemented"
	// ErrTypeInternal indicates an internal unexpected error.
	// ErrTypeInternal 表示内部意外错误。
	ErrTypeInternal ErrorType = "internal"
	// ErrTypeVCluster indicates an error related to vcluster operations.
	// ErrTypeVCluster 表示与 vcluster 操作相关的错误。
	ErrTypeVCluster ErrorType = "vcluster"
	// ErrTypeApplication indicates an error related to application deployment or management.
	// ErrTypeApplication 表示与应用程序部署或管理相关的错误。
	ErrTypeApplication ErrorType = "application"
	// ErrTypeDFX indicates an error related to DFX (Design for X) components like logging, metrics, tracing.
	// ErrTypeDFX 表示与 DFX（可观测性设计）组件（如日志、指标、追踪）相关的错误。
	ErrTypeDFX ErrorType = "dfx"
	// ErrTypeReliability indicates an error related to reliability components like backup, restore.
	// ErrTypeReliability 表示与可靠性组件（如备份、恢复）相关的错误。
	ErrTypeReliability ErrorType = "reliability"
)

// ChasiBodError is a custom error type for chasi-bod errors.
// ChasiBodError 是 chasi-bod 错误的自定义错误类型。
type ChasiBodError struct {
	Type    ErrorType `json:"type"`    // Error category / 错误类别
	Message string    `json:"message"` // Human-readable error message / 可读的错误消息
	Cause   error     `json:"cause"`   // Underlying error if any / 底层错误（如果有）
}

// Error returns the string representation of the error.
// Error 返回错误的字符串表示。
func (e *ChasiBodError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

// Unwrap returns the underlying cause error.
// Unwrap 返回底层原因错误。
func (e *ChasiBodError) Unwrap() error {
	return e.Cause
}

// New creates a new ChasiBodError.
// New 创建一个新的 ChasiBodError。
func New(errType ErrorType, message string) error {
	return &ChasiBodError{
		Type:    errType,
		Message: message,
	}
}

// NewWithCause creates a new ChasiBodError with a cause.
// NewWithCause 创建一个新的带有原因的 ChasiBodError。
func NewWithCause(errType ErrorType, message string, cause error) error {
	return &ChasiBodError{
		Type:    errType,
		Message: message,
		Cause:   cause,
	}
}

// IsChasiBodError checks if an error is a ChasiBodError of a specific type.
// IsChasiBodError 检查错误是否是特定类型的 ChasiBodError。
func IsChasiBodError(err error, errType ErrorType) bool {
	if cbErr, ok := err.(*ChasiBodError); ok {
		return cbErr.Type == errType
	}
	return false
}

// Is checks if an error matches a target error.
// This is part of the Go 1.13 errors.Is pattern.
// Is 检查错误是否与目标错误匹配。
// 这是 Go 1.13 errors.Is 模式的一部分。
func (e *ChasiBodError) Is(target error) bool {
	// Allow checking against another ChasiBodError of the same type
	// 允许检查同一类型的另一个 ChasiBodError
	if targetErr, ok := target.(*ChasiBodError); ok {
		return e.Type == targetErr.Type
	}
	// Fallback to checking the cause if target is not a ChasiBodError
	// 如果目标不是 ChasiBodError，则回退检查原因
	if e.Cause != nil {
		return e.Cause == target // Use standard error equality check
	}
	return false
}
