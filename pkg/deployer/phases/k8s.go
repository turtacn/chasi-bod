// Package phases defines the individual steps involved in deploying the chasi-bod platform.
// 包 phases 定义了部署 chasi-bod 平台的各个步骤。
package phases

import (
	"context"
	"time"

	//"fmt"
	"strings"

	"github.com/turtacn/chasi-bod/common/errors"
	"github.com/turtacn/chasi-bod/common/types/enum"
	"github.com/turtacn/chasi-bod/common/utils" // Assuming logger is here // 假设日志记录器在这里
	"github.com/turtacn/chasi-bod/pkg/config/model"
	// Assuming you have an SSH client utility
	// 假设您有一个 SSH 客户端工具
	// "github.com/turtacn/chasi-bod/pkg/sshutil"
	// Placeholder for Kubernetes client-go, needed for CNI/CSI deployment and node status check
	// Kubernetes client-go 的占位符，CNI/CSI 部署和节点状态检查需要它
	// "k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/rest" // Needed to build K8s client config // 需要它来构建 K8s 客户端配置
	// "k8s.io/client-go/tools/clientcmd" // Needed to load kubeconfig // 需要它来加载 kubeconfig
	// corev1 "k8s.io/api/core/v1" // Needed for Node object // 需要它来获取 Node 对象
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1" // Needed for ListOptions // 需要它来获取 ListOptions
	// "k8s.io/apimachinery/pkg/util/wait" // Needed for polling // 需要它来进行轮询
)

// K8sInstallPhase defines the interface for installing and initializing the Host Kubernetes cluster.
// K8sInstallPhase 定义了安装和初始化 Host Kubernetes 集群的接口。
// This phase uses kubeadm to initialize the control plane and join worker nodes.
// 此阶段使用 kubeadm 初始化控制平面并加入工作节点。
type K8sInstallPhase interface {
	// Run executes the Kubernetes installation phase.
	// Run 执行 Kubernetes 安装阶段。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// nodes: The configurations for all nodes in the cluster. / 集群中所有节点的配置。
	// clusterCfg: The overall cluster configuration. / 整体集群配置。
	// Returns an error if the phase fails.
	// 如果阶段失败则返回错误。
	Run(ctx context.Context, nodes []model.NodeConfig, clusterCfg *model.ClusterConfig) error

	// Add methods for joining/removing nodes if needed separate from full deploy
	// 如果需要与完整部署分开，则添加用于加入/移除节点的方法
	// JoinNode(ctx context.Context, nodeCfg *model.NodeConfig, clusterCfg *model.ClusterConfig, joinCommand string) error
	// RemoveNode(ctx context.Context, nodeCfg *model.NodeConfig, hostK8sClient kubernetes.Interface) error // This is also in deployer.Deployer
}

// NewK8sInstallPhase creates a new K8sInstallPhase instance.
// NewK8sInstallPhase 创建一个新的 K8sInstallPhase 实例。
// Returns a K8sInstallPhase implementation.
// 返回 K8sInstallPhase 实现。
func NewK8sInstallPhase() K8sInstallPhase {
	return &defaultK8sInstallPhase{}
}

// defaultK8sInstallPhase is a default implementation of the K8sInstallPhase.
// defaultK8sInstallPhase 是 K8sInstallPhase 的默认实现。
type defaultK8sInstallPhase struct{}

