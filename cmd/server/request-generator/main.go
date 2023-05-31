package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	hperf "github.com/takanoriyanagitani/go-http-perf/http-perf"
	wasm "github.com/takanoriyanagitani/go-http-perf/http-perf/wasm"
	wtperf "github.com/takanoriyanagitani/go-http-perf/http-perf/wasm/wasmtime"
	util "github.com/takanoriyanagitani/go-http-perf/util"
)

func must[T any](t T, e error) T {
	if nil != e {
		panic(e)
	}
	return t
}

var wasmLocation string = os.Getenv("ENV_WASM_PATH")
var wasm2reqgen wasm.Wasm2builder = wasm.Wasm2builder(wtperf.Wasm2builderSpeed)
var rgbuild hperf.UnixtimeMicros2RequestBuilder = must(wasm2reqgen.FromPath(wasmLocation))
var reqgen hperf.UnixtimeMicros2Request = must(rgbuild.Build())
var rawTime2req hperf.Time2RequestRaw = reqgen.ToTime2RequestRaw()
var rawTime2reqLock sync.Mutex

func time2req(t time.Time) (serialized []byte, e error) {
	rawTime2reqLock.Lock()
	defer rawTime2reqLock.Unlock()
	return rawTime2req(t)
}

func mustNil(e error) {
	if nil != e {
		log.Fatal(e)
	}
}

func main() {
	var s *http.Server = &http.Server{
		Addr:           ":53080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	var handler http.HandlerFunc = func(w http.ResponseWriter, _ *http.Request) {
		var now time.Time = time.Now()
		generated, e := time2req(now)
		util.Select(
			func() {
				w.WriteHeader(500)
				_, ew := w.Write([]byte("Unable to create a request bytes."))
				mustNil(ew)
				log.Printf("Unable to create a request bytes: %v\n", e)
			},
			func() {
				var hdr http.Header = w.Header()
				hdr.Set("Content-Type", "application/octet-stream")
				_, ew := w.Write(generated)
				mustNil(ew)
			},
			nil == e,
		)()
	}

	http.HandleFunc("/now2req", handler)

	log.Fatal(s.ListenAndServe())
}
