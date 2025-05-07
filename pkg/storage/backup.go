// Package storage provides storage-related utilities, including backup functions.
// 包 storage 提供了存储相关的工具，包括备份功能。
package storage

import (
	"context"
	"fmt"
	"github.com/turtacn/chasi-bod/common/errors"
	"github.com/turtacn/chasi-bod/common/utils" // Assuming logger and file ops are here // 假设日志记录器和文件操作在这里
	"github.com/turtacn/chasi-bod/pkg/config/model"
	"path/filepath" // Added for path joining // 添加用于路径拼接
	// Assuming you have an SSH client utility
	// 假设您有一个 SSH 客户端工具
	// "github.com/turtacn/chasi-bod/pkg/sshutil"
)

// sysctlBaseDir is the base directory for sysctl entries in /proc.
// sysctlBaseDir 是 /proc 中 sysctl 条目的基本目录。
const sysctlBaseDir = "/proc/sys" // This constant seems misplaced here, better in a common/system pkg if needed globally. // 这个常量放在这里似乎不合适，如果需要全局使用，最好放在 common/system 包中。

// BackupETCD performs an ETCD snapshot backup on a master node.
// BackupETCD 在主节点上执行 ETCD 快照备份。
// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
// nodeCfg: The configuration for the master node where ETCD is running. / 运行 ETCD 的主节点的配置。
// etcdctlCmdPrefix: The command prefix to execute etcdctl (e.g., "sudo etcdctl --endpoints=<...> --cacert=<...> --cert=<...> --key=<...>").
// etcdctlCmdPrefix: 执行 etcdctl 的命令前缀（例如，“sudo etcdctl --endpoints=<...> --cacert=<...> --cert=<...> --key=<...>”）。
// backupPathOnNode: The destination path on the remote node to save the snapshot. / 远程节点上保存快照的目标路径。
// Returns an error if the backup fails.
// 如果备份失败则返回错误。
func BackupETCD(ctx context.Context, nodeCfg *model.NodeConfig, etcdctlCmdPrefix string, backupPathOnNode string) error {
	utils.GetLogger().Printf("Performing ETCD backup on node %s to path %s...", nodeCfg.Address, backupPathOnNode)

	// TODO: Get SSH client for the master node
	// TODO: 获取主节点的 SSH 客户端
	// sshClient, err := sshutil.NewClient(nodeCfg.Address, nodeCfg.Port, nodeCfg.User, nodeCfg.Password, nodeCfg.PrivateKey)
	// if err != nil { return errors.NewWithCause(errors.ErrTypeNetwork, fmt.Sprintf("failed to create SSH client for %s", nodeCfg.Address), err) }
	// defer sshClient.Close()

	// Ensure the backup directory exists on the remote node
	// 确保远程节点上的备份目录存在
	remoteBackupDir := filepath.Dir(backupPathOnNode)
	utils.GetLogger().Printf("Ensuring remote backup directory %s exists on %s...", remoteBackupDir, nodeCfg.Address)
	// TODO: sshClient.RunCommand(ctx, fmt.Sprintf("sudo mkdir -p %s && sudo chmod 0700 %s", remoteBackupDir, remoteBackupDir)) // Create directory with restricted permissions
	// TODO: 在远程节点上执行 mkdir -p 并设置权限
	utils.GetLogger().Printf("Placeholder: Remote backup directory %s ensured on %s.", remoteBackupDir, nodeCfg.Address)

	// Construct the full backup command
	// 构建完整的备份命令
	backupCmd := fmt.Sprintf("%s snapshot save %s", etcdctlCmdPrefix, backupPathOnNode)

	utils.GetLogger().Printf("Executing ETCD backup command on %s: %s", nodeCfg.Address, backupCmd)
	// TODO: Execute the command remotely via SSH
	// TODO: 通过 SSH 远程执行命令
	// output, err := sshClient.RunCommand(ctx, backupCmd)
	// if err != nil {
	// 	return errors.NewWithCause(errors.ErrTypeSystem, fmt.Sprintf("failed to execute ETCD backup command on %s", nodeCfg.Address), err)
	// }
	// utils.GetLogger().Printf("ETCD backup command output on %s: %s", nodeCfg.Address, output)
	utils.GetLogger().Printf("Placeholder: Executed ETCD backup command on %s.", nodeCfg.Address)

	// Verify the backup file exists and is not empty on the remote node
	// 验证备份文件是否存在且不为空于远程节点上
	checkCmd := fmt.Sprintf("sudo test -s %s", backupPathOnNode) // Use test -s to check if file exists and is not empty
	utils.GetLogger().Printf("Verifying backup file existence and size on %s: %s", nodeCfg.Address, checkCmd)
	// TODO: Execute the check command remotely via SSH
	// TODO: 通过 SSH 远程执行检查命令
	// _, err = sshClient.RunCommand(ctx, checkCmd) // test command exits with 0 on success, non-zero on failure
	// if err != nil {
	// 	return errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("ETCD backup file %s not found or is empty on %s after backup", backupPathOnNode, nodeCfg.Address), err)
	// }
	utils.GetLogger().Printf("Placeholder: ETCD backup file %s verified on %s.", backupPathOnNode, nodeCfg.Address)

	utils.GetLogger().Printf("ETCD backup completed successfully on node %s to path %s.", nodeCfg.Address, backupPathOnNode)
	return nil
}

