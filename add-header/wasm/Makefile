.PHONY: pro-app
pro-app:
	tinygo build -o ./pro-app/main.wasm -scheduler=none -gc=custom -tags=custommalloc -target=wasi ./pro-app/main.go
