import { readFile } from "node:fs/promises";

const wasmName = "./rs_sin_fast.wasm";

const composeAsync = function (f, g) {
  return async function (t) {
    const u = await f(t);
    const v = await g(u);
    return v;
  };
};

const filename2buffer = async function (filename) {
  return readFile(
    filename,
    {
      encoding: null,
      flag: "r",
    },
  );
};

const buffer2arrayBuffer = (buffer) => buffer.buffer;

const filename2bytes = composeAsync(
  filename2buffer,
  buffer2arrayBuffer,
);

const wasmBytes2wasm = function (bytes) {
  return WebAssembly.instantiate(bytes, {});
};

const filename2wasm = composeAsync(
  filename2bytes,
  wasmBytes2wasm,
);

const module2instance = function (module) {
  return new WebAssembly.Instance(module, {});
};

const compare_dummy_sin = (i) => {
  const o = i / 32768.0;
  const x = o * Math.PI;
  return Math.sin(x);
};

const main = async function () {
  const wasm = await filename2wasm(wasmName);
  const {
    instance,
    module,
  } = wasm;
  const { exports } = instance;
  const { f32_sin_fast_u64 } = exports;
  let sum = 0.0;
  if (false) {
    for (let i = 0n; i < 16777216n; i++) { // 32 M calls / s @ Apple M2
      const f = f32_sin_fast_u64(i);
      sum += f;
    }
  } else {
    for (let i = 0; i < 16777216; i++) { // 59 M calls / s @ Apple M2
      const f = compare_dummy_sin(i);
      sum += f;
    }
  }
  console.info(sum);
  return;
};

Promise.resolve()
  .then(main)
  .catch(console.error);
