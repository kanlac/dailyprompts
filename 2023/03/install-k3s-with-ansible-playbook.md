# 简述通过 Ansible Playbook 自动化安装 K3S 的大致步骤？

1. 在各节点上安装 docker（可选）❶
2. 更新各节点上的 hosts 文件，以便使用域名访问
3. 启动容器，包括 postgres, registry, chartmuseum 等
4. 将 k3s 安装脚本同步到各个节点（离线环境可选）
5. 先在主节点完成 K3S 安装，创建 K8S 命名空间，创建 service account，并获取 k3s token 和 service account token（也可以把 kube-config 文件拷贝到 ansible controller）
6. 通过 token 添加其他服务节点，也就是管理节点（高可用性可选）
7. 通过 token 添加代理节点，也就是工作节点
8. 将 kube-config 文件拷贝到其他节点（如果需要在其它节点使用 `kubectl`，可选）

---

❶ K3S 默认使用 containerd 作为容器运行时，但也可以在执行 K3S 安装脚本时通过参数指定使用 Docker，安装的时候就会下载 Docker。然而如果需要在离线环境下安装，则可以提前手动通过 ansible playbook 在各节点上完成 Docker 的安装。