(() => {
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

  // TODO
  const us2json = (micros) => {
    return JSON.stringify({
      dummy: "data",
      sample: 42,
      c: 2.99792458,
      micros: micros + "",
      complex: {
        header: ["c1", "c2"],
        data: [
          [2, 0.599],
          [3, 3.776],
          [5, 0.634],
        ],
      },
    });
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
