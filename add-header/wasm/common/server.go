package common

import "github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"

// 获取请求 AppId
func GetAppId(conf *Config) string {
	for _, appIdKey := range conf.AppIdKey {
		if appId, err := proxywasm.GetHttpRequestHeader(appIdKey); err == nil {
			return appId
		}
	}
	return ""
}
