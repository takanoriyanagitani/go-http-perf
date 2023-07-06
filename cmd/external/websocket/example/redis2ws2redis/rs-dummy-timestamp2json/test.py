from typing import Callable

import sys

from wasmtime import Store, Module, Instance, Memory

store: Store = Store()
module: Module = Module.from_file(store.engine, "./rs_dummy_timestamp2json.wasm")
instance: Instance = Instance(store, module, [])

init = instance.exports(store)["init"]
out  = instance.exports(store)["output_ptr"](store)
mem  = instance.exports(store)["memory"]
t2j  = instance.exports(store)["timestamp2json"]

init(store)
jlen = t2j(store, 43)

copied = mem.read(store, out, out+jlen)
sys.stdout.buffer.write(copied)
