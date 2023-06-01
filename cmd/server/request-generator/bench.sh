#!/bin/sh

single(){
	# mem: 23 ~ 37 MB
	# cpu: ~250 %
	ENV_WASM_PATH=./rs_time2json.wasm \
		ENV_USE_POOL=no \
		./request-generator \
		&
}

multi(){
	# slower than single threaded(mutex guarded) version
	# mem: 23 ~ 167 MB
	# cpu: ~560 %
	ENV_WASM_PATH=./rs_time2json.wasm \
		ENV_USE_POOL=pool \
		./request-generator \
		&
}

#multi
single

sleep 5

bombardier \
	localhost:53080/now2req
