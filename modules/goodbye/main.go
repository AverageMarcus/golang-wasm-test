package main

import (
	"github.com/averagemarcus/wasmhook/pkg/wasmhelper"
	"github.com/averagemarcus/wasmhook/pkg/wasmhelper/types"
)

func main() {}

func init() {
	wasmhelper.Fn = func(input types.HookRequest) types.HookRequest {
		input.Name = "Goodbye, WASM"
		input.Result = false
		return input
	}
}
