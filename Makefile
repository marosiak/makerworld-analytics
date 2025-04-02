build:
	GOARCH=wasm GOOS=js go build -o web/app.wasm ./cmd/server
	go build -o output/app_binary ./cmd/server

build_static:
	cd docs && go run .././cmd/static/
	GOARCH=wasm GOOS=js go build -o ./docs/web/app.wasm ./cmd/server

run:
	GOARCH=wasm GOOS=js go build -o web/app.wasm ./cmd/server
	go run -o output/app_binary ./cmd/server