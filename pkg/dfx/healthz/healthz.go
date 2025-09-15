// Package healthz provides functionality for implementing and checking platform health.
// 包 healthz 提供了实现和检查平台健康的功能。
// This includes exposing platform component health endpoints and running checks on cluster resources.
// 这包括暴露平台组件健康端点和在集群资源上运行检查。
package healthz

import (
	"context"
	"encoding/json" // Added for JSON encoding // 添加用于 JSON 编码
	"fmt"
	"github.com/turtacn/chasi-bod/pkg/vcluster"
	"net/http"
	"time"

	//"net/http"
	//"strings" // Added for string joining // 添加用于字符串拼接
	"sync"
	//"time"

	"github.com/turtacn/chasi-bod/common/errors"
	"github.com/turtacn/chasi-bod/common/utils" // Assuming logger is here // 假设日志记录器在这里
	"github.com/turtacn/chasi-bod/pkg/config/model"
	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/apimachinery/pkg/util/wait" // Needed for polling // 需要用于轮询
)

// Manager defines the interface for managing platform health checks.
// Manager 定义了管理平台健康检查的接口。
type Manager interface {
	// StartHealthCheckServer starts an HTTP server to expose platform component health endpoints.
	// StartHealthCheckServer 启动一个 HTTP 服务器以暴露平台组件健康端点。
	// This is for checking the health of the chasi-bod management plane itself.
	// 这用于检查 chasi-bod 管理平面自身的健康状况。
	StartHealthCheckServer(listenAddr string) error

	// RunClusterChecks performs checks on the state of the Host Cluster and vclusters.
	// RunClusterChecks 在 Host Cluster 和 vcluster 的状态上执行检查。
	// This is for proactive monitoring of the deployed platform infrastructure.
	// 这用于主动监控已部署的平台基础设施。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// hostK8sClient: A Kubernetes client connected to the Host Cluster. / 连接到 Host 集群的 Kubernetes 客户端。
	// config: The platform configuration (needed to identify vclusters etc.). / 平台配置（需要用于识别 vcluster 等）。
	// Returns an error if critical checks fail.
	// 如果关键检查失败则返回错误。
	RunClusterChecks(ctx context.Context, hostK8sClient kubernetes.Interface, config *model.PlatformConfig) error

	// RegisterComponentCheck registers a health checker for a platform component.
	// RegisterComponentCheck 为平台组件注册一个健康检查器。
	RegisterComponentCheck(name string, checker HealthChecker)

	// GetOverallStatus gets the current aggregate health status.
	// GetOverallStatus 获取当前聚合健康状态。
	// It runs the component checks.
	// 它运行组件检查。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// Returns the overall health status.
	// 返回总体健康状态。
	GetOverallStatus(ctx context.Context) *OverallHealthStatus
}

// HealthChecker defines the interface for individual component health checks.
// HealthChecker 定义了单个组件健康检查的接口。
type HealthChecker interface {
	// Check performs a health check and returns a result.
	// Check 执行健康检查并返回结果。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// Returns the check result.
	// 返回检查结果。
	Check(ctx context.Context) HealthCheckResult
	// Name returns the name of the checker. / Name 返回检查器的名称。
	Name() string
}

// HealthCheckResult represents the result of a single health check.
// HealthCheckResult 表示单个健康检查的结果。
type HealthCheckResult struct {
	Name    string `json:"name"`    // Name of the check / 检查名称
	Status  string `json:"status"`  // Status (e.g., "Healthy", "Degraded", "Unhealthy", "Unknown") / 状态（例如，“Healthy”、“Degraded”、“Unhealthy”、“Unknown”）
	Message string `json:"message"` // Optional message / 可选消息
	Error   string `json:"error"`   // Optional error details / 可选错误详情
}

// OverallHealthStatus represents the aggregate health status of the platform.
// OverallHealthStatus 表示平台的总体健康状态。
type OverallHealthStatus struct {
	Status  string              `json:"status"`  // Aggregate status / 聚合状态
	Results []HealthCheckResult `json:"results"` // Results of individual checks / 单个检查的结果
}

