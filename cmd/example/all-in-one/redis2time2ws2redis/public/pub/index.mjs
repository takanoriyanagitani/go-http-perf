(async () => {
  const _sent = document.getElementById("sent");
  const _recv = document.getElementById("recv");

  const incText = (prev) => {
    const i = prev - 0;
    const next = i + 1;
    return next;
  };

  const incContent = (ele) => {
    const next = incText(ele.textContent);
    ele.textContent = next;
  };

  const incSent = (_) => incContent(_sent);
  const incRecv = (_) => incContent(_recv);

  const { host } = window.location;
  const wpath = "ws://" + host + "/redis2ws2redis";
  console.info(wpath);
  const ws = new WebSocket(wpath);
  ws.binaryType = "arraybuffer";

  const compose = (f, g) => {
    return (t) => {
      const u = f(t);
      const v = g(u);
      return v;
    };
  };

  const wasmRes = await fetch("/pub/time2wasm2json.wasm");
  const wasmBuf = await wasmRes.arrayBuffer();
  const { module, instance } = await WebAssembly.instantiate(
    wasmBuf,
    {
      env: {
        sin_f64: Math.sin,
        sin_f32: Math.sin,
      },
    },
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
  const { buffer } = memory;
  const decoder = new TextDecoder("utf-8");

  const us2json = (micros) => {
    const sz = unixtime2json(BigInt(Date.now() * 1000));
    const json = buffer.slice(baseAddr, baseAddr + sz);
    const decoded = decoder.decode(json);
    return decoded;
  };

  const data2us = (data) => {
    const v = new DataView(data);
    const us = v.getBigInt64(0);
    return us;
  };

  const data2us2json = compose(
    data2us,
    us2json,
  );

  const send2ws = (dat) => {
    ws.send(dat);
    incSent();
  };

  ws.onmessage = (evt) => {
    incRecv();
    const { data } = evt;
    const { byteLength } = data;
    if (byteLength < 8) return send2ws("");
    const j = data2us2json(data);
    send2ws(j);
  };
})();
