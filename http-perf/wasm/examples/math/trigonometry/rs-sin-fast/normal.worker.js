(() => {
  self.onmessage = (evt) => {
    const { data } = evt;
    const { calls } = data;
    return Promise.resolve()
      .then((_) => {
        const f32_sin_slow_dummy = Math.sin;
        const start = performance.now();
        let sum = 0.0;
        for (let i = 0; i < calls; i++) {
          sum += f32_sin_slow_dummy(i);
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
