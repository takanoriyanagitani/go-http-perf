from typing import List
from dataclasses import dataclass

import math
import functools

import wasmtime.loader

#import wasm_sine
import rs_sin_fast

I4TO_F5: float = 1.0 / 32768.0

def f32_sin_fast_u64(x: int)->float:
	return rs_sin_fast.f32_sin_fast_u64(x)

def f32_sin_slow_u64(x: int)->float:
	i4: int = x & 0xffff
	f4: float = float(i4)
	xf: float = f4 * I4TO_F5
	return math.sin(xf)

@dataclass
class CompareSine:
	input: int
	fast: float
	slow: float

def input2compare(input: int)->CompareSine:
	return CompareSine(
		input=input,
		fast=f32_sin_fast_u64(input),
		slow=f32_sin_slow_u64(input),
	)

inputs: List[int] = [
	0, 1, 2,
	16, 17, 18,
	128, 129, 130,
	1024, 1025, 1026,
	16384, 16385, 16386,
	32767, 32768, 32769,
	65535, 65536, 65537,
]

mapd: List[CompareSine] = list(map(input2compare, inputs))

for compare in mapd:
	diff: float = compare.fast - compare.slow
	print(dict(
		input=compare.input,
		diff=diff,
	))
