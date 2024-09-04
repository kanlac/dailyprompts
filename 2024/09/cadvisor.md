# 介绍 cAdvisor 以及它的部署方式

## 简介

cAdvisor（Container Advisor）是一个开源的容器资源使用和性能分析工具，主要用于收集、聚合、处理和导出运行中容器的信息。它由 Google 开发，现在是 Kubernetes 生态系统中的一个重要组件。

## 特点

- 低开销：设计为对系统性能影响最小
- 实时监控：提供近实时的资源使用信息
- 自包含：无需额外的依赖即可运行
- 支持 Web UI：提供简单的 Web 界面查看统计信息
- 支持导出数据到多种监控系统，如 Prometheus、InfluxDB 等

## 部署方式

1. 在 Kubernetes 中，cAdvisor 通常作为 Kubelet 的一部分运行，不需要单独部署
2. 作为独立的 Docker 容器部署
3. 在 Kubernetes 集群中部署：
    - 可以作为 DaemonSet 部署在集群中的每个节点上
    - 这种方式适用于需要更详细监控或自定义配置的场景
4. 二进制部署：
    - 可以直接在主机上运行 cAdvisor 二进制文件
    - 适用于非容器化环境或特殊需求场景