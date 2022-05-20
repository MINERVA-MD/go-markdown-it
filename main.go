package main

import (
	"go-markdown-it/pkg"
	"syscall/js"
)

// function definition
func parse(this js.Value, i []js.Value) interface{} {
	// Process Results
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{Html: true, Typography: true, LangPrefix: "language-"})

	return js.ValueOf(md.Render(i[0].String(), &pkg.Env{}))
}

func registerCallbacks() {
	js.Global().Set("parse", js.FuncOf(parse))
}

func main() {
	c := make(chan struct{}, 0)

	println("WASM Go Initialized")
	// register functions
	registerCallbacks()
	<-c
}
