# 静态 Pod

由 kubelet 直接管理的 Pod，位于 /etc/kubernetes/manifests，不受 API server 控制（API server 自身也属于静态 Pod），但后者还是能看到它，名字会加上节点名后缀。

控制平面组件（如 api-server, controller-manager, scheduler, etcd）通常作为静态 Pod 运行。

由于静态 Pod 直接受 kubelet 管理，不经过 API 服务器，因此它们无法访问 ServiceAccount, ConfigMap, Secret 等 API 对象，也不支持通过 Kubernetes API 的大多数动态特性，如环境变量注入等。
