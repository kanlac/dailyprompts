# Ingress 四层端口开启方法

1. 确保容器监听端口并配置好 Service
2. Ingress ConfigMap
    - TCP: `kubectl edit -n ingress-nginx cm/nginx-ingress-tcp`
    - UDP: `kubectl edit -n ingress-nginx cm/nginx-ingress-udp`
    
    ```yaml
    ---
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: nginx-ingress-tcp  # or nginx-ingress-udp
      namespace: ingress-nginx
    data:
    	# 在这里加上要开启的端口配置, 格式:
    	# "INGRESS_LISTNING_PORT": NAMESPACE/SERVICE_NAME:SERVICE_PORT
    	"20000": ncompass/nlc-fluentd:20000
    ```
    
3. Ingress Service
    
    `kubectl edit svc/ingress-nginx-controller -n ingress-nginx`
    
    ```yaml
    ports:
    	# 添加配置列表
      - name: fluentd-20000
        port: 20000       # 节点端口
        protocol: UDP     # TCP or UDP
        targetPort: 20000 # Ingress 监听端口
    ```
    
4. 检查 service 配置是否成功
    
    `kubectl get svc ingress-nginx-controller -n ingress-nginx`
    
5. 检查端口是否开启
    - TCP: `nc -zv localhost 20000`
    - UDP: `echo "hello message" | nc -u localhost 20000`
6. 若端口未成功开启，检查 ingress 日志：
    
    `kubectl logs -f -n ingress-nginx -l app.kubernetes.io/name=ingress-nginx`
    

---

### 附：Ingress 开启四层转发

`kubectl edit ds/ingress-nginx-controller -n ingress-nginx`

添加启动参数：

```yaml
- --tcp-services-configmap=$(POD_NAMESPACE)/nginx-ingress-tcp
- --udp-services-configmap=$(POD_NAMESPACE)/nginx-ingress-udp
```
