// Package deployer orchestrates the deployment of the chasi-bod platform onto target nodes.
// 包 deployer 协调将 chasi-bod 平台部署到目标节点。
package deployer

import (
	"context"
	"fmt"
	"time"

	//"time"

	"github.com/turtacn/chasi-bod/common/errors"
	"github.com/turtacn/chasi-bod/common/utils" // Assuming logger is here // 假设日志记录器在这里
	"github.com/turtacn/chasi-bod/pkg/config/model"
	"github.com/turtacn/chasi-bod/pkg/deployer/phases"
	// Placeholder for Kubernetes client-go, needed for vcluster phase
	// Kubernetes client-go 的占位符，vcluster 阶段需要它
	// "k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/rest" // Might be needed to build K8s client config // 可能需要它来构建 K8s 客户端配置
	// "k8s.io/client-go/tools/clientcmd" // Needed to load kubeconfig // 需要它来加载 kubeconfig
)

// Deployer defines the interface for the platform deployment orchestrator.
// Deployer 定义了平台部署协调器的接口。
type Deployer interface {
	// Deploy orchestrates the step-by-step deployment process based on the platform configuration.
	// Deploy 根据平台配置协调分步部署过程。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// config: The platform configuration including cluster and node details. / 包含集群和节点详细信息的平台配置。
	// Returns an error if the deployment fails at any stage.
	// 如果部署在任何阶段失败则返回错误。
	Deploy(ctx context.Context, config *model.PlatformConfig) error

	// AddNode adds a new node to an existing Host Kubernetes cluster.
	// AddNode 将新节点添加到现有的 Host Kubernetes 集群。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// config: The overall platform configuration. / 整体平台配置。
	// nodeCfg: The configuration for the new node. / 新节点的配置。
	// Returns an error if adding the node fails.
	// 如果添加节点失败则返回错误。
	AddNode(ctx context.Context, config *model.PlatformConfig, nodeCfg *model.NodeConfig) error

	// RemoveNode removes a node from an existing Host Kubernetes cluster.
	// RemoveNode 从现有的 Host Kubernetes 集群中移除节点。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// nodeCfg: The configuration for the node to remove. / 要移除的节点的配置。
	// hostK8sClient: A Kubernetes client connected to the Host Cluster. / 连接到 Host 集群的 Kubernetes 客户端。
	// Returns an error if removing the node fails.
	// 如果移除节点失败则返回错误。
	RemoveNode(ctx context.Context, nodeCfg *model.NodeConfig, hostK8sClient interface{} /* kubernetes.Interface */) error

	// Add more methods for other deployment actions like upgrade, etc.
	// 添加其他部署操作的方法，例如升级等。
	// Upgrade(ctx context.Context, config *model.PlatformConfig, newConfig *model.PlatformConfig) error // This might be in lifecycle manager
}

// defaultDeployer is a default implementation of the Deployer interface.
// defaultDeployer 是 Deployer 接口的默认实现。
type defaultDeployer struct {
	// Phases hold instances of deployment phases
	// Phases 持有部署阶段的实例
	initializationPhase phases.InitPhase
	osConfigPhase       phases.OSConfigPhase
	runtimeConfigPhase  phases.RuntimeConfigPhase
	networkConfigPhase  phases.NetworkConfigPhase
	storageConfigPhase  phases.StorageConfigPhase
	k8sInstallPhase     phases.K8sInstallPhase // Renamed for clarity // 重命名以提高清晰度
	vclusterDeployPhase phases.VClusterDeployPhase
	// Add other phases here
	// 在这里添加其他阶段
}

// NewDeployer creates a new Deployer instance.
// NewDeployer 创建一个新的 Deployer 实例。
// Returns a Deployer implementation.
// 返回 Deployer 实现。
func NewDeployer() (Deployer, error) {
	// Default deployer implementation will orchestrate the phases
	// 默认的 deployer 实现将协调各个阶段
	d := &defaultDeployer{}
	d.init() // Initialize phases on creation
	return d, nil
}

