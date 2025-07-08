package common

import (
	"fmt"

	"github.com/tidwall/gjson"
)

const (
	// 社区版 0
	CE = iota
	// 专业版 1
	PE
)

type Config struct {
	// 获取专业版列表的服务名字（envoy 的 cluster name）
	ApiService string
	// 访问服务时 http 头的中的 host/authority，对就后端服务的虚拟主机配置的域名，其实对于获取这个专业列来说这个没用，proxywasm.DispatchHttpCall 函数的 http 头中没有 authority 值会报错
	ApiHost string
	ApiUri  string
	// 专业版列表
	ProApps map[string]struct{}
	// 是否是专业版标志，将在 http 头增加此 key
	ProAppFlag string
	// 专业版默认标识
	ProAppFlagDefaultValue int8
	// 访问时从 http header 提取 appid 是使用的 key
	AppIdKey []string
	// 多长时间更新一次专业版列表
	TickTime uint32
}

func NewConfig() *Config {
	return &Config{
		ApiService: "Professional-headless.k8s.ccluster.local",
		ApiHost:    "Professional-headless.k8s.ccluster.local",
		ApiUri:     "/Professional/list-pro-apps",
		ProApps:    make(map[string]struct{}),
		// 默为 http header 增加专业版标识 key
		ProAppFlag: "pro-app",
		// 默认值为社区版
		ProAppFlagDefaultValue: CE,
		// 默认从 http header 中获取 appid 的 key
		AppIdKey: []string{"app-id"},
		// 默认每1小时更新一次专业版列表
		TickTime: 1000 * 60 * 60, // 1 小时
	}
}

func ParsePluginConfiguration(data []byte) (*Config, error) {
	config := NewConfig()
	if !gjson.ValidBytes(data) {
		return nil, fmt.Errorf("the plugin configuration is not a valid json: %v", data)
	}

	jsonData := gjson.ParseBytes(data)

	if jsonData.Get("tick_time").Exists() {
		config.TickTime = uint32(jsonData.Get("tick_time").Uint())
	}

	if jsonData.Get("api_service").Exists() {
		config.ApiService = jsonData.Get("api_service").String()
	}

	if jsonData.Get("api_host").Exists() {
		config.ApiHost = jsonData.Get("api_host").String()
	}

	if jsonData.Get("api_url").Exists() {
		config.ApiUri = jsonData.Get("api_url").String()
	}

	if jsonData.Get("pro_app_flag").Exists() {
		config.ProAppFlag = jsonData.Get("pro_app_flag").String()
	}

	if jsonData.Get("pro_app_flag_default_value").Exists() {
		config.ProAppFlagDefaultValue = int8(jsonData.Get("pro_app_flag_default_value").Int())
	}

	if jsonData.Get("app_id_key").Exists() {
		jsonData.Get("app_id_key").ForEach(func(key, value gjson.Result) bool {
			config.AppIdKey = append(config.AppIdKey, value.String())
			return true
		})
	}

	return config, nil
}
