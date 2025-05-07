// Package phases defines the individual steps involved in deploying the chasi-bod platform.
// 包 phases 定义了部署 chasi-bod 平台的各个步骤。
package phases

import (
	"context"
	"github.com/turtacn/chasi-bod/common/utils" // Assuming logger is here // 假设日志记录器在这里
	"github.com/turtacn/chasi-bod/pkg/config/model"
	"strings"
	// Assuming you have an SSH client utility
	// 假设您有一个 SSH 客户端工具
	// "github.com/turtacn/chasi-bod/pkg/sshutil"
)

// NetworkConfigPhase defines the interface for configuring network settings on target nodes.
// NetworkConfigPhase 定义了在目标节点上配置网络设置的接口。
// This includes interfaces, IP addresses, routes, and potentially firewall rules.
// 这包括接口、IP 地址、路由以及可能的防火墙规则。
type NetworkConfigPhase interface {
	// Run executes the network configuration phase for a single node.
	// Run 为单个节点执行网络配置阶段。
	// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
	// nodeCfg: The configuration for the specific node. / 特定节点的配置。
	// clusterCfg: The overall cluster configuration including network details. / 整体集群配置，包括网络详细信息。
	// Returns an error if the phase fails for the node.
	// 如果阶段在节点上失败则返回错误。
	Run(ctx context.Context, nodeCfg *model.NodeConfig, clusterCfg *model.ClusterConfig) error
}

// NewNetworkConfigPhase creates a new NetworkConfigPhase instance.
// NewNetworkConfigPhase 创建一个新的 NetworkConfigPhase 实例。
// Returns a NetworkConfigPhase implementation.
// 返回 NetworkConfigPhase 实现。
func NewNetworkConfigPhase() NetworkConfigPhase {
	return &defaultNetworkConfigPhase{}
}

// defaultNetworkConfigPhase is a default implementation of the NetworkConfigPhase.
// defaultNetworkConfigPhase 是 NetworkConfigPhase 的默认实现。
type defaultNetworkConfigPhase struct{}

// Run executes the network configuration phase.
// Run 执行网络配置阶段。
func (p *defaultNetworkConfigPhase) Run(ctx context.Context, nodeCfg *model.NodeConfig, clusterCfg *model.ClusterConfig) error {
	utils.GetLogger().Printf("Running NetworkConfigPhase for node %s", nodeCfg.Address)

	// TODO: Get SSH client for the node
	// TODO: 获取节点的 SSH 客户端
	// sshClient, err := sshutil.NewClient(...)
	// if err != nil { return errors.NewWithCause(...) }
	// defer sshClient.Close()

	// Step 1: Configure network interfaces (IP addresses, netmask, gateway)
	// 步骤 1：配置网络接口（IP 地址、子网掩码、网关）
	if len(clusterCfg.Network.Interfaces) > 0 {
		utils.GetLogger().Printf("Configuring network interfaces on %s...", nodeCfg.Address)
		// This is highly OS-dependent (netplan, networkd, ifconfig, ip command)
		// This involves generating configuration files and applying them remotely.
		// 这高度依赖于操作系统（netplan, networkd, ifconfig, ip 命令）
		// 这涉及生成配置文件并远程应用它们。
		// TODO: Implement remote interface configuration via SSH based on OS
		// TODO: 根据操作系统通过 SSH 实现远程接口配置
		for _, ifaceCfg := range clusterCfg.Network.Interfaces {
			utils.GetLogger().Printf("Configuring interface %s on %s with IPs %s, Gateway %s",
				ifaceCfg.Name, nodeCfg.Address, strings.Join(ifaceCfg.IPAddrs, ", "), ifaceCfg.Gateway)
			// Example commands (Debian/Ubuntu):
			// - Generate netplan YAML locally
			// - Copy YAML to /etc/netplan/ via SSH/SFTP
			// - sshClient.RunCommand(ctx, "sudo netplan apply")
			// Example commands (CentOS/RHEL):
			// - Generate ifcfg file locally
			// - Copy to /etc/sysconfig/network-scripts/ via SSH/SFTP
			// - sshClient.RunCommand(ctx, "sudo systemctl restart network")
			utils.GetLogger().Printf("Placeholder: Interface %s configured on %s.", ifaceCfg.Name, nodeCfg.Address)
		}
		utils.GetLogger().Printf("Network interfaces configured successfully on %s.", nodeCfg.Address)
	} else {
		utils.GetLogger().Printf("No network interfaces defined to configure on %s.", nodeCfg.Address)
	}

	// Step 2: Configure routes (if necessary, beyond default gateway)
	// 步骤 2：配置路由（如果需要，超出默认网关）
	utils.GetLogger().Printf("Placeholder: Configuring routes on %s.", nodeCfg.Address)
	// TODO: Implement remote route configuration via SSH using `ip route add` or similar
	// TODO: 通过 SSH 使用 `ip route add` 或类似命令实现远程路由配置

	// Step 3: Configure firewall rules (e.g., firewalld, iptables, nftables)
	// 步骤 3：配置防火墙规则（例如，firewalld, iptables, nftables）
	// Need to open ports for K8s components, SSH, NodePort range etc.
	// This is also highly OS-dependent.
	// 需要打开 K8s 组件、SSH、NodePort 范围等的端口
	// 这也高度依赖于操作系统。
	utils.GetLogger().Printf("Configuring firewall on %s...", nodeCfg.Address)
	// TODO: Implement remote firewall configuration via SSH based on OS and required ports
	// TODO: 根据操作系统和所需端口通过 SSH 实现远程防火墙配置
	utils.GetLogger().Printf("Placeholder: Firewall configured on %s.", nodeCfg.Address)

	// Step 4: Configure DNS resolution (e.g., /etc/resolv.conf)
	// 步骤 4：配置 DNS 解析（例如，/etc/resolv.conf）
	utils.GetLogger().Printf("Configuring DNS resolution on %s...", nodeCfg.Address)
	// This involves creating or updating /etc/resolv.conf.
	// 这涉及创建或更新 /etc/resolv.conf。
	// TODO: Implement remote DNS configuration via SSH
	utils.GetLogger().Printf("Placeholder: DNS resolution configured on %s.", nodeCfg.Address)

	// Step 5: Wait for network to be ready and accessible (optional but recommended)
	// 步骤 5：等待网络就绪并可访问（可选但推荐）
	utils.GetLogger().Printf("Waiting for network to be ready on %s...", nodeCfg.Address)
	// This might involve checking reachability to a known external service or another node via SSH.
	// 这可能涉及通过 SSH 检查对已知外部服务或其他节点的的可达性。
	utils.GetLogger().Printf("Placeholder: Network is ready on %s.", nodeCfg.Address)

	utils.GetLogger().Printf("NetworkConfigPhase completed successfully for node %s", nodeCfg.Address)
	return nil
}

// TODO: Implement sshutil package for remote command execution and file transfer
// TODO: 实现 sshutil 包用于远程命令执行和文件传输
// TODO: Implement helper functions to generate OS-specific network config files (netplan, ifcfg etc.)
// TODO: 实现生成操作系统特定网络配置文件的辅助函数（netplan, ifcfg 等）
