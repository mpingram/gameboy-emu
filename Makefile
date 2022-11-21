BINARY_NAME=gameboy-emu

build:
	GOARCH=wasm GOOS=js go build -o dist/bundle.wasm main.go

serve:
	python3 -m http.server --directory dist

run: build serve