// BackupConfigFiles backs up specified configuration files from a remote node.
// BackupConfigFiles 从远程节点备份指定的配置文件。
// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
// nodeCfg: The configuration for the remote node. / 远程节点的配置。
// filePaths: A list of file paths on the remote node to back up. / 要备份的远程节点上的文件路径列表。
// backupDirLocal: The destination directory on the local machine to save the backups. / 本地机器上保存备份的目标目录。
// Returns an error if the backup fails.
// 如果备份失败则返回错误。
func BackupConfigFiles(ctx context.Context, nodeCfg *model.NodeConfig, filePaths []string, backupDirLocal string) error {
	if len(filePaths) == 0 {
		utils.GetLogger().Println("No configuration files specified for backup.")
		return nil
	}
	utils.GetLogger().Printf("Backing up configuration files from node %s to local directory %s...", nodeCfg.Address, backupDirLocal)

	// TODO: Get SSH client for the node
	// TODO: 获取节点的 SSH 客户端
	// sshClient, err := sshutil.NewClient(nodeCfg.Address, nodeCfg.Port, nodeCfg.User, nodeCfg.Password, nodeCfg.PrivateKey)
	// if err != nil { return errors.NewWithCause(errors.ErrTypeNetwork, fmt.Sprintf("failed to create SSH client for %s", nodeCfg.Address), err) }
	// defer sshClient.Close()

	// Ensure local backup directory exists
	// 确保本地备份目录存在
	if err := utils.MkdirAll(backupDirLocal, 0755); err != nil { // Assuming utils.MkdirAll exists
		return errors.NewWithCause(errors.ErrTypeIO, fmt.Sprintf("failed to create local backup directory %s", backupDirLocal), err)
	}

	for _, filePath := range filePaths {
		remotePath := filePath
		localPath := filepath.Join(backupDirLocal, nodeCfg.Address, filepath.Base(filePath)) // Save in local dir structure by node/filename

		// Ensure the node-specific subdirectory exists locally
		// 确保本地存在节点特定的子目录
		nodeLocalDir := filepath.Join(backupDirLocal, nodeCfg.Address)
		if err := utils.MkdirAll(nodeLocalDir, 0755); err != nil {
			utils.GetLogger().Printf("Warning: Failed to create local node backup directory %s: %v", nodeLocalDir, err)
			continue // Skip this file if local dir creation fails
		}

		utils.GetLogger().Printf("Copying file from %s:%s to local %s...", nodeCfg.Address, remotePath, localPath)
		// TODO: Copy file from remote to local via SSH/SFTP
		// TODO: 通过 SSH/SFTP 从远程复制文件到本地
		// Example: err := sshClient.CopyFromRemote(ctx, remotePath, localPath) // Assuming CopyFromRemote exists
		// if err != nil {
		// 	utils.GetLogger().Printf("Warning: Failed to backup file %s from %s: %v", remotePath, nodeCfg.Address, err)
		// 	// Decide if this is a fatal error or just a warning
		// 	// 决定这是否是致命错误或只是警告
		// } else {
		// 	utils.GetLogger().Printf("File %s from %s backed up to local %s.", remotePath, nodeCfg.Address, localPath)
		// }
		utils.GetLogger().Printf("Placeholder: File %s from %s backed up to local %s.", remotePath, nodeCfg.Address, localPath)
	}

	utils.GetLogger().Printf("Configuration file backup completed successfully from node %s.", nodeCfg.Address)
	return nil
}

