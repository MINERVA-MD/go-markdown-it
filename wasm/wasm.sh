$Env:GOOS = "js"; $Env:GOARCH = "wasm"; go build -o wasm/main.wasm
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./wasm

#$Env:GOOS = "windows"; $Env:GOARCH = "amd64";
