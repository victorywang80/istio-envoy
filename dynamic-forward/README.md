# 按需动态转发
本配置实现了下述流程中的 3-4 步骤
存在两个服务：lobby 和 exampleServer。lobby 负责分配并调度用户访问具体的内部 exampleServer：
1. 用户首次访问 lobby 服务
2. lobby 返回包含目标服务器信息的 URL，格式示例：`https://www.exampleserver.com/v1/product/server/?idn=2887647234`（其中 `idn=2887647234` 表示十进制格式的 IP 地址，对应标准 IP 格式为 `172.17.0.2`）
3. 用户在返回的 URL 中附加业务参数发起后续请求（`idn` 参数位置不固定，可位于 URL 任意位置）
4. Ingress 网关提取 URL 中的 `idn` 参数，执行以下操作：
   - 进行安全验证（包括 IP 白名单检查和参数合法性校验）
   - 将十进制 IP 转换为标准 dotted decimal 格式
   - 将请求转发至该 IP 对应的 exampleServer 服务实例

### 前置要求
- Istio 1.10+ 版本（已启用 Sidecar 自动注入）
- Kubernetes 1.19+ 集群环境
- Helm 3.0+ 客户端工具

### 部署步骤
1. 基于 Istio 部署 helm-chart：
2. helm-chart 示例中未包含 IngressGateway 实例的创建，用户需根据自身环境创建，并确保配置了正确的监听器（Listener）
3. 示例中的证书由 Let's Encrypt 提供，用户需在集群中部署 cert-manager 进行证书管理，并配置 DNS01 验证方式的 Issuer/ClusterIssuer
4. 配置 IP 地址范围：
   - **十进制 IP 转换说明**：`2887647234` = 172×2^24 + 17×2^16 + 0×2^8 + 2 → 标准格式 `172.17.0.2`
   - 示例范围：`2887647234-2887647235`（对应标准 IP：`172.17.0.2-172.17.0.3`）
   - 生产建议：使用 CIDR 格式（如 `10.244.0.0/16`）
   - 实现示例（Calico IP 池）：
     ```yaml
     apiVersion: projectcalico.org/v3
     kind: IPPool
     metadata:
       name: exampleserver-ip-pool
     spec:
       cidr: 10.244.0.0/16
       blockSize: 26
     ```
   - 验证命令：`kubectl get pods -o wide | grep exampleServer`
   - 在服务部署的 Deployment 中用注释引用分配的 IPPool 名称：
     ```yaml
     metadata:
       annotations:
         # 引用示例中的 IP 池
         projectcalico.org/ippool: exampleserver-ip-pool
     ```