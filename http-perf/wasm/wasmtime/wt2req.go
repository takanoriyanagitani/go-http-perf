package wt2req

import (
	"bytes"
	"errors"
	"runtime"

	"github.com/bytecodealliance/wasmtime-go/v9"

	perf "github.com/takanoriyanagitani/go-http-perf/http-perf"
	util "github.com/takanoriyanagitani/go-http-perf/util"
)

var ErrInvalidInt error = errors.New("invalid integer")
var ErrInvalidFunc error = errors.New("invalid micros2request")
var ErrInvalidExtern error = errors.New("invalid extern")
var ErrInvalidMemory error = errors.New("invalid memory")

type wasm2module func(wasm []byte) (*wasmtime.Module, error)

type engine2builder struct {
	engine   *wasmtime.Engine
	time2req string
	address  string
	wasm     []byte
}

func (w engine2builder) build() (module2builder, error) {
	var w2m wasm2module = util.CurryErr(wasmtime.NewModule)(w.engine)
	var store wasmtime.Storelike = wasmtime.NewStore(w.engine)
	module, em := w2m(w.wasm)
	var time2req string = w.time2req
	var address string = w.address
	return module2builder{
		time2req,
		address,
		store,
		module,
	}, em
}

type module2builder struct {
	time2req string
	address  string
	store    wasmtime.Storelike
	module   *wasmtime.Module
}

func (m module2builder) module2instance() (*wasmtime.Instance, error) {
	return wasmtime.NewInstance(m.store, m.module, nil)
}

func (m module2builder) build() (wasmtime2requestBuilder, error) {
	instance, ei := m.module2instance()
	name2func := util.CurryErrIII(getFunc)(instance)(m.store)
	var time2req string = m.time2req
	var address string = m.address
	var store wasmtime.Storelike = m.store
	return wasmtime2requestBuilder{
		time2req,
		address,
		store,
		instance,
		name2func,
	}, ei
}

type wasmtime2requestBuilder struct {
	time2req  string
	address   string
	store     wasmtime.Storelike
	instance  *wasmtime.Instance
	name2func func(name string) (*wasmtime.Func, error)
}

func any2i5(a any) (int32, error) {
	switch i := a.(type) {
	case int32:
		return i, nil
	default:
		return -1, ErrInvalidInt
	}
}

func (b wasmtime2requestBuilder) func2addr(f *wasmtime.Func) (int32, error) {
	return util.ComposeErr(
		func(s wasmtime.Storelike) (any, error) { return f.Call(s) },
		any2i5,
	)(b.store)
}

func (b wasmtime2requestBuilder) getAddress() (int32, error) {
	return util.ComposeErr(
		b.name2func,
		b.func2addr,
	)(b.address)
}

func getFunc(i *wasmtime.Instance, s wasmtime.Storelike, name string) (*wasmtime.Func, error) {
	var f *wasmtime.Func = i.GetFunc(s, name)
	return util.Select(
		func() (*wasmtime.Func, error) { return nil, ErrInvalidFunc },
		func() (*wasmtime.Func, error) { return f, nil },
		nil != f,
	)()
}

func (b wasmtime2requestBuilder) getTime2Req() (*wasmtime.Func, error) {
	return b.name2func(b.time2req)
}

func extern2memory(ex *wasmtime.Extern) (*wasmtime.Memory, error) {
	var m *wasmtime.Memory = ex.Memory()
	return util.Select(
		func() (*wasmtime.Memory, error) { return nil, ErrInvalidMemory },
		func() (*wasmtime.Memory, error) { return m, nil },
		nil != m,
	)()
}

func (b wasmtime2requestBuilder) getExtern(name string) (*wasmtime.Extern, error) {
	var ex *wasmtime.Extern = b.instance.GetExport(b.store, name)
	return util.Select(
		func() (*wasmtime.Extern, error) { return nil, ErrInvalidExtern },
		func() (*wasmtime.Extern, error) { return ex, nil },
		nil != ex,
	)()
}

func (b wasmtime2requestBuilder) getMemory() (*wasmtime.Memory, error) {
	return util.ComposeErr(
		b.getExtern,
		extern2memory,
	)("memory")
}

func (b wasmtime2requestBuilder) build() (wasmtime2request, error) {
	micros2req, ef := b.getTime2Req()
	memory, em := b.getMemory()
	address, ea := b.getAddress()
	var buffer *bytes.Buffer = new(bytes.Buffer)
	var e error = errors.Join(
		ef,
		em,
		ea,
	)
	var store wasmtime.Storelike = b.store
	return wasmtime2request{
		micros2req,
		store,
		memory,
		address,
		buffer,
	}, e
}

type wasmtime2request struct {
	micros2req *wasmtime.Func
	store      wasmtime.Storelike
	memory     *wasmtime.Memory
	address    int32
	buffer     *bytes.Buffer
}

func (w wasmtime2request) us2req(micros int64) (any, error) {
	return w.micros2req.Call(w.store, micros)
}

func (w wasmtime2request) toMicros2req() perf.UnixtimeMicros2Request {
	var u2q func(int64) (int32, error) = util.ComposeErr(
		w.us2req,
		any2i5,
	)
	return func(micros int64) (raw []byte, e error) {
		length, e := u2q(micros)
		return util.Select(
			func() ([]byte, error) { return nil, e },
			func() ([]byte, error) {
				var start int32 = w.address
				var end int32 = w.address + length
				var dat []byte = w.memory.UnsafeData(w.store)
				w.buffer.Reset()
				_, _ = w.buffer.Write(dat[start:end]) // no error (or panic)
				runtime.KeepAlive(w.memory)
				return w.buffer.Bytes(), nil
			},
			nil == e,
		)()
	}
}
