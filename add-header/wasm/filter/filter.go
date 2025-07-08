package filter

import (
	"wasm/common"
	"wasm/internal"
)

func FilterHeader(ctx *internal.App, conf *common.Config) {
	if len(ctx.ProApps) != 0 {
		if _, ok := ctx.ProApps[ctx.AppId]; ok {
			ctx.ProAppValue = common.PE
			return
		} else {
			ctx.ProAppValue = conf.ProAppFlagDefaultValue
		}
	}
}
