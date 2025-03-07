# 离线集群镜像管理

## 离线集群镜像管理涉及哪些问题？

下述子标题。

## 如何管理部署清单中的镜像地址？

### 集群预置应用中的镜像地址

可以确定有哪些镜像，kustomize 加前缀部署。

### 产品中的镜像地址

无法确定有哪些镜像，不加前缀，不修改部署清单，所有节点配置容器运行时 rewrite 所有镜像仓库地址，这样可以确保开发环境和生产环境使用相同的部署方式。

## 如何准备镜像？

### 集群预置应用中的镜像准备

1. 通过维护一个镜像列表，自动化完成预置应用（ingress-nginx, kube-prometheus 等）镜像的打包，制作集群安装包时，CI 从文件服务器拉取镜像包
2. 安装加载镜像后，通过同一份镜像列表打标签加前缀、推送
3. kustomize 加前缀部署

### 产品安装包中的镜像准备

需要考虑的问题：如何提取产品中的所有镜像？包括 templates 和 values，格式不确定。

两个思路：

1. 渲染模版后按路径替换——缺点：会导致无法直接使用 helm 命令部署，增加调试复杂性
2. 正则匹配替换——缺点：依赖确定的格式

最终选择方案——编写正则替换函数。打包时可以用来提取镜像列表，部署时可以用于加前缀；做好测试用例；避免在应用中使用 YAML 变量语法：

```go
// 逐行读取文件，匹配 regexPattern，并将匹配部分替换为 repl（若非空）
func RegexReplaceFile(filepath, regexPattern string, repl string) ([]MatchedLine, error)
```
