# Istio Envoy
这个仓库会围绕 Istio/Envoy 进一些开发和配置，以实现一些特殊的业务场景需求。

## 特点
- 所有的特殊功能都通过 Istio 的 EnfoyFilter 来实现，会涉及一些 Lua/Golang 的编码
- 每个目录对应一个功能配置，配置会使用 Helm 进行模板化

## 许可证
[MIT](LICENSE)
