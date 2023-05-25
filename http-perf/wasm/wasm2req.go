package wasm2req

import (
	perf "github.com/takanoriyanagitani/go-http-perf/http-perf"
)

type Time2ReqWasmBuilder func(wasmBytes []byte) perf.Time2RequestRaw
