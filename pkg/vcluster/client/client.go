// Package client provides functionality to obtain Kubernetes clients for virtual clusters.
// 包 client 提供了获取虚拟集群 Kubernetes 客户端的功能。
// It abstracts the connection method (e.g., via Host Service, port-forwarding).
// 它抽象了连接方法（例如，通过 Host Service，端口转发）。
package client

import (
	"context"
	"encoding/base64"
	"fmt"
	//"time"

	"github.com/turtacn/chasi-bod/common/errors"
	"github.com/turtacn/chasi-bod/common/utils" // Assuming logger is here // 假设日志记录器在这里

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1" // Added for meta types // 添加用于元数据类型
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest" // Added for rest.Config // 添加用于 rest.Config
	// You might need to import clientcmd for loading kubeconfig if not connecting via service
	// 如果不是通过服务连接，您可能需要导入 clientcmd 以加载 kubeconfig
	// "k8s.io/client-go/tools/clientcmd"
	// Import client-go exec for port-forwarding (if using that method)
	// 导入用于端口转发的 client-go exec（如果使用此方法）
	// "k8s.io/client-go/tools/remotecommand"
	// Import client-go portforwarder (if using that method)
	// 导入 client-go portforwarder（如果使用此方法）
	// "k8s.io/client-go/tools/portforward"
	// Import kubelet client for port-forwarding (if using that method)
	// 导入用于端口转发的 kubelet 客户端（如果使用此方法）
	// "net/http" // Needed for portforwarding // 端口转发需要
	// "net/url" // Needed for portforwarding // 端口转发需要
	// "k8s.io/apimachinery/pkg/runtime/schema" // Might be needed for rest.Config // rest.Config 可能需要
	// "k8s.io/apimachinery/pkg/runtime" // Might be needed for rest.Config // rest.Config 可能需要
)

