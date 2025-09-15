# Chasi-Bod ( chaass-board )

## Build, Share and Run Multi-Tenant Business Systems on Shared Kubernetes Cluster with Virtual Cluster Isolation, Simplified

[![License](https://img.shields.io/github/license/turtacn/chasi-bod)](https://github.com/turtacn/chasi-bod/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/turtacn/chasi-bod)](https://goreportcard.com/report/github.com/turtacn/chasi-bod)
[![Release](https://img.shields.io/github/v/release/turtacn/chasi-bod)](https://github.com/turtacn/chasi-bod/releases/latest)

## ✨ What is Chasi-Bod?

`Chasi-Bod` (pronounced similar to "cha-si-board") is an open-source project inspired by the concepts of [sealer](https://github.com/sealerio/sealer) [1] and [vcluster](https://github.com/loft-sh/vcluster) [2]. It provides a powerful and simplified platform to **fuse multiple complex business systems** onto a **single, shared underlying Kubernetes cluster**, while ensuring **strong isolation** between them by leveraging virtual Kubernetes clusters (`vcluster`).

Think of `chasi-bod` as the "body" built upon the "chassis" (the underlying hardware/VMs) that provides a robust, segmented structure for different "components" (your business systems) to reside in securely and efficiently.

In today's cloud-native landscape, deploying many independent business systems often leads to infrastructure sprawl, high costs, and operational complexity. While traditional Kubernetes namespaces offer some isolation, they fall short for multi-tenant scenarios requiring stricter control-plane separation and tenancy-specific configurations.

`Chasi-bod` addresses these challenges by:

* **Building** a unified, image-based platform artifact that bundles the OS, Kubernetes, container runtime, and core platform components.
* **Sharing** this platform artifact for consistent and repeatable deployments.
* **Running** the platform as a Host Kubernetes cluster.
* **Managing** the lifecycle of isolated virtual Kubernetes clusters (`vclusters`) within the Host cluster for each business system or group.
* **Simplifying** the deployment, configuration, and operations of both the platform and the business systems running on it.

## 🚀 Features & Goals

* **Image-Based Platform Management:** Package the entire infrastructure stack (OS, K8s, Runtime, `vcluster` components) into a versioned image for simplified deployment and upgrade.
* **Multi-Tenancy with vCluster Isolation:** Provide each business system with its own isolated virtual Kubernetes cluster, enhancing security and flexibility compared to namespaces.
* **Simplified Business System Onboarding:** Define clear standards and templates for deploying applications into their respective `vclusters`.
* **Unified Configuration Management:** Centralize and streamline complex configurations for networking, storage, system parameters (`sysctl`), and resource partitioning.
* **Optimized for Mixed Workloads:** Consider and provide guidance/configurations for different application types (CPU, I/O, Memory intensive) running side-by-side.
* **Full Lifecycle Management:** Support building, distributing, deploying, upgrading, scaling, and operating the entire stack.
* **Built-in DFX Capabilities:** Incorporate design for excellence, including observability, reliability, and testability.

## 💡 Core Concepts

* **Chasi-Bod Platform Image:** The portable artifact (ISO, OVA, QCOW2, etc.) containing everything needed to boot up the Host Cluster. Inspired by sealer's ClusterImage.
* **Host Cluster:** The foundational Kubernetes cluster running directly on the infrastructure nodes, managed by Chasi-Bod. It hosts the `vcluster` pods.
* **vCluster (Virtual Cluster):** A lightweight Kubernetes cluster running *inside* a namespace on the Host Cluster. Each business system interacts with its own vCluster API.
* **Business System Application:** The actual application workload deployed within a specific vCluster.
* **Chasi-Bod Management Plane:** The core logic of the Chasi-Bod project responsible for building images, deploying/managing the Host Cluster, provisioning `vclusters`, and orchestrating application deployments.

## 🏗️ Architecture Overview

`Chasi-Bod` employs a layered architecture, covering everything from the bootloader up to the application deployment layer.

At a high level, it consists of:

1.  **Infrastructure Layer:** Physical/Virtual Machines, OS, Filesystem, Network, Storage.
2.  **Host Kubernetes Layer:** The shared Kubernetes cluster running on the infrastructure.
3.  **Virtual Cluster / Application Layer:** Isolated `vcluster` instances and the business applications running within them.
4.  **Chasi-Bod Management Plane:** The control logic orchestrating build, deploy, and manage operations across all layers.

For a detailed breakdown, please refer to the [architecture documentation / 架构设计](docs/architecture.md).

## 🧠 How It Works

1.  **Build:** Define your platform requirements (OS, K8s version, CNI/CSI plugins, base tools) using Chasi-Bod's configuration. Chasi-Bod builds a reproducible Platform Image.
2.  **Share:** Distribute the Platform Image (e.g., upload to a repository).
3.  **Run (Deploy):** Use the Chasi-Bod CLI or API to deploy the Platform Image onto bare metal or VMs. This automatically sets up the Host OS and the Host Kubernetes Cluster.
4.  **Manage:** Use Chasi-Bod to provision isolated `vclusters` for each business system based on predefined templates.
5.  **Deploy Applications:** Use Chasi-Bod to deploy business system applications (e.g., Helm Charts) into their designated `vclusters`. Chasi-Bod ensures proper configuration injection and lifecycle management within the virtualized environment.

## 🌱 Getting Started

This guide will walk you through the basic steps to get `chasi-bod` up and running.

### Prerequisites

*   Go 1.18+
*   Docker
*   A running Kubernetes cluster (e.g., kind, Minikube, or any other standard cluster) to act as the "host" cluster.

### 1. Build the `chasi-bod` CLI

Clone the repository and build the `chasi-bod` binary:

```bash
git clone https://github.com/turtacn/chasi-bod.git
cd chasi-bod
go build -o chasi-bod cmd/chasi-bod/main.go
```

### 2. Create a Configuration File

Create a file named `chasi-bod.yaml` with the following content. This file defines a `vcluster` named `my-vcluster`.

```yaml
apiVersion: chasi-bod.io/v1alpha1
kind: PlatformConfig
metadata:
  name: chasi-bod-platform
vclusters:
  my-vcluster:
    name: my-vcluster
    namespace: vcluster-my-vcluster
    kubernetesVersion: v1.27.3
```

### 3. Create a vCluster

Now, use the `chasi-bod` CLI to create the `vcluster` in your host Kubernetes cluster. Make sure your `kubeconfig` is pointing to your host cluster.

```bash
./chasi-bod --config chasi-bod.yaml vcluster create my-vcluster --wait
```

This command will:
1.  Create a namespace `vcluster-my-vcluster` in your host cluster.
2.  Deploy the `vcluster` components into that namespace.
3.  Wait for the `vcluster` to be ready.

### 4. Interact with the vCluster

Once the `vcluster` is created, you can get a `kubeconfig` to interact with it:

```bash
./chasi-bod vcluster connect my-vcluster > vcluster-kubeconfig.yaml
```

You can now use this `kubeconfig` with `kubectl` to interact with your `vcluster`:

```bash
kubectl --kubeconfig=vcluster-kubeconfig.yaml get namespaces
```

### 5. Delete the vCluster

To clean up, you can delete the `vcluster`:

```bash
./chasi-bod vcluster delete my-vcluster --wait
```

## 🤝 Contributing

We welcome contributions! Please see the [CONTRIBUTING.md](CONTRIBUTING.md) file (Coming Soon) for details on how to get involved.

## 📄 License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.

## 📞 Contact & Community

(Coming Soon)

* How to reach the maintainers.
* Links to community channels (Slack, WeChat, Mailing List, etc.).

---

## References

- [1] sealer - Build, Share and Run Both Your Kubernetes Cluster and Distributed Applications (Project under CNCF) https://github.com/sealerio/sealer
- [2] vCluster - Create fully functional virtual Kubernetes clusters https://github.com/loft-sh/vcluster