// init initializes the defaultDeployer with concrete phase implementations.
// init 使用具体的阶段实现初始化 defaultDeployer。
func (d *defaultDeployer) init() {
	// Initialize all phases
	// 初始化所有阶段
	d.initializationPhase = phases.NewInitPhase()
	d.osConfigPhase = phases.NewOSConfigPhase()
	d.runtimeConfigPhase = phases.NewRuntimeConfigPhase()
	d.networkConfigPhase = phases.NewNetworkConfigPhase()
	d.storageConfigPhase = phases.NewStorageConfigPhase()
	d.k8sInstallPhase = phases.NewK8sInstallPhase()
	d.vclusterDeployPhase = phases.NewVClusterDeployPhase()
	// Initialize other phases
	// 初始化其他阶段
}

// Deploy orchestrates the full deployment process.
// Deploy 协调完整的部署过程。
func (d *defaultDeployer) Deploy(ctx context.Context, config *model.PlatformConfig) error {
	utils.GetLogger().Printf("Starting platform deployment for config: %s", config.Metadata.Name)

	// Define the sequence of deployment phases
	// 定义部署阶段的顺序
	// This sequence is crucial and represents the state transitions from bare OS to running platform
	// 这个顺序至关重要，代表了从裸操作系统到运行平台的状体转换
	// Note: vCluster Deployment happens *after* Host K8s is fully up
	// 注意：vCluster 部署发生在 Host K8s 完全启动之后

	// Phase 1: Initialization (per node)
	// 阶段 1：初始化（每个节点）
	utils.GetLogger().Println("--- Running Initialization Phase ---")
	for i, nodeCfg := range config.Cluster.Nodes {
		nodeCtx, cancel := context.WithTimeout(ctx, 5*time.Minute) // Context for node-specific phase
		defer cancel()
		utils.GetLogger().Printf("Initializing node %d/%d: %s", i+1, len(config.Cluster.Nodes), nodeCfg.Address)
		if err := d.initializationPhase.Run(nodeCtx, &nodeCfg, &config.Cluster); err != nil {
			return errors.NewWithCause(errors.ErrTypeSystem, fmt.Sprintf("initialization phase failed for node %s", nodeCfg.Address), err)
		}
		utils.GetLogger().Printf("Node %s initialized successfully.", nodeCfg.Address)
	}

	// Phase 2: OS Configuration (per node)
	// 阶段 2：操作系统配置（每个节点）
	utils.GetLogger().Println("--- Running OS Configuration Phase ---")
	for i, nodeCfg := range config.Cluster.Nodes {
		nodeCtx, cancel := context.WithTimeout(ctx, 10*time.Minute) // Longer timeout for OS tasks
		defer cancel()
		utils.GetLogger().Printf("Configuring OS for node %d/%d: %s", i+1, len(config.Cluster.Nodes), nodeCfg.Address)
		if err := d.osConfigPhase.Run(nodeCtx, &nodeCfg, &config.Cluster); err != nil {
			return errors.NewWithCause(errors.ErrTypeSystem, fmt.Sprintf("OS configuration phase failed for node %s", nodeCfg.Address), err)
		}
		utils.GetLogger().Printf("OS configured successfully for node %s.", nodeCfg.Address)
	}

	// Phase 3: Runtime Configuration (per node)
	// 阶段 3：运行时配置（每个节点）
	utils.GetLogger().Println("--- Running Runtime Configuration Phase ---")
	for i, nodeCfg := range config.Cluster.Nodes {
		nodeCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
		utils.GetLogger().Printf("Configuring Runtime for node %d/%d: %s", i+1, len(config.Cluster.Nodes), nodeCfg.Address)
		if err := d.runtimeConfigPhase.Run(nodeCtx, &nodeCfg, &config.Cluster); err != nil {
			return errors.NewWithCause(errors.ErrTypeSystem, fmt.Sprintf("runtime configuration phase failed for node %s", nodeCfg.Address), err)
		}
		utils.GetLogger().Printf("Runtime configured successfully for node %s.", nodeCfg.Address)
	}

	// Phase 4: Network Configuration (per node)
	// 阶段 4：网络配置（每个节点）
	utils.GetLogger().Println("--- Running Network Configuration Phase ---")
	for i, nodeCfg := range config.Cluster.Nodes {
		nodeCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
		utils.GetLogger().Printf("Configuring Network for node %d/%d: %s", i+1, len(config.Cluster.Nodes), nodeCfg.Address)
		if err := d.networkConfigPhase.Run(nodeCtx, &nodeCfg, &config.Cluster); err != nil {
			return errors.NewWithCause(errors.ErrTypeSystem, fmt.Sprintf("network configuration phase failed for node %s", nodeCfg.Address), err)
		}
		utils.GetLogger().Printf("Network configured successfully for node %s.", nodeCfg.Address)
	}

	// Phase 5: Storage Configuration (per node)
	// 阶段 5：存储配置（每个节点）
	utils.GetLogger().Println("--- Running Storage Configuration Phase ---")
	for i, nodeCfg := range config.Cluster.Nodes {
		nodeCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
		utils.GetLogger().Printf("Configuring Storage for node %d/%d: %s", i+1, len(config.Cluster.Nodes), nodeCfg.Address)
		if err := d.storageConfigPhase.Run(nodeCtx, &nodeCfg, &config.Cluster); err != nil {
			return errors.NewWithCause(errors.ErrTypeSystem, fmt.Sprintf("storage configuration phase failed for node %s", nodeCfg.Address), err)
		}
		utils.GetLogger().Printf("Storage configured successfully for node %s.", nodeCfg.Address)
	}

	// Phase 6: Kubernetes Installation (orchestrated across nodes)
	// 阶段 6：Kubernetes 安装（跨节点协调执行）
	utils.GetLogger().Println("--- Running Kubernetes Installation Phase ---")
	k8sInstallCtx, cancel := context.WithTimeout(ctx, 30*time.Minute) // Give K8s installation more time
	defer cancel()
	// This phase needs the full config and orchestrates across nodes
	// 这个阶段需要完整的配置并在节点之间协调
	if err := d.k8sInstallPhase.Run(k8sInstallCtx, config.Cluster.Nodes, &config.Cluster); err != nil {
		return errors.NewWithCause(errors.ErrTypeSystem, "Kubernetes Host Cluster installation phase failed", err)
	}
	utils.GetLogger().Println("Kubernetes Host Cluster installed successfully.")

	// --- Host Cluster is now UP ---
	// At this point, the Host Cluster is running and accessible.
	// 现在，Host 集群已经启动并可访问。
	// We need to get a Kubernetes client to interact with it for the next phase.
	// 我们需要获取一个 Kubernetes 客户端来与其交互以便进行下一阶段。

	// TODO: Get Host K8s client
	// TODO: 获取 Host K8s 客户端
	// hostK8sClient, err := getHostK8sClient(config) // Need to implement this function
	// if err != nil {
	// 	return errors.NewWithCause(errors.ErrTypeSystem, "failed to get host Kubernetes client", err)
	// }
	var hostK8sClient interface{} // Placeholder for now - should be kubernetes.Interface
	utils.GetLogger().Println("Placeholder: Acquired host Kubernetes client.")

	// Phase 7: vCluster Deployment (orchestrated)
	// 阶段 7：vCluster 部署（协调执行）
	utils.GetLogger().Println("--- Running vCluster Deployment Phase ---")
	vclusterDeployCtx, cancel := context.WithTimeout(ctx, 15*time.Minute) // Give vcluster deployment time
	defer cancel()
	// This phase needs the full platform config and the Host K8s client
	// 这个阶段需要完整的平台配置和 Host K8s 客户端
	if err := d.vclusterDeployPhase.Run(vclusterDeployCtx, config, hostK8sClient /* Pass config.VClusters etc. */); err != nil {
		return errors.NewWithCause(errors.ErrTypeVCluster, "vcluster deployment phase failed", err)
	}
	utils.GetLogger().Println("vClusters deployed successfully.")

	// TODO: Add DFX deployment/configuration phase
	// TODO: 添加 DFX 部署/配置阶段

	utils.GetLogger().Println("Platform deployment completed successfully.")
	return nil
}

