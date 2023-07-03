#!/bin/sh

which docker | fgrep --silent docker || alias docker=podman

tag=redis:7.0.11-alpine3.18

docker \
	run \
	--detach \
	--name redis \
	--publish 127.0.0.1:6379:6379 \
	"${tag}"
