package main

import (
	"syscall/js"
)

func localStore(key string, value interface{}) {
	js.Global().Get("localStorage").Set(key, value)
}

func localLoad(key string) js.Value {
	return js.Global().Get("localStorage").Get(key)
}
