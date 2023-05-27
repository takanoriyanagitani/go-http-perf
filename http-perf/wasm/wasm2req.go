package wasm2req

import (
	"os"

	perf "github.com/takanoriyanagitani/go-http-perf/http-perf"
	util "github.com/takanoriyanagitani/go-http-perf/util"
)

type Time2ReqWasmBuilder func(wasmBytes []byte) perf.Time2RequestRaw

type Wasm2builder func(wasm []byte) (perf.UnixtimeMicros2RequestBuilder, error)

func (b Wasm2builder) FromPath(wasmname string) (perf.UnixtimeMicros2RequestBuilder, error) {
	return util.ComposeErr(
		os.ReadFile,
		b,
	)(wasmname)
}
