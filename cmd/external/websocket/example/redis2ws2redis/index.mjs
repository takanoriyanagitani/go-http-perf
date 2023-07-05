(()=>{

	const ws = new WebSocket("ws://localhost:7058/redis2ws2redis");
	ws.binaryType = "arraybuffer";

	const onEmptyNew = () => {
		return () => {
			ws.send("hw");
		};
	};

	const onDataNew = data => {
		return () => {
			console.info(data);
			ws.send("hw");
		};
	};

	ws.onmessage = evt => {
		const { data } = evt;
		const { byteLength } = data;
		const f = 0 == byteLength ? onEmptyNew() : onDataNew(data);
		f();
	};

})()
