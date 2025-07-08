package internal

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type PluginContext struct {
	types.DefaultPluginContext
	CallBack func(numHeaders, bodySize, numTrailers int)
}

type HttpContext struct {
	types.DefaultHttpContext
}

type App struct {
	ProApps     map[string]struct{}
	ProAppValue int8
	AppId       string
}
