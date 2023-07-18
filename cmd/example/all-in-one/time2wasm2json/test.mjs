import { readFile } from "node:fs/promises";

const main = async () => {
  const wasm = await readFile(
    "./time2wasm2json.wasm",
    {
      encoding: null,
      flag: "r",
    },
  );
  const iobj = {
    env: {
      sin_f64: Math.sin,
      sin_f32: Math.sin,
    },
  };
  const { instance, module } = await WebAssembly.instantiate(
    wasm,
    iobj,
  );
  const { exports } = instance;
  const {
    memory,
    init_internal,
    output2ptr,
    unixtime2json,
  } = exports;
  init_internal();
  const baseAddr = output2ptr();
  const sz = unixtime2json(BigInt(Date.now() * 1000));
  const { buffer } = memory;
  const json = buffer.slice(baseAddr, baseAddr + sz);
  const decoder = new TextDecoder("utf-8");
  const decoded = decoder.decode(json);
  const parsed = JSON.parse(decoded);
  return JSON.stringify(parsed);
};

Promise.resolve()
  .then(main)
  .then(console.info)
  .catch(console.error);
