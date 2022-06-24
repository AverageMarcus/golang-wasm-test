.PHONY: build
build:
	@tinyjson -all ./src/pkg/wasmhelper/types/types.go
	@find ./modules -maxdepth 1 -mindepth 1 -type d | xargs -I {} -n 1 sh -c 'cd {} && tinygo build -o  ./module.wasm -scheduler=none -target wasi .'

.PHONY: run
run:
	@cd src && go run .
