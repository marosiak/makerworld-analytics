build:
	GOARCH=wasm GOOS=js go build -o web/app.wasm
	go build -o output/app_binary

run: output
	./output/app_binary