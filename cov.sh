#!/bin/sh

out=./coverage.txt
vis=./coverate.html

go \
	test \
	-race \
	-covermode=atomic \
	-coverprofile="${out}" \
	./...

go \
	tool \
	cover \
	-html="${out}" \
	-o "${vis}"
