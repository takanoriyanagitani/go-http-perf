import functools
import pathlib
import json

from wasmtime import Store, Module, Instance, Func

import request_pb2

def bytes2reqNew(buf, src):
	buf.Clear()
	buf.ParseFromString(src)
	return buf

def curry(f):
	return lambda x: lambda y: f(x,y)

def compose(f, g):
	return lambda x: functools.reduce(
		lambda state, f: f(state),
		[f, g],
		x,
	)

bytes2req = curry(bytes2reqNew)(request_pb2.Request())
bytearray2req = compose(bytes, bytes2req)

filepath2bytes = lambda p: pathlib.Path(p).read_bytes()

store = Store()
wasm = filepath2bytes("./rs_time2json.wasm")
module = Module(store.engine, wasm)
instance = Instance(store, module, [])

addr = instance.exports(store)["addr"](store)
time2req = instance.exports(store)["time2req"]
end = time2req(store, 3776)
memory = instance.exports(store)["memory"]
copied = memory.read(
	store,
	addr,
	addr+end,
)
req = bytearray2req(copied)
print(json.loads(req.body))
