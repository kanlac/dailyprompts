# 卸载 K3S 集群有哪些步骤
1. 备份需要的数据和文件
2. 确保所有应用已经停止
3. 先清除所有命名空间的资源，再删除所有命名空间
4. 执行 k3s 卸载脚本
5. 最后，清空 datapoint 中的 kine 表（包含 bootstrap 信息，如果不清理会影响下次安装）