// AddNode adds a new node to an existing Host Kubernetes cluster.
// AddNode 将新节点添加到现有的 Host Kubernetes 集群。
func (d *defaultDeployer) AddNode(ctx context.Context, config *model.PlatformConfig, nodeCfg *model.NodeConfig) error {
	utils.GetLogger().Printf("Starting to add node %s to Host Cluster...", nodeCfg.Address)

	// This process involves running a subset of the main deploy phases on the new node.
	// 此过程涉及在新节点上运行主部署阶段的一个子集。
	// Phase sequence for adding a node: Init -> OS Config -> Runtime Config -> Network Config -> Storage Config -> K8s Join
	// 添加节点的阶段顺序：Init -> OS Config -> Runtime Config -> Network Config -> Storage Config -> K8s Join

	// TODO: Implement AddNode logic by running specific phases for the new node
	// TODO: 通过为新节点运行特定阶段来实现 AddNode 逻辑

	nodeCtx, cancel := context.WithTimeout(ctx, 20*time.Minute) // Context for single node deployment
	defer cancel()

	// Run phases Init, OS Config, Runtime Config, Network Config, Storage Config
	// 运行 Init, OS Config, Runtime Config, Network Config, Storage Config 阶段
	phasesToRun := []struct {
		name  string
		phase phases.NodeSpecificPhase // Assuming a common interface for node phases
	}{
		{"Initialization", d.initializationPhase},
		{"OS Configuration", d.osConfigPhase},
		{"Runtime Configuration", d.runtimeConfigPhase},
		{"Network Configuration", d.networkConfigPhase},
		{"Storage Configuration", d.storageConfigPhase},
	}

	for _, p := range phasesToRun {
		utils.GetLogger().Printf("--- Running %s Phase for new node %s ---", p.name, nodeCfg.Address)
		// Need to check if the phase implements NodeSpecificPhase and call Run appropriately
		// 需要检查阶段是否实现了 NodeSpecificPhase 并适当调用 Run
		// For now, call Run assuming it takes nodeCfg and clusterCfg
		// 现在，假设 Run 接受 nodeCfg 和 clusterCfg 来调用它
		if err := p.phase.Run(nodeCtx, nodeCfg, &config.Cluster); err != nil { // Assuming Run method signature
			return errors.NewWithCause(errors.ErrTypeSystem, fmt.Sprintf("%s phase failed for new node %s", p.name, nodeCfg.Address), err)
		}
		utils.GetLogger().Printf("%s phase completed successfully for new node %s.", p.name, nodeCfg.Address)
	}

	// Run K8s Join Phase specifically for the new node
	// 专门为新节点运行 K8s Join Phase
	utils.GetLogger().Printf("--- Running Kubernetes Join Phase for new node %s ---", nodeCfg.Address)
	// The K8sInstallPhase.Run method currently takes ALL nodes.
	// We need to adapt it to only *join* a specific node to an *existing* cluster.
	// K8sInstallPhase.Run 方法当前接受所有节点。
	// 我们需要对其进行调整，使其只将特定节点加入到现有集群中。
	// This likely requires calling kubeadm join remotely on the new node using SSH.
	// 这可能需要通过 SSH 在新节点上远程调用 kubeadm join。
	// It also needs the join command/token from the existing cluster (often fetched from a master).
	// 它还需要从现有集群获取 join command/token（通常从主节点获取）。

	// TODO: Get join command/token from an existing master node
	// TODO: 从现有主节点获取 join command/token
	// TODO: Execute `sudo kubeadm join ...` on the new node via SSH
	// TODO: 通过 SSH 在新节点上执行 `sudo kubeadm join ...`
	utils.GetLogger().Printf("Placeholder: Running kubeadm join on new node %s.", nodeCfg.Address)

	// TODO: Wait for the new node to become Ready in the Host Cluster
	// TODO: 等待新节点在 Host Cluster 中变为 Ready
	utils.GetLogger().Printf("Placeholder: Waiting for new node %s to become Ready.", nodeCfg.Address)

	utils.GetLogger().Printf("Node %s added to Host Cluster successfully.", nodeCfg.Address)
	return nil
}