// GetVClusterClient returns a Kubernetes client configured to interact with the specified virtual cluster.
// It finds the vcluster server within the host cluster and establishes a connection.
// GetVClusterClient 返回配置用于与指定虚拟集群交互的 Kubernetes 客户端。
// 它在 host 集群中找到 vcluster 服务器并建立连接。
// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
// vclusterName: The name of the virtual cluster. / 虚拟集群的名称。
// hostK8sClient: A Kubernetes client connected to the Host Cluster. / 连接到 Host 集群的 Kubernetes 客户端。
// Returns a Kubernetes client for the virtual cluster and an error.
// 返回虚拟集群的 Kubernetes 客户端和错误。
func GetVClusterClient(ctx context.Context, vclusterName string, hostK8sClient kubernetes.Interface) (kubernetes.Interface, error) {
	utils.GetLogger().Printf("Attempting to get client for vcluster '%s'", vclusterName)

	// Determine the host namespace where the vcluster is running
	// 确定 vcluster 正在运行的 host 命名空间
	// This should ideally come from a stored state or the config if available,
	// but convention (e.g., vcluster-<name>) is common.
	// 这理想情况下应来自存储的状态或配置（如果可用），
	// 但约定（例如，vcluster-<name>）是常见的。
	hostNamespace := fmt.Sprintf("vcluster-%s", vclusterName) // TODO: Get namespace reliably from config or stored state

	// Option 1: Use the vcluster's Service in the Host Cluster (simpler, recommended if possible)
	// Vcluster typically exposes its API server via a ClusterIP Service in the host namespace.
	// We can build a rest.Config targeting this service.
	// 选项 1：使用 Host 集群中的 vcluster Service（更简单，如果可能则推荐）
	// Vcluster 通常通过 host 命名空间中的 ClusterIP Service 暴露其 API 服务器。
	// 我们可以构建一个针对此服务的 rest.Config。
	utils.GetLogger().Printf("Attempting to connect to vcluster '%s' via Host Service in namespace '%s'", vclusterName, hostNamespace)

	// The service name is usually the same as the vcluster name
	// Service 名称通常与 vcluster 名称相同
	vclusterServiceName := vclusterName // TODO: Confirm naming convention

	// Build the API server URL using the service name and host namespace
	// 使用 Service 名称和 host 命名空间构建 API 服务器 URL
	// Format: https://<service-name>.<namespace>.svc.cluster.local:<port>
	// The port is usually 443 for HTTPS
	// 格式：https://<service-name>.<namespace>.svc.cluster.local:<port>
	// 端口通常是 443 用于 HTTPS
	vclusterAPIServerURL := fmt.Sprintf("https://%s.%s.svc.cluster.local:443", vclusterServiceName, hostNamespace)
	utils.GetLogger().Printf("Vcluster API server URL: %s", vclusterAPIServerURL)

	// Need to build a rest.Config that trusts the *vcluster's* CA certificate.
	// The vcluster CA cert is typically stored as a Secret in the host namespace.
	// Secret name is usually "vc-certs-{{ .Release.Name }}" or similar.
	// 需要构建一个信任虚拟集群的 CA 证书的 rest.Config。
	// 虚拟集群的 CA 证书通常作为 Secret 存储在 host 命名空间中。
	// Secret 名称通常为 "vc-certs-{{ .Release.Name }}" 或类似名称。
	certsSecretName := fmt.Sprintf("vc-certs-%s", vclusterName) // TODO: Confirm secret name

	utils.GetLogger().Printf("Fetching vcluster CA certificate from Secret '%s' in host namespace '%s'", certsSecretName, hostNamespace)
	certsSecret, err := hostK8sClient.CoreV1().Secrets(hostNamespace).Get(ctx, certsSecretName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.NewWithCause(errors.ErrTypeVCluster, fmt.Sprintf("failed to get vcluster certs secret '%s/%s'", hostNamespace, certsSecretName), err)
	}

	caCert, ok := certsSecret.Data["ca.crt"] // TODO: Confirm key name in Secret
	if !ok || len(caCert) == 0 {
		return nil, errors.New(errors.ErrTypeVCluster, fmt.Sprintf("vcluster CA certificate not found in secret '%s/%s'", hostNamespace, certsSecretName))
	}

	// Fetch the token from the secret
	tokenSecretName := fmt.Sprintf("vc-token-%s", vclusterName)
	utils.GetLogger().Printf("Fetching vcluster token from Secret '%s' in host namespace '%s'", tokenSecretName, hostNamespace)
	tokenSecret, err := hostK8sClient.CoreV1().Secrets(hostNamespace).Get(ctx, tokenSecretName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.NewWithCause(errors.ErrTypeVCluster, fmt.Sprintf("failed to get vcluster token secret '%s/%s'", hostNamespace, tokenSecretName), err)
	}

	token, ok := tokenSecret.Data["token"]
	if !ok || len(token) == 0 {
		return nil, errors.New(errors.ErrTypeVCluster, fmt.Sprintf("vcluster token not found in secret '%s/%s'", hostNamespace, tokenSecretName))
	}

	// Now build the rest.Config
	// 现在构建 rest.Config
	vclusterConfig := &rest.Config{
		Host:        vclusterAPIServerURL,
		BearerToken: string(token),
		TLSClientConfig: rest.TLSClientConfig{
			CAData: caCert, // Trust the vcluster's CA // 信任 vcluster 的 CA
		},
	}

	// Create the vcluster client
	// 创建 vcluster 客户端
	vClient, err := kubernetes.NewForConfig(vclusterConfig)
	if err != nil {
		return nil, errors.NewWithCause(errors.ErrTypeVCluster, fmt.Sprintf("failed to create client for vcluster '%s'", vclusterName), err)
	}

	utils.GetLogger().Printf("Successfully obtained client for vcluster '%s' via Host Service.", vclusterName)
	return vClient, nil

	// Option 2: Use Port-Forwarding (More complex, less recommended for production use cases)
	// This involves finding the vcluster pod in the host namespace, setting up a port-forward,
	// and then building a rest.Config targeting the local forwarded port.
	// This requires more complex client-go logic (using portforward package).
	// 选项 2：使用端口转发（更复杂，生产用例不推荐）
	// 这涉及在 host 命名空间中查找 vcluster pod，设置端口转发，
	// 然后构建一个针对本地转发端口的 rest.Config。
	// 这需要更复杂的 client-go 逻辑（使用 portforward 包）。
	// utils.GetLogger().Printf("Attempting to get client for vcluster '%s' via Port-Forwarding...", vclusterName)
	// TODO: Implement port-forwarding logic using k8s.io/client-go/tools/portforward
	// return nil, errors.New(errors.ErrTypeNotImplemented, fmt.Sprintf("getting client for vcluster '%s' via port-forwarding not implemented yet", vclusterName))
}

