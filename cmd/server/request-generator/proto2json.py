import sys
import functools
import json
import operator

import request_pb2

reducer = lambda state, f: f(state)

curry = lambda f: lambda x: lambda y: f(x,y)

def bytes2requestNew(buf, data):
	buf.Clear()
	buf.ParseFromString(data)
	return buf

bytes2request = curry(bytes2requestNew)(request_pb2.Request())

functools.reduce(
	reducer,
	[
		lambda rdr: rdr.read(),
		bytes2request,
		operator.attrgetter("body"),
		json.loads,
		json.dumps,
		print,
	],
	sys.stdin.buffer,
)