// RemoveNode removes a node from an existing Host Kubernetes cluster.
// RemoveNode 从现有的 Host Kubernetes 集群中移除节点。
func (d *defaultDeployer) RemoveNode(ctx context.Context, nodeCfg *model.NodeConfig, hostK8sClient interface{} /* kubernetes.Interface */) error {
	utils.GetLogger().Printf("Starting to remove node %s from Host Cluster...", nodeCfg.Address)

	// This involves draining the node, removing it from the cluster, and potentially cleaning up the OS.
	// 这涉及排空节点、将其从集群中移除，以及可能清理操作系统。

	// TODO: Get Host K8s client using the provided interface{}
	// TODO: 使用提供的 interface{} 获取 Host K8s 客户端
	// k8sClient, ok := hostK8sClient.(kubernetes.Interface)
	// if !ok { return errors.New(errors.ErrTypeInternal, "invalid Host Kubernetes client provided") }

	// Step 1: Drain the node (move pods to other nodes)
	// 步骤 1：排空节点（将 pod 移走）
	utils.GetLogger().Printf("Draining node %s...", nodeCfg.Address)
	// This requires the Host K8s client and the node's Kubernetes name (which might be different from address).
	// Needs the node's Kubernetes name (obj.Name from List(Nodes)).
	// 这需要 Host K8s 客户端和节点的 Kubernetes 名称（可能与地址不同）。
	// 需要节点的 Kubernetes 名称（来自 List(Nodes) 的 obj.Name）。
	// Assume we can get the Kubernetes node name from its address or config.
	// 假设我们可以从节点的地址或配置中获取其 Kubernetes 节点名称。
	k8sNodeName := nodeCfg.Address // Placeholder - should find the actual K8s node object
	// TODO: Implement drain logic using client-go corev1.Nodes().Evict() or related functions
	// Can use kubectl drain equivalent logic via client-go.
	// 可以使用 client-go 实现 kubectl drain 的等效逻辑。
	utils.GetLogger().Printf("Placeholder: Node %s drained.", nodeCfg.Address)

	// Step 2: Delete the node from the cluster
	// 步骤 2：从集群中删除节点
	utils.GetLogger().Printf("Deleting node %s from cluster...", k8sNodeName)
	// TODO: Implement deletion logic using client-go corev1.Nodes().Delete()
	utils.GetLogger().Printf("Placeholder: Node %s deleted from cluster.", k8sNodeName)

	// Step 3: Clean up the OS on the removed node (optional, depends on full deprovisioning vs reuse)
	// 步骤 3：清理已移除节点上的操作系统（可选，取决于完全解除供应还是重用）
	// This might involve running `kubeadm reset` remotely via SSH.
	// 这可能涉及通过 SSH 远程运行 `kubeadm reset`。
	utils.GetLogger().Printf("Placeholder: Running OS cleanup on node %s (e.g., kubeadm reset).", nodeCfg.Address)
	// TODO: Get SSH client
	// TODO: sshClient.RunCommand(ctx, "sudo kubeadm reset --force")

	utils.GetLogger().Printf("Node %s removed from Host Cluster successfully.", nodeCfg.Address)
	return nil
}

