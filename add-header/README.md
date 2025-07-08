# add-header

## 概述
基于 Istio/Envoy WASM 过滤器实现的动态请求头注入模块，支持根据请求上下文和外部配置动态添加路由决策所需的 HTTP 头字段。

## 功能特点
- **动态头注入**：根据请求特征和缓存的版本信息自动添加路由头
- **高效缓存**：定时从内部 API 同步版本配置并本地缓存，当内部调用服务失效时会使用之前的脏缓存数据，减少外部依赖
- **WASM 实现**：使用 Proxy-Wasm Go SDK 开发，轻量高效且资源占用低
- **灵活配置**：支持通过 Helm 参数自定义缓存周期、头字段名称和 API 端点

## 技术实现
### 核心组件
- **WASM 过滤器**：`wasm/pro-app/main.go` - 实现请求处理和头注入逻辑
- **配置更新器**：`wasm/updater/updater.go` - 定时从 API 获取并更新配置
- **Helm 图表**：`helm-chart/` - 提供完整部署模板和参数配置

### 工作流程
1. **配置同步**：定期从内部 API (`/professional/list-pro-apps`) 获取版本标识列表
2. **请求处理**：解析入站请求头 `x-id`
3. **缓存查询**：检查 app_id 是否在专业版列表中
4. **头字段注入**：若匹配则添加 `pro-app: 1` 头
5. **路由决策**：后端服务根据注入的头字段进行版本路由

## 前置要求
- Kubernetes 1.19+ 集群
- Istio 1.10+ 已安装并启用 Sidecar 注入
- Helm 3.0+ 客户端
- TinyGo 0.29.0+ (用于本地开发 WASM 模块)

## 部署步骤
### 1. 构建 WASM 模块
```bash
cd wasm
make pro-app
# 或使用 Docker 构建
docker build -f Dockerfile.proapp -t add-header-wasm:latest .
```
### 2. 部署 WASM 模块(2种方式)
- 将 main.wasm 文件复制到所有部署 IngressGateway 实例的节点，修改 IngressGateway 部署的 YAML 模板，挂载对应的目录并在配置中引用（本示例的 Helm 部署采用此方式）
- 使用 Istio WasmPlugin 资源进行配置部署（需 Istio 高版本支持），配置方式请参考 [Istio 官方文档](https://istio.io/latest/docs/reference/config/proxy_extensions/wasm-plugin/)

## 配置示例
详见 [Helm 配置示例](helm-chart/values.yaml)

## 开发指南
### 本地测试
1. 启动测试服务器
```bash
cd wasm/test
go run http_server.go
```
2. 发送测试请求
```bash
curl -H "X-ID: YTRm5Nzjx51vM5MHsUqKApFGEC9MeoJt" http://localhost:9070
```