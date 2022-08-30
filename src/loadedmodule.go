package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type LoadedModule struct {
	r   wazero.Runtime
	ctx context.Context
	mod api.Module
}

func NewLoadedModule(wasmBytes []byte) LoadedModule {
	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	if _, err := wasi_snapshot_preview1.Instantiate(ctx, r); err != nil {
		log.Panicln(err)
	}
	mod, err := r.InstantiateModuleFromBinary(ctx, wasmBytes)
	if err != nil {
		log.Panicln(err)
	}

	return LoadedModule{
		r:   r,
		ctx: ctx,
		mod: mod,
	}
}

func (m LoadedModule) Close() {
	m.r.Close(m.ctx)
}

func (m *LoadedModule) Run(input string) ([]interface{}, error) {
	ptr, size := m.getInputString(input)

	ptrs, err := m.mod.ExportedFunction("Run").Call(m.ctx, ptr, size)
	if err != nil {
		return nil, err
	}
	results := []interface{}{}
	for _, ptr := range ptrs {
		val, _ := m.getReturnString(ptr)
		results = append(results, val)
	}

	return results, nil
}

func (m *LoadedModule) getReturnString(ptr uint64) (string, error) {
	if bytes, ok := m.mod.Memory().Read(m.ctx, uint32(ptr>>32), uint32(ptr)); !ok {
		return "", fmt.Errorf("Failed to read string from memory pointer")
	} else {
		return string(bytes), nil
	}
}

func (m *LoadedModule) getInputString(str string) (uint64, uint64) {
	malloc := m.mod.ExportedFunction("malloc")
	strSize := uint64(len(str))
	results, err := malloc.Call(m.ctx, strSize)
	if err != nil {
		log.Panicln(err)
	}

	inPtr := results[0]

	if !m.mod.Memory().Write(m.ctx, uint32(inPtr), []byte(str)) {
		log.Panicf("Failed to set memory")
	}

	return inPtr, strSize
}
