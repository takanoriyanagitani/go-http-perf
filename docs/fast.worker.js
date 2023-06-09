(() => {
  const compose = (f, g) => {
    return (t) => {
      const u = f(t);
      const v = g(u);
      return v;
    };
  };

  const composeAsync = (f, g) => {
    return async (t) => {
      return Promise.resolve(t)
        .then(f)
        .then(g);
    };
  };

  const url2response = async (url) => {
    const response = await fetch(url);
    return response;
  };

  const response2buf = async (res) => res.arrayBuffer();
  const buf2wasm = (buf) => WebAssembly.instantiate(buf, {});
  const response2buf2wasm = composeAsync(
    response2buf,
    buf2wasm,
  );

  const response2wasm = async (response) =>
    WebAssembly.instantiateStreaming(response, {});

  const url2wasm = composeAsync(
    url2response,
    response2buf2wasm,
  );

  const wasmUrl = "/go-http-perf/rs_sin_fast.wasm";

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
