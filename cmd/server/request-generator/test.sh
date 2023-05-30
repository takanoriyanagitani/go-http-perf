#!/bin/sh

url=localhost:53080/now2req

curl --silent "${url}" |
  python3 proto2json.py |
  jq .
