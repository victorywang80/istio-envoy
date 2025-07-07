# Istio Envoy
这个仓库围绕 Istio/Envoy 进行开发和配置，以实现特殊业务场景需求。

## 特点
- 所有特殊功能通过 Istio 的 EnvoyFilter 实现，涉及 Lua/Golang 编码
- 每个目录对应一个功能配置，使用 Helm 进行模板化

## 目录结构
- [dynamic-forward](./dynamic-forward): 动态提取 URL 参数中的十进制 IP 地址，以此作为转发目标(通过 Lua 实现 IP 地址转换、安全检查及动态赋值 "X-Host-Port")。适用于以下场景：
  > 需要访问的服务部署在 Kubernetes 集群中，服务 Pod 地址不固定。用户访问服务的地址由调度服务提供，调度服务知晓需访问的所有 Pod 地址，并希望指定用户访问特定 Pod 地址。这种场景多见于游戏领域，例如大厅分配固定房间后，用户的所有后续访问需转发至该固定 Pod 地址，同时利用 Kubernetes 实现便捷的扩缩容及 SSL 终止。详细说明参考 [dynamic-forward](./dynamic-forward/README.md)


## 许可证
[MIT](LICENSE)
