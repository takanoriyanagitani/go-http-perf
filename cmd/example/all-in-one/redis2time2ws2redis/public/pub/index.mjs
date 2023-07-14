(() => {

	const { host } = window.location;
	const wpath = "ws://" + host + "/redis2ws2redis";
	console.info(wpath)
	const ws = new WebSocket(wpath);
	ws.binaryType = "arraybuffer";

	const compose = (f, g) => {
		return t => {
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
		})
	};

	const data2us = data => {
		const v = new DataView(data);
		const us = v.getBigInt64(0);
		return us;
	};

	const data2us2json = compose(
		data2us,
		us2json,
	);

	ws.onmessage = evt => {
		const { data } = evt;
		const { byteLength } = data;
		if(byteLength < 8) return ws.send("");
		const j = data2us2json(data);
		ws.send(j);
	};

})()
