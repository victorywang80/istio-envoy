package updater

import (
	"wasm/common"
	"wasm/internal"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tidwall/gjson"
)

func Update(ctx internal.PluginContext, conf *common.Config) {
	hs := [][2]string{
		{":method", "GET"}, {":path", conf.ApiUri}, {"accept", "*/*"}, {":authority", conf.ApiHost},
	}
	if _, err := proxywasm.DispatchHttpCall(conf.ApiService, hs, []byte(""), nil, 5000, ctx.CallBack); err != nil {
		proxywasm.LogCriticalf("dispatch httpcall failed: %v", err)
	}
	proxywasm.LogInfof("OnTick method is called after dispatchHttpCall ")
}

func HandleApps(apps []byte) map[string]struct{} {
	ProApps := make(map[string]struct{})

	if !gjson.ValidBytes(apps) {
		// 没有获取到专业版列表, 返回空 map
		return ProApps
	}

	jsonData := gjson.ParseBytes(apps)
	jsonData.Get("results.#.app_id").ForEach(func(key, value gjson.Result) bool {
		ProApps[value.String()] = struct{}{}
		return true
	})

	return ProApps
}