// Run executes the Kubernetes installation phase.
// Run 执行 Kubernetes 安装阶段。
func (p *defaultK8sInstallPhase) Run(ctx context.Context, nodes []model.NodeConfig, clusterCfg *model.ClusterConfig) error {
	utils.GetLogger().Println("Running K8sInstallPhase...")

	// Identify master and worker nodes
	// 识别主节点和工作节点
	var masterNodes, workerNodes []model.NodeConfig
	for _, node := range nodes {
		isMaster := false
		for _, role := range node.Roles {
			if role == enum.RoleMaster {
				isMaster = true
				break
			}
		}
		if isMaster {
			masterNodes = append(masterNodes, node)
		} else {
			workerNodes = append(workerNodes, node)
		}
	}

	if len(masterNodes) == 0 {
		return errors.New(errors.ErrTypeValidation, "no master nodes specified in cluster configuration")
	}

	// Step 1: Initialize Kubernetes control plane on master nodes (using kubeadm)
	// 步骤 1：在主节点上初始化 Kubernetes 控制平面（使用 kubeadm）
	// This involves running `kubeadm init` on the first master, and `kubeadm join` on others if HA
	// 这涉及在第一个主节点上运行 `kubeadm init`，如果是高可用则在其他节点上运行 `kubeadm join`
	utils.GetLogger().Printf("Initializing Kubernetes control plane on master nodes: %s", formatNodeAddresses(masterNodes))

	// TODO: Implement kubeadm init/join logic via SSH
	// TODO: 通过 SSH 实现 kubeadm init/join 逻辑
	// You'll need to:
	// - Generate kubeadm config file (YAML) based on clusterCfg
	// - Copy config file to master node
	// - Execute `sudo kubeadm init --config=<config-file>` on the first master via SSH
	// - Wait for init to complete
	// - Extract join command and token from the output
	// - For HA masters, execute `sudo kubeadm join ...` with the --control-plane flag via SSH
	// - Copy kubeconfig from the first master to a known location (e.g., /etc/kubernetes/admin.conf) on the master, and potentially download it locally for the deployer to use.

	utils.GetLogger().Println("Placeholder: Running kubeadm init on master nodes...")
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Minute): // Simulate kubeadm init time
		utils.GetLogger().Println("Placeholder: Kubeadm init completed on master nodes.")
	}

	// Placeholder for join command and token - should be extracted from kubeadm init output
	// join 命令和 token 的占位符 - 应从 kubeadm init 输出中提取
	joinCommandPlaceholder := "kubeadm join <ControlPlaneEndpoint>:<ControlPlanePort> --token <token> --discovery-token-ca-cert-hash sha256:<hash>"
	if len(masterNodes) > 1 {
		joinCommandPlaceholder = "kubeadm join <ControlPlaneEndpoint>:<ControlPlanePort> --token <token> --discovery-token-ca-cert-hash sha256:<hash> --control-plane --certificate-key <key>"
	}
	utils.GetLogger().Printf("Placeholder: Extracted join command: %s", joinCommandPlaceholder)

	// Step 2: Join worker nodes to the cluster (using kubeadm)
	// 步骤 2：将工作节点加入集群（使用 kubeadm）
	if len(workerNodes) > 0 {
		utils.GetLogger().Printf("Joining worker nodes to the cluster: %s", formatNodeAddresses(workerNodes))
		// TODO: Implement kubeadm join logic via SSH on worker nodes
		// TODO: 通过 SSH 在工作节点上实现 kubeadm join 逻辑
		// For each worker node:
		// - Execute `sudo <joinCommandPlaceholder>`
		// - Wait for node to join and become Ready (can use K8s client after deploying CNI/CSI)
		utils.GetLogger().Println("Placeholder: Running kubeadm join on worker nodes...")
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(len(workerNodes) * 2 * time.Minute): // Simulate join time per node
			utils.GetLogger().Println("Placeholder: Kubeadm join completed on worker nodes.")
		}
	} else {
		utils.GetLogger().Println("No worker nodes to join.")
	}

	// --- Host Cluster Control Plane is now UP, but potentially not fully functional without CNI/CSI ---

	// TODO: Get Host K8s client (requires kubeconfig which should be available after init)
	// TODO: 获取 Host K8s 客户端（需要 kubeconfig，在 init 后应该可用）
	// hostK8sClient, err := getHostK8sClientForInstallPhase(ctx, masterNodes[0]) // Assuming a helper to get client after init
	// if err != nil { return errors.NewWithCause(...) }
	var hostK8sClient interface{} // Placeholder
	if hostK8sClient == nil {

	}

	// Step 3: Deploy CNI (usually applied after cluster is up and client is available)
	// 步骤 3：部署 CNI（通常在集群启动并客户端可用后应用）
	utils.GetLogger().Printf("Deploying CNI plugin %s...", clusterCfg.Network.Plugin)
	// TODO: Implement CNI deployment (e.g., apply YAML manifest via kubectl or K8s client)
	// This might involve:
	// - Getting the K8s client for the Host Cluster (needs kubeconfig from master)
	// - Loading the CNI manifest (e.g., read from a file in the built image or download)
	// - Applying the manifest using the hostK8sClient
	utils.GetLogger().Println("Placeholder: CNI plugin deployed.")

	// Step 4: Deploy CSI (usually applied after cluster is up and client is available)
	// 步骤 4：部署 CSI（通常在集群启动并客户端可用后应用）
	if clusterCfg.Storage.DefaultStorageClass != "" || len(clusterCfg.Storage.StorageClasses) > 0 {
		utils.GetLogger().Printf("Deploying CSI driver for storage classes...")
		// TODO: Implement CSI deployment (similar to CNI) using hostK8sClient
		// This depends on the chosen CSI driver.
		// TODO: 使用 hostK8sClient 实现 CSI 部署（类似于 CNI）
		// 这取决于所选的 CSI 驱动。
		utils.GetLogger().Println("Placeholder: CSI driver deployed.")
	} else {
		utils.GetLogger().Println("No storage classes defined, skipping CSI deployment.")
	}

	// Step 5: Wait for all nodes to be Ready (including masters and workers, and CNI/CSI components might need to be up)
	// 步骤 5：等待所有节点就绪（包括主节点和工作节点，以及 CNI/CSI 组件可能需要启动）
	utils.GetLogger().Println("Waiting for all nodes to become Ready...")
	// TODO: Implement waiting logic using K8s client
	// Poll the status of Nodes in the Host Cluster until all desired nodes are in the Ready state.
	// 轮询 Host Cluster 中节点的状态，直到所有期望的节点都处于 Ready 状态。
	// Example:
	// err = wait.PollUntilContextTimeout(ctx, 5*time.Second, 10*time.Minute, func(ctx context.Context) (bool, error) {
	// 	nodeList, listErr := hostK8sClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	// 	if listErr != nil { return false, listErr }
	// 	readyNodes := 0
	// 	expectedNodes := len(nodes) // Total number of nodes in the cluster config
	// 	for _, node := range nodeList.Items {
	// 		for _, condition := range node.Status.Conditions {
	// 			if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue {
	// 				readyNodes++
	// 				break
	// 			}
	// 		}
	// 	}
	// 	utils.GetLogger().Printf("Waiting for nodes: %d/%d ready", readyNodes, expectedNodes)
	// 	return readyNodes == expectedNodes, nil
	// })
	// if err != nil { return errors.NewWithCause(...) }
	utils.GetLogger().Println("Placeholder: All nodes are Ready.")

	utils.GetLogger().Println("K8sInstallPhase completed successfully.")
	return nil
}

