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
var rgpool *u2rPool = builder2pool(rgbuild)

var rawTime2req hperf.Time2RequestRaw = reqgen.ToTime2RequestRaw()
var rawTime2reqLock sync.Mutex

type u2rPool struct{ pool sync.Pool }

func builder2pool(b hperf.UnixtimeMicros2RequestBuilder) *u2rPool {
	return &u2rPool{
		pool: sync.Pool{
			New: func() any {
				var um2r hperf.UnixtimeMicros2Request = must(b.Build())
				var t2rr hperf.Time2RequestRaw = um2r.ToTime2RequestRaw()
				return t2rr
			},
		},
	}
}

func (p *u2rPool) getRaw() any                    { return p.pool.Get() }
func (p *u2rPool) get() hperf.Time2RequestRaw     { return p.getRaw().(hperf.Time2RequestRaw) }
func (p *u2rPool) put(um2r hperf.Time2RequestRaw) { p.pool.Put(um2r) }
func (p *u2rPool) time2request(t time.Time) (raw []byte, e error) {
	var t2r hperf.Time2RequestRaw = p.get()
	defer p.put(t2r)
	return t2r(t)
}

func pool2t2r(p *u2rPool) hperf.Time2RequestRaw { return p.time2request }

func time2reqSingle(t time.Time) (serialized []byte, e error) {
	rawTime2reqLock.Lock()
	defer rawTime2reqLock.Unlock()
	return rawTime2req(t)
}

var usePool bool = "pool" == os.Getenv("ENV_USE_POOL")

var time2req hperf.Time2RequestRaw = util.Select(
	time2reqSingle,
	pool2t2r(rgpool),
	usePool,
)

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
