// Package utils provides common utility functions.
// 包 utils 提供了常用的工具函数。
package utils

import (
	"log"
	"os"
	"sync"
)

// Logger is the global logger instance.
// Logger 是全局日志记录器实例。
var (
	logger *log.Logger
	once   sync.Once
)

// InitLogger initializes the global logger. It should be called early in the application lifecycle.
// InitLogger 初始化全局日志记录器。应在应用程序生命周期的早期调用。
// prefix: The prefix for log messages. / 日志消息的前缀。
// logFlags: Flags to customize the output format (e.g., log.LstdFlags, log.Lshortfile). / 用于自定义输出格式的标志（例如，log.LstdFlags, log.Lshortfile）。
func InitLogger(prefix string, logFlags int) {
	once.Do(func() {
		logger = log.New(os.Stdout, prefix, logFlags)
		// You might want to configure logging to a file here as well
		// 您可能也想在这里配置日志输出到文件
		// logFile, err := os.OpenFile(constants.DefaultLogDir+"/chasi-bod.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		// if err != nil {
		// 	log.Printf("Failed to open log file: %v", err)
		// 	// Continue logging to stdout/stderr
		// } else {
		// 	// You can use io.MultiWriter to log to both stdout and file
		// 	// 您可以使用 io.MultiWriter 同时将日志输出到 stdout 和文件
		// 	// mw := io.MultiWriter(os.Stdout, logFile)
		// 	// logger.SetOutput(mw)
		// }
	})
}

// GetLogger returns the global logger instance. InitLogger must be called first.
// GetLogger 返回全局日志记录器实例。必须先调用 InitLogger。
func GetLogger() *log.Logger {
	if logger == nil {
		// Fallback to default logger or panic, depending on desired behavior
		// 回退到默认日志记录器或 panic，取决于期望的行为
		// log.Println("Warning: Logger not initialized. Using default standard logger.")
		// return log.Default()
		panic("Logger not initialized. Call InitLogger first.")
	}
	return logger
}

// Example log levels (can be expanded with a proper logging library)
// 示例日志级别（可以使用适当的日志库进行扩展）
// func Info(format string, v ...interface{}) { GetLogger().Printf("INFO: "+format, v...) }
// func Warn(format string, v ...interface{}) { GetLogger().Printf("WARN: "+format, v...) }
// func Error(format string, v ...interface{}) { GetLogger().Printf("ERROR: "+format, v...) }
// func Fatal(format string, v ...interface{}) { GetLogger().Fatalf("FATAL: "+format, v...) } // Calls os.Exit(1)
