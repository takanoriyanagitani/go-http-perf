#!/bin/sh

cargo \
	build \
	--profile release-wasm \
	--target wasm32-unknown-unknown