// defaultManager is a default implementation of the Healthz Manager.
// defaultManager 是 Healthz Manager 的默认实现。
type defaultManager struct {
	componentCheckers map[string]HealthChecker // Registered component checkers / 已注册的组件检查器
	mu                sync.RWMutex             // Mutex for accessing checkers map / 用于访问 checkers 映射的互斥锁
	healthCheckServer *http.Server             // HTTP server for health endpoint / 健康端点的 HTTP 服务器
}

// NewManager creates a new Healthz Manager.
// NewManager 创建一个新的 Healthz Manager。
// Returns a Healthz Manager implementation.
// 返回 Healthz Manager 实现。
func NewManager() Manager {
	return &defaultManager{
		componentCheckers: make(map[string]HealthChecker),
	}
}

// RegisterComponentCheck registers a health checker for a platform component.
// RegisterComponentCheck 为平台组件注册一个健康检查器。
func (m *defaultManager) RegisterComponentCheck(name string, checker HealthChecker) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.componentCheckers[name]; exists {
		utils.GetLogger().Printf("Warning: Health checker for component '%s' already registered. Overwriting.", name)
	}
	m.componentCheckers[name] = checker
	utils.GetLogger().Printf("Registered health checker for component: %s", name)
}

// StartHealthCheckServer starts an HTTP server to expose component health.
// StartHealthCheckServer 启动一个 HTTP 服务器以暴露组件健康状况。
func (m *defaultManager) StartHealthCheckServer(listenAddr string) error {
	utils.GetLogger().Printf("Starting platform health check server on %s", listenAddr)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", m.handleOverallHealthCheck)
	// Optional: Add individual component endpoints like /healthz/logger
	// 可选：添加单个组件端点，例如 /healthz/logger
	// mux.HandleFunc("/healthz/", m.handleComponentHealthCheck) // Requires parsing path // 需要解析路径

	m.healthCheckServer = &http.Server{
		Addr:    listenAddr,
		Handler: mux,
		// Add timeouts if needed
		// ReadTimeout: 5 * time.Second,
		// WriteTimeout: 10 * time.Second,
		// IdleTimeout: 15 * time.Second,
	}

	// Run the server in a goroutine so it doesn't block
	// 在 goroutine 中运行服务器，使其不阻塞
	go func() {
		utils.GetLogger().Printf("Health check server listening on %s", listenAddr)
		if err := m.healthCheckServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Log the error and potentially signal application shutdown
			// 记录错误并可能发出应用程序关闭信号
			utils.GetLogger().Fatalf("Health check server failed: %v", err)
		}
		utils.GetLogger().Println("Health check server stopped.")
	}()

	utils.GetLogger().Println("Platform health check server started.")
	return nil
}

// StopHealthCheckServer gracefully stops the HTTP health check server.
// StopHealthCheckServer 优雅地停止 HTTP 健康检查服务器。
// ctx: Context for graceful shutdown. / 用于优雅关闭的上下文。
func (m *defaultManager) StopHealthCheckServer(ctx context.Context) error {
	if m.healthCheckServer == nil {
		utils.GetLogger().Println("Health check server not started.")
		return nil
	}
	utils.GetLogger().Println("Stopping health check server...")
	err := m.healthCheckServer.Shutdown(ctx)
	if err != nil {
		return errors.NewWithCause(errors.ErrTypeSystem, "failed to stop health check server", err)
	}
	utils.GetLogger().Println("Health check server shut down gracefully.")
	return nil
}

// handleOverallHealthCheck handles requests to the /healthz endpoint.
// handleOverallHealthCheck 处理对 /healthz 端点的请求。
func (m *defaultManager) handleOverallHealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.GetLogger().Println("Received request to /healthz")

	// Get the overall status by running component checks
	// 通过运行组件检查获取总体状态
	overallStatus := m.GetOverallStatus(r.Context())

	w.Header().Set("Content-Type", "application/json")
	statusCode := http.StatusOK
	if overallStatus.Status != "Healthy" { // Assuming "Healthy" is the desired state
		statusCode = http.StatusInternalServerError // Use 500 for unhealthy status
	}
	w.WriteHeader(statusCode)

	// Marshal results to JSON and write to response
	// 将结果序列化为 JSON 并写入响应
	err := json.NewEncoder(w).Encode(overallStatus)
	if err != nil {
		utils.GetLogger().Printf("Error encoding health check response: %v", err)
		// Attempt to write a generic error response
		// 尝试写入通用错误响应
		http.Error(w, `{"status": "Unknown", "message": "Error encoding response"}`, http.StatusInternalServerError)
	}

	utils.GetLogger().Printf("/healthz response: Status=%s, Results Count=%d", overallStatus.Status, len(overallStatus.Results))
}

