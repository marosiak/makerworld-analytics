build:
	GOARCH=wasm GOOS=js go build -o web/app.wasm ./cmd/server
	go build -o output/app_binary ./cmd/server


build_optimized_wasm:
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o web/app.wasm ./cmd/server

	wasm-opt ./web/app.wasm \
		--enable-bulk-memory \
		-Oz \
		--vacuum \
		--remove-unused-module-elements \
		--dce \
		--output=./web/app.final.wasm

	

build_static:
	cd docs && go run .././cmd/static/
	GOARCH=wasm GOOS=js go build -o ./docs/web/app.wasm ./cmd/server

run:
	GOARCH=wasm GOOS=js go build -o web/app.wasm ./cmd/server
	go run -o output/app_binary ./cmd/server