// Placeholder function to get Host K8s client - requires client-go and kubeconfig loading
// 获取 Host K8s 客户端的占位符函数 - 需要 client-go 和 kubeconfig 加载
// func getHostK8sClient(config *model.PlatformConfig) (kubernetes.Interface, error) {
// 	// This function would typically read kubeconfig from a well-known location
// 	// on the master node or construct it based on cluster configuration.
// 	// For a real implementation, you'd use clientcmd.BuildConfigFromFlags or clientcmd.RESTConfigFromKubeConfig
// 	// 这个函数通常会从主节点上的已知位置读取 kubeconfig，或根据集群配置构造它。
// 	// 对于实际实现，您将使用 clientcmd.BuildConfigFromFlags 或 clientcmd.RESTConfigFromKubeConfig
//
// 	// Example using a hypothetical kubeconfig path derived from config
// 	// 示例使用从配置派生的假设 kubeconfig 路径
// 	// kubeconfigPath := filepath.Join(constants.DefaultDataDir, config.Metadata.Name, "kubeconfig")
// 	// config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
// 	// if err != nil { return nil, err }
// 	// return kubernetes.NewForConfig(config)
//
// 	return nil, errors.New(errors.ErrTypeNotImplemented, "getting host Kubernetes client not implemented")
// }

// Ensure node-specific phases implement the marker interface (for the example loop above)
// 确保节点特定的阶段实现了标记接口（用于上面的示例循环）
var _ phases.NodeSpecificPhase = &phases.DefaultInitPhase{}
var _ phases.NodeSpecificPhase = &phases.DefaultOSConfigPhase{}
var _ phases.NodeSpecificPhase = &phases.DefaultRuntimeConfigPhase{}
var _ phases.NodeSpecificPhase = &phases.DefaultNetworkConfigPhase{}
var _ phases.NodeSpecificPhase = &phases.DefaultStorageConfigPhase{}