// GetOverallStatus gets the current aggregate health status by running component checks.
// GetOverallStatus 通过运行组件检查获取当前聚合健康状态。
func (m *defaultManager) GetOverallStatus(ctx context.Context) *OverallHealthStatus {
	utils.GetLogger().Println("Running component health checks...")

	// Run all registered component checks
	// 运行所有已注册的组件检查
	results := m.runComponentChecks(ctx)

	// Determine overall status based on results
	// 根据结果确定总体状态
	overallStatus := m.determineOverallStatus(results)

	utils.GetLogger().Printf("Overall health status: %s", overallStatus.Status)
	return overallStatus
}

// runComponentChecks executes all registered health checkers.
// runComponentChecks 执行所有已注册的健康检查器。
func (m *defaultManager) runComponentChecks(ctx context.Context) []HealthCheckResult {
	m.mu.RLock()
	checkers := make([]HealthChecker, 0, len(m.componentCheckers))
	for _, checker := range m.componentCheckers {
		checkers = append(checkers, checker)
	}
	m.mu.RUnlock()

	results := make([]HealthCheckResult, len(checkers))
	var wg sync.WaitGroup

	for i, checker := range checkers {
		wg.Add(1)
		go func(index int, c HealthChecker) {
			defer wg.Done()
			// Run check with a limited context derived from the input context
			// 使用从输入上下文派生的有限上下文运行检查
			checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second) // Example timeout per check // 示例每个检查的超时时间
			defer cancel()
			results[index] = c.Check(checkCtx)
		}(i, checker)
	}

	wg.Wait()
	return results
}

// determineOverallStatus aggregates results from individual checks.
// determineOverallStatus 聚合单个检查的结果。
func (m *defaultManager) determineOverallStatus(results []HealthCheckResult) *OverallHealthStatus {
	overall := &OverallHealthStatus{
		Status:  "Healthy", // Start assuming healthy // 开始假定健康
		Results: results,
	}

	// Simple aggregation: if any check is Unhealthy, overall is Unhealthy.
	// If any check is Degraded, and none are Unhealthy, overall is Degraded.
	// 简单的聚合：如果任何检查为 Unhealthy，则总体为 Unhealthy。
	// 如果任何检查为 Degraded，且没有 Unhealthy 的检查，则总体为 Degraded。
	hasDegraded := false
	for _, result := range results {
		if result.Status == "Unhealthy" {
			overall.Status = "Unhealthy"
			return overall // Found a critical issue, no need to check further // 发现关键问题，无需进一步检查
		}
		if result.Status == "Degraded" {
			hasDegraded = true
		}
		// Add other states like "Unknown" if applicable
		// 如果适用，添加其他状态，例如“Unknown”
	}

	if hasDegraded {
		overall.Status = "Degraded"
	}

	return overall
}

