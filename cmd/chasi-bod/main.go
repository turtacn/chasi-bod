// main package is the entry point for the chasi-bod command-line application.
// main 包是 chasi-bod 命令行应用程序的入口点。
package main

import (
	"log"
	"os"

	"github.com/turtacn/chasi-bod/cmd/chasi-bod/cli" // Assuming cli package exists // 假设 cli 包存在
	"github.com/turtacn/chasi-bod/common/utils"      // Assuming common.utils.InitLogger exists // 假设 common.utils.InitLogger 存在
)

// main is the entry function for the chasi-bod CLI.
// main 是 chasi-bod CLI 的入口函数。
func main() {
	// Initialize the global logger early
	// 尽早初始化全局日志记录器
	// You might want to configure logging output based on flags later
	// 稍后您可能希望根据标志配置日志记录输出
	utils.InitLogger("chasi-bod: ", log.LstdFlags|log.Lshortfile)

	// Execute the root command of the CLI
	// 执行 CLI 的根命令
	if err := cli.Execute(); err != nil {
		// Use the initialized logger to report the error
		// 使用已初始化的日志记录器报告错误
		// Cobra's Execute() often handles errors by printing to stderr,
		// but this allows for centralized error handling/formatting.
		// Cobra 的 Execute() 通常通过打印到 stderr 处理错误，
		// 但这允许集中处理/格式化错误。
		utils.GetLogger().Fatalf("Error executing command: %v", err)
		os.Exit(1) // Exit with a non-zero status code on error
	}

	// Exit with a zero status code on success
	// 成功时以零状态码退出
	os.Exit(0)
}
