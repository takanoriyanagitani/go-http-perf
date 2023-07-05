package ws

type Seed2bytes[T any] func(seed T) (serialized []byte, e error)
