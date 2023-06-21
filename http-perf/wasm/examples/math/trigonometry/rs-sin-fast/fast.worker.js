(() => {
  const url2response = async (url) => {
    const response = await fetch(url);
    return response;
  };

  const response2wasm = async (response) =>
    WebAssembly.instantiateStreaming(response, {});

  const composeAsync = (f, g) => {
    return async (t) => {
      const u = await f(t);
      const v = await g(u);
      return v;
    };
  };

  const url2wasm = composeAsync(
    url2response,
    response2wasm,
  );

  const wasmUrl = "/rs_sin_fast.wasm";

  self.onmessage = (evt) => {
    const { data } = evt;
    const { calls } = data;
    return Promise.resolve(wasmUrl)
      .then(url2wasm)
      .then((result) => {
        const {
          module,
          instance,
        } = result;
        const { exports } = instance;
        const { f32_sin_fast_u64 } = exports;
        const start = performance.now();
        let sum = 0.0;
        for (let i = 0n; i < calls; i++) {
          sum += f32_sin_fast_u64(i);
        }
        const end = performance.now();
        const diff = end - start;
        self.postMessage({
          ms: diff,
          sum,
        });
      });
  };
})();