// GetVClusterKubeConfig generates a kubeconfig string for the specified virtual cluster.
// This is often needed for users or external tools to interact with the vcluster.
// GetVClusterKubeConfig 为指定的虚拟集群生成一个 kubeconfig 字符串。
// 用户或外部工具通常需要此功能来与 vcluster 交互。
// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
// vclusterName: The name of the virtual cluster. / 虚拟集群的名称。
// hostK8sClient: A Kubernetes client connected to the Host Cluster. / 连接到 Host 集群的 Kubernetes 客户端。
// Returns the kubeconfig content as a byte slice and an error.
// 返回 kubeconfig 内容作为字节切片和错误。
func GetVClusterKubeConfig(ctx context.Context, vclusterName string, hostK8sClient kubernetes.Interface) ([]byte, error) {
	utils.GetLogger().Printf("Attempting to get kubeconfig for vcluster '%s'", vclusterName)

	// The vcluster API server URL and CA certificate are needed.
	// We already fetched the CA cert in GetVClusterClient.
	// vcluster API 服务器 URL 和 CA 证书是必需的。
	// 我们已经在 GetVClusterClient 中获取了 CA 证书。
	// The API server URL is the service name in the host cluster.
	// API 服务器 URL 是 host 集群中的服务名称。
	hostNamespace := fmt.Sprintf("vcluster-%s", vclusterName)                                                      // TODO: Get namespace reliably
	vclusterServiceName := vclusterName                                                                            // TODO: Confirm naming convention
	vclusterAPIServerURL := fmt.Sprintf("https://%s.%s.svc.cluster.local:443", vclusterServiceName, hostNamespace) // Use service name

	// Fetch the CA certificate from the host Secret
	// 从 host Secret 获取 CA 证书
	certsSecretName := fmt.Sprintf("vc-certs-%s", vclusterName) // TODO: Confirm secret name
	certsSecret, err := hostK8sClient.CoreV1().Secrets(hostNamespace).Get(ctx, certsSecretName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.NewWithCause(errors.ErrTypeVCluster, fmt.Sprintf("failed to get vcluster certs secret '%s/%s' for kubeconfig", hostNamespace, certsSecretName), err)
	}
	caCertData, ok := certsSecret.Data["ca.crt"] // TODO: Confirm key name
	if !ok || len(caCertData) == 0 {
		return nil, errors.New(errors.ErrTypeVCluster, fmt.Sprintf("vcluster CA certificate not found in secret '%s/%s' for kubeconfig", hostNamespace, certsSecretName))
	}

	// Fetch the token from the secret
	tokenSecretName := fmt.Sprintf("vc-token-%s", vclusterName)
	tokenSecret, err := hostK8sClient.CoreV1().Secrets(hostNamespace).Get(ctx, tokenSecretName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.NewWithCause(errors.ErrTypeVCluster, fmt.Sprintf("failed to get vcluster token secret '%s/%s' for kubeconfig", hostNamespace, tokenSecretName), err)
	}
	token, ok := tokenSecret.Data["token"]
	if !ok || len(token) == 0 {
		return nil, errors.New(errors.ErrTypeVCluster, fmt.Sprintf("vcluster token not found in secret '%s/%s' for kubeconfig", hostNamespace, tokenSecretName))
	}

	// Placeholder: Building a basic kubeconfig structure with CA cert and endpoint.
	// Authentication details are missing in this simplified placeholder.
	// 占位符：使用 CA 证书和端点构建基本的 kubeconfig 结构。
	// 此简化占位符中缺少身份验证详细信息。
	kubeconfigTemplate := `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: %s
    server: %s
  name: %s
contexts:
- context:
    cluster: %s
    user: %s
    namespace: default # Or a specific namespace
  name: %s
current-context: %s
kind: Config
preferences: {}
users:
- name: %s
  user:
    token: %s
`
	// Example user name (can be anything descriptive)
	// 示例用户名称（可以是任何描述性的）
	userName := fmt.Sprintf("vcluster-user-%s", vclusterName)
	clusterNameInKubeconfig := vclusterName
	contextName := vclusterName

	// Base64 encode the CA cert data for embedding in kubeconfig
	// 对 CA 证书数据进行 Base64 编码，以便嵌入 kubeconfig 中
	caCertBase64 := base64.StdEncoding.EncodeToString(caCertData) // Requires "encoding/base64"

	// Format the template with actual values
	// 使用实际值格式化模板
	kubeconfigContent := fmt.Sprintf(kubeconfigTemplate,
		caCertBase64,            // certificate-authority-data
		vclusterAPIServerURL,    // server
		clusterNameInKubeconfig, // cluster name in clusters
		clusterNameInKubeconfig, // cluster name in context
		userName,                // user name in context
		contextName,             // context name
		contextName,             // current-context
		userName,                // user name in users
		string(token),           // token
	)

	utils.GetLogger().Printf("Generated kubeconfig for vcluster '%s'.", vclusterName)

	return []byte(kubeconfigContent), nil
}
