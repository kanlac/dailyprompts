# 镜像中的敏感数据处理

制作镜像时，如何存储加密内容？如何避免敏感数据打到镜像里面？

1. 使用 `.dockerignore` 文件来防止敏感文件在构建过程中被打包进镜像
2. 编排层面通过环境变量注入，避免在 Dockerfile 中硬编码
3. 在 Kubernetes 中使用 Secret 来外部挂载。在 `Pod` 的定义中，你可以通过 `env` 或者 `volume` 挂载来使用这些 Secrets
    
    ```yaml
    env:
    - name: SECRET_KEY
      valueFrom:
        secretKeyRef:
          name: my-secret
          key: SECRET_KEY
    ```
4. 挂载卷
5. 使用外部的加密存储服务，如 HashiCorp Vault、AWS Secrets Manager、Azure Key Vault 等
6. 如果在构建镜像过程中确实需要临时访问敏感数据，可以使用**多阶段构建**。在第一阶段处理敏感数据（如下载依赖、生成证书等），然后在最后阶段中将这些数据从镜像中移除
    - 示例
        
        ```
        # Stage 1: 构建阶段
        FROM golang:alpine AS builder
        # 下载依赖、编译代码等需要密钥的操作
        RUN --mount=type=secret,id=my_secret_key \\
            ./build.sh
        
        # Stage 2: 生成最终镜像
        FROM alpine
        COPY --from=builder /app /app
        # 不包含任何敏感数据
        
        ```