

static: wasm binary
	rm -rf docs
	bin/program -static docs
	cp -a web/* docs/web

run: wasm binary
	bin/program

wasm:
	GOARCH=wasm GOOS=js go build -ldflags="-s -w" -o web/app.wasm

binary: bin
	go build -ldflags "-X main.GitCommit=$(GIT_COMMIT)" -o bin/program

bin:
	mkdir -p bin