package main

import (
	"github.com/averagemarcus/wasmhook/pkg/wasmhelper"
	"github.com/averagemarcus/wasmhook/pkg/wasmhelper/types"
)

func main() {
	wasmhelper.Fn = func(input types.HookRequest) types.HookRequest {
		input.Name = "Hello, WASM"
		input.Result = true
		return input
	}
}
