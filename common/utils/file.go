// Package utils provides common utility functions.
// 包 utils 提供了常用的工具函数。
package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/turtacn/chasi-bod/common/errors"
)

// PathExists checks if a file or directory exists at the given path.
// PathExists 检查给定路径是否存在文件或目录。
// path: The path to check. / 要检查的路径。
// Returns true if the path exists, false otherwise, and an error if checking failed.
// 如果路径存在则返回 true，否则返回 false，检查失败时返回错误。
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	// Other error, might be permissions etc.
	// 其他错误，可能是权限等。
	return false, errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to check path existence for %s", path), err)
}

// IsDir checks if the path is a directory.
// IsDir 检查路径是否为目录。
// path: The path to check. / 要检查的路径。
// Returns true if it's a directory, false otherwise, and an error if checking failed or path doesn't exist.
// 如果是目录则返回 true，否则返回 false，检查失败或路径不存在时返回错误。
func IsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, errors.New(errors.ErrTypeNotFound, fmt.Sprintf("path %s not found", path))
		}
		return false, errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to stat path %s", path), err)
	}
	return info.IsDir(), nil
}

// IsFile checks if the path is a regular file.
// IsFile 检查路径是否为普通文件。
// path: The path to check. / 要检查的路径。
// Returns true if it's a file, false otherwise, and an error if checking failed or path doesn't exist.
// 如果是文件则返回 true，否则返回 false，检查失败或路径不存在时返回错误。
func IsFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, errors.New(errors.ErrTypeNotFound, fmt.Sprintf("path %s not found", path))
		}
		return false, errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to stat path %s", path), err)
	}
	return info.Mode().IsRegular(), nil
}

// CopyFile copies a file from src to dest.
// CopyFile 将文件从 src 复制到 dest。
// src: Source file path. / 源文件路径。
// dest: Destination file path. / 目标文件路径。
// Returns an error if the copy failed.
// 复制失败时返回错误。
func CopyFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to open source file %s", src), err)
	}
	defer sourceFile.Close()

	// Ensure destination directory exists
	// 确保目标目录存在
	destDir := filepath.Dir(dest)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to create destination directory %s", destDir), err)
	}

	destFile, err := os.Create(dest)
	if err != nil {
		return errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to create destination file %s", dest), err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to copy file from %s to %s", src, dest), err)
	}

	// Copy file permissions
	// 复制文件权限
	sourceInfo, err := os.Stat(src)
	if err != nil {
		GetLogger().Printf("Warning: Failed to get source file info %s to copy permissions: %v", src, err)
		// Continue without copying permissions if failed
		// 如果失败，则不复制权限并继续
	} else {
		if err := os.Chmod(dest, sourceInfo.Mode()); err != nil {
			GetLogger().Printf("Warning: Failed to copy file permissions to %s: %v", dest, err)
			// Continue if changing permissions failed
			// 如果更改权限失败，则继续
		}
	}

	return nil
}

// CopyDir recursively copies a directory from src to dest.
// CopyDir 递归地将目录从 src 复制到 dest。
// src: Source directory path. / 源目录路径。
// dest: Destination directory path. / 目标目录路径。
// Returns an error if the copy failed.
// 复制失败时返回错误。
func CopyDir(src, dest string) error {
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to stat source directory %s", src), err)
	}

	// Create the destination directory with source permissions
	// 使用源权限创建目标目录
	if err := os.MkdirAll(dest, sourceInfo.Mode()); err != nil {
		return errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to create destination directory %s", dest), err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to read source directory %s", src), err)
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectory
			// 递归复制子目录
			if err := CopyDir(sourcePath, destPath); err != nil {
				return err // Return error from recursive call
			}
		} else {
			// Copy file
			// 复制文件
			if err := CopyFile(sourcePath, destPath); err != nil {
				return err // Return error from CopyFile
			}
		}
	}

	return nil
}

// ReadFileContent reads the content of a file.
// ReadFileContent 读取文件的内容。
// path: The file path. / 文件路径。
// Returns the content as a byte slice and an error if reading failed.
// 返回内容作为字节切片，以及读取失败时的错误。
func ReadFileContent(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to read file %s", path), err)
	}
	return content, nil
}

// WriteFileContent writes content to a file.
// WriteFileContent 将内容写入文件。
// path: The file path. / 文件路径。
// content: The content to write. / 要写入的内容。
// perm: File permissions. / 文件权限。
// Returns an error if writing failed.
// 写入失败时返回错误。
func WriteFileContent(path string, content []byte, perm os.FileMode) error {
	// Ensure destination directory exists
	// 确保目标目录存在
	destDir := filepath.Dir(path)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to create directory %s for writing file", destDir), err)
	}

	err := os.WriteFile(path, content, perm)
	if err != nil {
		return errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to write file %s", path), err)
	}
	return nil
}

// RemovePath removes a file or directory. If it's a directory, it is removed recursively.
// RemovePath 删除文件或目录。如果是目录，则递归删除。
// path: The path to remove. / 要删除的路径。
// Returns an error if removal failed.
// 删除失败时返回错误。
func RemovePath(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to remove path %s", path), err)
	}
	return nil
}

// MkdirAll creates a directory and any necessary parent directories.
// MkdirAll 创建一个目录和所有必要的父目录。
// path: The directory path to create. / 要创建的目录路径。
// perm: The permissions to use for the created directories. / 用于创建目录的权限。
// Returns an error if the creation failed.
// 如果创建失败则返回错误。
func MkdirAll(path string, perm os.FileMode) error {
	err := os.MkdirAll(path, perm)
	if err != nil {
		return errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to create directory path %s", path), err)
	}
	return nil
}
