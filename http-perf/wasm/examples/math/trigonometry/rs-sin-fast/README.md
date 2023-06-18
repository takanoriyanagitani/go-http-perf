# Sine-like function

## Benchmark

- Intel Core i7-8700 3.2 GHz
- node v18.14.2 / linux

|      type | calls | user | sys  | cpu% | total | calls/s | ratio |
| --------: | :---: | :--: | :--: | :--: | :---: | :-----: | :---: |
|  Math.sin |  65K  | 0.04 | 0.00 | 103  | 0.038 | 1.72 M  | 100%  |
|  Math.sin | 131K  | 0.03 | 0.01 | 102  | 0.039 | 3.36 M  | 195%  |
|  Math.sin |  1M   | 0.05 | 0.00 | 101  | 0.055 | 19.07 M | 1109% |
|  Math.sin |  16M  | 0.34 | 0.01 | 100  | 0.346 | 48.49 M | 2819% |
| sine-like |  65K  | 0.04 | 0.00 | 102  | 0.043 | 1.52 M  |  88%  |
| sine-like | 131K  | 0.04 | 0.01 | 101  | 0.045 | 2.91 M  | 169%  |
| sine-like |  1M   | 0.09 | 0.01 | 100  | 0.100 | 10.49 M | 610%  |
| sine-like |  16M  | 1.77 | 0.02 | 100  | 1.788 | 9.38 M  | 545%  |

## theta: -pi .. pi

```
f(x) = a x^2 + bx
b = 4/pi
am: a =  (2/pi)^2 (-pi/2 <= x <    0)
ap: a = -(2/pi)^2 (    0 <= x < pi/2)
      = -am
sgn: 0.0 (x < 0), 1.0 (0 <= x)

a = ap sgn(x) + (1-sgn(x)) am
  = -am sgn(x) + am - am sgn(x)
  = am(1 - 2 sgn(x))
  = am - 2 am sgn(x)
anii: -2 am
a = am + anii sgn(x)
```

## theta: 0 .. pi

```
f(x) = a(x-b)^2+c ~ sin(x)
b = pi/2
f(pi/2) = c = 1
f(x) = a(x-0.5pi)^2 + 1
f(0) = ab^2 + 1 = 0
ab^2 = -1
a = -1/b^2
  = -1/(0.5pi)^2
  = -1/0.25pi^2
  = -4/pi^2
  = -(2/pi)^2
f(x) = -(2/pi)^2(x-pi/2)^2+1
     = -(2x-pi)^2/pi^2 + 1
     = -((2x-pi)/pi)^2 + 1
     = -(2/pi)^2 x^2 + (4/pi) x
```

## theta: -pi .. 0

```
g(x) = (2/pi)^2 (x+pi)^2 - (4/pi)(x+pi) ~ sin(x)
     = (2/pi)^2 x^2 + (4/pi) x
```

## i16 -> f32

|  i16   |   f32    |
| :----: | :------: |
|   0    |  0.000   |
| 32767  | 0.999... |
| -32768 |  -1.000  |

f32(i) = i/32768.0