// RunClusterChecks performs checks on the state of the Host Cluster and vclusters.
// This is more about checking the deployed infrastructure state, not the chasi-bod management plane itself.
// RunClusterChecks 在 Host Cluster 和 vcluster 的状态上执行检查。
// 这更多是关于检查已部署基础设施的状态，而不是 chasi-bod 管理平面本身。
func (m *defaultManager) RunClusterChecks(ctx context.Context, hostK8sClient kubernetes.Interface, config *model.PlatformConfig) error {
	utils.GetLogger().Println("Running cluster health checks...")

	var clusterCheckErrors []error

	// Step 1: Check Host Cluster components status
	// 步骤 1：检查 Host Cluster 组件状态
	utils.GetLogger().Println("Checking Host Cluster component statuses...")
	// Use hostK8sClient to list and check component statuses (e.g., etcd, controller-manager, scheduler)
	// 使用 hostK8sClient 列出并检查组件状态（例如，etcd, controller-manager, scheduler）
	// componentStatuses, err := hostK8sClient.CoreV1().ComponentStatuses().List(ctx, metav1.ListOptions{})
	// if err != nil {
	// 	utils.GetLogger().Printf("Error listing Host Cluster ComponentStatuses: %v", err)
	// 	clusterCheckErrors = append(clusterCheckErrors, fmt.Errorf("failed to list component statuses: %w", err))
	// } else {
	// 	for _, status := range componentStatuses.Items {
	// 		isHealthy := false
	// 		for _, condition := range status.Conditions {
	// 			if condition.Type == corev1.ComponentHealthy && condition.Status == corev1.ConditionTrue {
	// 				isHealthy = true
	// 				break
	// 			}
	// 		}
	// 		if !isHealthy {
	// 			utils.GetLogger().Printf("Component status unhealthy: %s", status.Name)
	// 			clusterCheckErrors = append(clusterCheckErrors, fmt.Errorf("host cluster component '%s' is unhealthy", status.Name))
	// 		} else {
	// 			utils.GetLogger().Printf("Component status healthy: %s", status.Name)
	// 		}
	// 	}
	// }
	utils.GetLogger().Println("Placeholder: Host Cluster component statuses checked.")

	// Step 2: Check Host Cluster node readiness
	// 步骤 2：检查 Host Cluster 节点就绪状态
	utils.GetLogger().Println("Checking Host Cluster node readiness...")
	// Use hostK8sClient to list nodes and check their Ready condition
	// 使用 hostK8sClient 列出节点并检查其 Ready 条件
	// nodeList, err := hostK8sClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	// if err != nil {
	// 	utils.GetLogger().Printf("Error listing Host Cluster Nodes: %v", err)
	// 	clusterCheckErrors = append(clusterCheckErrors, fmt.Errorf("failed to list nodes: %w", err))
	// } else {
	// 	for _, node := range nodeList.Items {
	// 		isReady := false
	// 		for _, condition := range node.Status.Conditions {
	// 			if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue {
	// 				isReady = true
	// 				break
	// 			}
	// 		}
	// 		if !isReady {
	// 			utils.GetLogger().Printf("Node not ready: %s", node.Name)
	// 			clusterCheckErrors = append(clusterCheckErrors, fmt.Errorf("host cluster node '%s' is not ready", node.Name))
	// 		} else {
	// 			utils.GetLogger().Printf("Node ready: %s", node.Name)
	// 		}
	// 	}
	// }
	utils.GetLogger().Println("Placeholder: Host Cluster node readiness checked.")

	// Step 3: Check core Host Cluster deployments/daemonsets (e.g., CNI, CSI, CoreDNS)
	// 步骤 3：检查核心 Host Cluster 部署/daemonsets（例如，CNI, CSI, CoreDNS）
	utils.GetLogger().Println("Checking core Host Cluster deployments...")
	// Identify critical deployments/daemonsets (CNI, CSI, CoreDNS, ingress controller if used etc.)
	// Identify their namespaces (often kube-system)
	// Use hostK8sClient.AppsV1().Deployments/DaemonSets().Get() and check .Status.ReadyReplicas
	// 识别关键的 deployments/daemonsets（CNI, CSI, CoreDNS, 如果使用的 ingress controller 等）
	// 识别它们的命名空间（通常是 kube-system）
	// 使用 hostK8sClient.AppsV1().Deployments/DaemonSets().Get() 并检查 .Status.ReadyReplicas
	// coreDeployments := []string{"coredns", "kubernetes-dashboard", "..."} // Example names // 示例名称
	// for _, deployName := range coreDeployments {
	// 	deploy, err := hostK8sClient.AppsV1().Deployments("kube-system").Get(ctx, deployName, metav1.GetOptions{}) // Assuming namespace kube-system // 假设命名空间 kube-system
	// 	if err != nil {
	// 		if !apierrors.IsNotFound(err) { // Only report if not just not found
	// 			utils.GetLogger().Printf("Error getting core deployment '%s': %v", deployName, err)
	// 			clusterCheckErrors = append(clusterCheckErrors, fmt.Errorf("failed to get core deployment '%s': %w", deployName, err))
	// 		}
	// 		continue // Move to next deployment
	// 	}
	// 	if deploy.Status.ReadyReplicas < deploy.Status.Replicas {
	// 		utils.GetLogger().Printf("Core deployment '%s' not ready (%d/%d replicas)", deployName, deploy.Status.ReadyReplicas, deploy.Status.Replicas)
	// 		clusterCheckErrors = append(clusterCheckErrors, fmt.Errorf("core deployment '%s' is not ready", deployName))
	// 	} else {
	// 		utils.GetLogger().Printf("Core deployment '%s' is ready.", deployName)
	// 	}
	// }
	// TODO: Add checks for DaemonSets similarly
	utils.GetLogger().Println("Placeholder: Core Host Cluster deployments checked.")

	// Step 4: Check vcluster health
	// 步骤 4：检查 vcluster 健康状况
	if len(config.VClusters) > 0 {
		utils.GetLogger().Printf("Checking health of %d vclusters...", len(config.VClusters))
		// This requires getting a client for each vcluster and checking its /healthz or status.
		// Use the vclusterManager for this.
		// 这需要为每个 vcluster 获取一个客户端，并检查其 /healthz 或状态。
		// 使用 vclusterManager 执行此操作。
		vclusterManager := vcluster.NewManager(hostK8sClient, "pkg/vcluster/chart/vcluster") // Create a vcluster manager instance // 创建一个 vcluster manager 实例
		for name := range config.VClusters {                                                   // Iterate through vcluster names from config // 遍历配置中的 vcluster 名称
			utils.GetLogger().Printf("Checking vcluster '%s' health...", name)
			// Use the vclusterManager's WaitForReady logic or a dedicated check
			// Use a non-blocking check here, not PollUntilContextTimeout
			// 使用 vclusterManager 的 WaitForReady 逻辑或专用检查
			// 在此处使用非阻塞检查，而不是 PollUntilContextTimeout
			vClient, getClientErr := vclusterManager.GetVClusterClient(ctx, name)
			if getClientErr != nil {
				utils.GetLogger().Printf("Error getting client for vcluster '%s' for health check: %v", name, getClientErr)
				clusterCheckErrors = append(clusterCheckErrors, fmt.Errorf("failed to get client for vcluster '%s': %w", name, getClientErr))
				continue // Move to next vcluster
			}
			// Perform a basic check on the vcluster API server
			// 对 vcluster API 服务器执行基本检查
			_, healthzErr := vClient.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
			if healthzErr != nil {
				utils.GetLogger().Printf("Vcluster '%s' API server check failed: %v", name, healthzErr)
				clusterCheckErrors = append(clusterCheckErrors, fmt.Errorf("vcluster '%s' API server is not healthy: %w", name, healthzErr))
			} else {
				utils.GetLogger().Printf("Vcluster '%s' API server is healthy.", name)
			}
			// TODO: Add more specific vcluster health checks if needed
			// TODO: 如果需要，添加更具体的 vcluster 健康检查
		}
		utils.GetLogger().Println("Vcluster health checks completed.")
	} else {
		utils.GetLogger().Println("No vclusters defined, skipping vcluster health checks.")
	}

	// Aggregate results and return an error if critical issues are found
	// 聚合结果，如果发现关键问题则返回错误
	if len(clusterCheckErrors) > 0 {
		errorMessage := fmt.Sprintf("Cluster health checks failed. Found %d issues:\n", len(clusterCheckErrors))
		for _, err := range clusterCheckErrors {
			errorMessage += fmt.Sprintf("- %v\n", err)
		}
		return errors.New(errors.ErrTypeSystem, errorMessage)
	}

	utils.GetLogger().Println("Cluster health checks completed successfully.")
	return nil
}

// TODO: Implement a HealthChecker for the Logger utility
// TODO: 实现 Logger 工具的 HealthChecker
// type LoggerHealthChecker struct{}
// func (c *LoggerHealthChecker) Check(ctx context.Context) HealthCheckResult { ... }
// func (c *LoggerHealthChecker) Name() string { return "logger" }
// In the main application, register this checker: manager.RegisterComponentCheck("logger", &LoggerHealthChecker{})

// TODO: Implement other HealthChecker implementations for key platform components (e.g., Config Loader, SSH client pool status)
// TODO: 实现其他关键平台组件的 HealthChecker 实现（例如，Config Loader, SSH 客户端池状态）