// formatNodeAddresses is a helper to format a list of node addresses for logging.
// formatNodeAddresses 是一个辅助函数，用于格式化节点地址列表以便记录日志。
func formatNodeAddresses(nodes []model.NodeConfig) string {
	addresses := make([]string, len(nodes))
	for i, node := range nodes {
		addresses[i] = node.Address
	}
	return strings.Join(addresses, ", ")
}

// TODO: Implement helper functions for remote command execution, file transfer, etc. using sshutil
// TODO: 使用 sshutil 实现远程命令执行、文件传输等辅助函数
// TODO: Implement helper function to generate kubeadm config YAML
// TODO: 实现生成 kubeadm 配置 YAML 的辅助函数
// TODO: Implement helper function to get K8s client for Host Cluster after init (needs kubeconfig from master)
// TODO: 实现获取 Host 集群 K8s 客户端的辅助函数（在 init 后需要从主节点获取 kubeconfig）
// func getHostK8sClientForInstallPhase(ctx context.Context, masterNodeCfg *model.NodeConfig) (kubernetes.Interface, error) { ... }
// TODO: Implement helper function to drain a node using K8s client
// TODO: 实现使用 K8s 客户端排空节点的辅助函数
// TODO: Implement helper function to remove a node from cluster using K8s client
// TODO: 实现使用 K8s 客户端从集群移除节点的辅助函数
// TODO: Implement helper function to wait for a node to be Ready using K8s client
// TODO: 实现使用 K8s 客户端等待节点就绪的辅助函数
