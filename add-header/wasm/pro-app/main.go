package main

import (
	"strconv"
	"wasm/common"
	"wasm/filter"
	"wasm/internal"
	"wasm/updater"

	_ "github.com/wasilibs/nottinygc"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

//export sched_yield
func sched_yield() int32 {
	return 0
}

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	internal.PluginContext
	app       *internal.App
	config    *common.Config
	contextID uint32
}

type proApp struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	internal.HttpContext
	app       *internal.App
	config    *common.Config
	contextID uint32
}

func (ctx *vmContext) OnVMStart(vmConfigurationSize int) types.OnVMStartStatus {
	return types.OnVMStartStatusOK
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{
		contextID: contextID,
		app:       new(internal.App),
	}
}

// Override types.DefaultPluginContext.
func (ctx *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &proApp{
		contextID: contextID,
		config:    ctx.config,
		app:       ctx.app,
	}
}

// Override types.DefaultPluginContext.
func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	data, err := proxywasm.GetPluginConfiguration()
	if err != nil {
		proxywasm.LogCriticalf("error reading vm configuration: %v", err)
		return types.OnPluginStartStatusFailed
	}

	config, err := common.ParsePluginConfiguration(data)
	if err != nil {
		proxywasm.LogCriticalf("error parsing plugin configuration: %v", err)
		return types.OnPluginStartStatusFailed
	}
	ctx.config = config

	// 专业版列表请求回调函数定义
	ctx.CallBack = func(numHeaders, bodySize, numTrailers int) {
		result, err := proxywasm.GetHttpCallResponseBody(0, bodySize)
		if err != nil {
			proxywasm.LogError(err.Error())
			return
		}

		tempProApps := updater.HandleApps(result)
		if len(tempProApps) != 0 {
			ctx.app.ProApps = tempProApps
			// 获取的专业版列表多存一份到 ctx.config.ProApps
			ctx.config.ProApps = tempProApps
		} else {
			// 获取失败时使用备份的列表内容
			ctx.app.ProApps = ctx.config.ProApps
		}
	}

	if err := proxywasm.SetTickPeriodMilliSeconds(ctx.config.TickTime); err != nil {
		proxywasm.LogCriticalf("failed to set tick period: %v", err)
		return types.OnPluginStartStatusFailed
	}
	proxywasm.LogInfof("set tick period milliseconds: %d", ctx.config.TickTime)

	return types.OnPluginStartStatusOK
}

// Override types.DefaultPluginContext.
func (ctx *pluginContext) OnTick() {
	updater.Update(ctx.PluginContext, ctx.config)
}

// Override types.DefaultHttpContext.
func (ctx *proApp) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	// 从 http 头中获取 appid
	ctx.app.AppId = common.GetAppId(ctx.config)

	// 判断过滤是否是专业版
	filter.FilterHeader(ctx.app, ctx.config)

	// 增加专业版 http 头标志
	if err := proxywasm.AddHttpRequestHeader(ctx.config.ProAppFlag, strconv.Itoa(int(ctx.app.ProAppValue))); err != nil {
		proxywasm.LogCriticalf("failed to set pro flag to header: %v", err)
	}

	return types.ActionContinue
}
