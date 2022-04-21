// Copyright 2020-2021 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
	duplicateBody bool
}

// Override types.DefaultPluginContext.
func (ctx *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &duplicateBodyContext{}
}

func (ctx *duplicateBodyContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	value, _ := proxywasm.GetHttpRequestHeader("x-duplicate")
	if value == "true" {
		ctx.duplicateBody = true
	} else {
		ctx.duplicateBody = false
	}
	return types.ActionContinue
}

type duplicateBodyContext struct {
	// embed the default plugin context
	// so that you don't need to reimplement all the methods by yourself.
	types.DefaultHttpContext
	duplicateBody bool
}

// Override types.DefaultHttpContext.
func (ctx *duplicateBodyContext) OnHttpResponseBody(bodySize int, endOfStream bool) types.Action {
	if !endOfStream {
		// Wait until we see the entire body to duplicate.
		return types.ActionPause
	}

	if ctx.duplicateBody {
		body, _ := proxywasm.GetHttpResponseBody(0, bodySize)
		err := proxywasm.AppendHttpResponseBody(body)
		if err != nil {
			return 0
		}
	}
	return types.ActionContinue
}