// RestoreETCD restores ETCD from a snapshot on a master node.
// This is a destructive operation and requires careful handling.
// RestoreETCD 在主节点上从快照恢复 ETCD。
// 这是一个破坏性操作，需要小心处理。
// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
// nodeCfg: The configuration for the master node. / 主节点的配置。
// etcdctlCmdPrefix: The command prefix to execute etcdctl. / 执行 etcdctl 的命令前缀。
// snapshotPathOnNode: The path to the snapshot file on the remote node. / 远程节点上的快照文件路径。
// dataDirOnNode: The ETCD data directory on the remote node. / 远程节点上的 ETCD 数据目录。
// Returns an error if the restoration fails.
// 如果恢复失败则返回错误。
func RestoreETCD(ctx context.Context, nodeCfg *model.NodeConfig, etcdctlCmdPrefix string, snapshotPathOnNode string, dataDirOnNode string) error {
	utils.GetLogger().Printf("Restoring ETCD on node %s from snapshot %s to data dir %s...", nodeCfg.Address, snapshotPathOnNode, dataDirOnNode)

	// TODO: Implement ETCD restoration logic via SSH
	// This involves:
	// - Stopping kube-apiserver and etcd services on the master node.
	// - Running `etcdctl snapshot restore <snapshotPathOnNode> --data-dir=<dataDirOnNode> --endpoints=<...>` remotely.
	// - Adjusting file permissions/ownership remotely.
	// - Starting etcd and kube-apiserver services.
	// - Re-adding other control plane components if they were external.
	// 这涉及：
	// - 在主节点上停止 kube-apiserver 和 etcd 服务。
	// - 远程运行 `etcdctl snapshot restore <snapshotPathOnNode> --data-dir=<dataDirOnNode> --endpoints=<...>`。
	// - 远程调整文件权限/所有权。
	// - 启动 ETCD 和 kube-apiserver 服务。
	// - 如果其他控制平面组件是外部的，则重新添加它们。
	utils.GetLogger().Println("Placeholder: ETCD restoration logic needs careful implementation.")
	return errors.New(errors.ErrTypeNotImplemented, "ETCD restoration not implemented yet")
}

