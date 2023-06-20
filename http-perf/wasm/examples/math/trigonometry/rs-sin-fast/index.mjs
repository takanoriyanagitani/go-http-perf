(() => {
  const fast = document.getElementById("fast");
  const slow = document.getElementById("slow");

  const fms = document.getElementById("fast-ms");
  const sms = document.getElementById("slow-ms");

  const sel = document.getElementById("calls");

  const fw = new Worker("fast.worker.js");
  fw.onmessage = (evt) => {
    const { data } = evt;
    const { ms } = data;
    fms.textContent = ms;
  };
  fast.onclick = (_) => fw.postMessage({ calls: sel.value - 0 });

  const sw = new Worker("normal.worker.js");
  sw.onmessage = (evt) => {
    const { data } = evt;
    const { ms } = data;
    sms.textContent = ms;
  };
  slow.onclick = (_) => sw.postMessage({ calls: sel.value - 0 });
})();