// RestoreConfigFiles restores configuration files to a remote node.
// RestoreConfigFiles 将配置文件恢复到远程节点。
// ctx: Context for cancellation and timeouts. / 用于取消和超时的上下文。
// nodeCfg: The configuration for the remote node. / 远程节点的配置。
// backupDirLocal: The local directory containing the backup files (structured by node/filename). / 包含备份文件的本地目录（按节点/文件名结构化）。
// originalFilePaths: A list of original destination paths on the remote node for the files in backupDirLocal.
// originalFilePaths: 远程节点上文件的原始目标路径列表（对应于 backupDirLocal 中的文件）。
// Returns an error if the restoration fails.
// 如果恢复失败则返回错误。
func RestoreConfigFiles(ctx context.Context, nodeCfg *model.NodeConfig, backupDirLocal string, originalFilePaths []string) error {
	if len(originalFilePaths) == 0 {
		utils.GetLogger().Println("No configuration files specified for restoration.")
		return nil
	}
	utils.GetLogger().Printf("Restoring configuration files to node %s from local directory %s...", nodeCfg.Address, backupDirLocal)

	// TODO: Get SSH client for the node
	// TODO: 获取节点的 SSH 客户端
	// sshClient, err := sshutil.NewClient(nodeCfg.Address, nodeCfg.Port, nodeCfg.User, nodeCfg.Password, nodeCfg.PrivateKey)
	// if err != nil { return errors.NewWithCause(errors.ErrTypeNetwork, fmt.Sprintf("failed to create SSH client for %s", nodeCfg.Address), err) }
	// defer sshClient.Close()

	nodeLocalBackupDir := filepath.Join(backupDirLocal, nodeCfg.Address)

	for _, originalRemotePath := range originalFilePaths {
		// Assume local file has the same base name as the original remote file within the node's backup dir
		// 假设本地文件与节点备份目录中的原始远程文件具有相同的基本名称
		localPath := filepath.Join(nodeLocalBackupDir, filepath.Base(originalRemotePath))
		remotePath := originalRemotePath

		exists, err := utils.PathExists(localPath) // Assuming utils.PathExists exists
		if err != nil {
			utils.GetLogger().Printf("Warning: Failed to check existence of local backup file %s: %v", localPath, err)
			continue // Skip this file if local check fails
		}
		if !exists {
			utils.GetLogger().Printf("Warning: Local backup file %s not found, skipping restoration for %s.", localPath, remotePath)
			continue
		}

		utils.GetLogger().Printf("Copying file from local %s to %s:%s...", localPath, nodeCfg.Address, remotePath)
		// TODO: Copy file from local to remote via SSH/SFTP
		// Need to handle permissions and ownership on the remote side after copying (sudo cp, chown, chmod)
		// This often involves using `sudo cp` to overwrite the original file and then `sudo chown`/`sudo chmod`.
		// 需要在复制后处理远程的权限和所有权（sudo cp, chown, chmod）
		// 这通常涉及使用 `sudo cp` 覆盖原始文件，然后使用 `sudo chown`/`sudo chmod`。
		// Example: err := sshClient.CopyToRemote(ctx, localPath, remotePath, "sudo cp -p %s %s && sudo chown --reference=%s %s && sudo chmod --reference=%s %s") // Use --reference to copy original ownership/permissions
		// If original permissions/ownership are not known from backup, hardcode them or use a default.
		// 如果备份不知道原始权限/所有权，则对其进行硬编码或使用默认值。
		// Example simple copy and set permissions:
		// err := sshClient.CopyToRemote(ctx, localPath, "/tmp/temp-file") // Copy to temp location first
		// if err != nil { ... continue }
		// chmodChownCmd := fmt.Sprintf("sudo mv /tmp/temp-file %s && sudo chown root:root %s && sudo chmod 0644 %s", remotePath, remotePath, remotePath) // Example: Set owner root:root and 0644
		// _, err = sshClient.RunCommand(ctx, chmodChownCmd)
		// if err != nil { utils.GetLogger().Printf("Warning: Failed to move/chown/chmod %s on %s: %v", remotePath, nodeCfg.Address, err) }

		utils.GetLogger().Printf("Placeholder: File %s restored to %s.", remotePath, nodeCfg.Address)
	}

	utils.GetLogger().Printf("Configuration file restoration completed successfully to node %s.", nodeCfg.Address)
	return nil
}

// TODO: Implement sshutil package for remote command execution and file transfer
// TODO: 实现 sshutil 包用于远程命令执行和文件传输
// TODO: Implement helper functions for transferring backup files to/from remote storage (S3, NFS etc.)
// TODO: 实现将备份文件传输到/从远程存储（S3、NFS 等）的辅助函数
