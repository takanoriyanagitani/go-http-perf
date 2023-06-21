# Sine-like function

## Benchmarks

### browser / macOS Ventura(Apple M2)

| type                | calls | elapsed | calls / s |
|:-------------------:|:-----:|:-------:|:---------:|
| firefox / sine-like |  65K  |   4 ms  |  16 M     |
| firefox / Math.sin  |  65K  |   1 ms  |  66 M     |
| firefox / sine-like | 131K  |   3 ms  |  44 M     |
| firefox / Math.sin  | 131K  |   3 ms  |  44 M     |
| firefox / sine-like |   1M  |  17 ms  |  62 M     |
| firefox / Math.sin  |   1M  |  14 ms  |  75 M     |
| firefox / sine-like |  16M  | 120 ms  | 140 M     |
| firefox / Math.sin  |  16M  | 137 ms  | 122 M     |
| firefox / sine-like | 167M  | 1.04 s  | 161 M     |
| firefox / Math.sin  | 167M  | 1.05 s  | 160 M     |
| safari  / sine-like |  65K  |  14 ms  |   5 M     |
| safari  / Math.sin  |  65K  |   3 ms  |  22 M     |
| safari  / sine-like | 131K  |  19 ms  |   7 M     |
| safari  / Math.sin  | 131K  |   5 ms  |  26 M     |
| safari  / sine-like |   1M  |  76 ms  |  14 M     |
| safari  / Math.sin  |   1M  |  17 ms  |  62 M     |
| safari  / sine-like |  16M  | 731 ms  |  23 M     |
| safari  / Math.sin  |  16M  | 102 ms  | 164 M     |
| safari  / sine-like | 167M  | 7.92 s  |  21 M     |
| safari  / Math.sin  | 167M  | 970 ms  | 173 M     |
| chrome  / sine-like |  65K  |  17 ms  |   4 M     |
| chrome  / Math.sin  |  65K  |   5 ms  |  13 M     |
| chrome  / sine-like | 131K  |  11 ms  |  12 M     |
| chrome  / Math.sin  | 131K  |   5 ms  |  26 M     |
| chrome  / sine-like |   1M  |  60 ms  |  17 M     |
| chrome  / Math.sin  |   1M  |  25 ms  |  42 M     |
| chrome  / sine-like |  16M  | 515 ms  |  32 M     |
| chrome  / Math.sin  |  16M  | 128 ms  | 131 M     |
| chrome  / sine-like | 167M  | 5.02 s  |  33 M     |
| chrome  / Math.sin  | 167M  | 2.62 s  |  64 M     |

### node v18.16.0 / macOS Ventura(Apple M2)

- "fast" version: up to 32 M calls / s @ Apple M2
- "slow" version: up to 59 M calls / s @ Apple M2

### node v18.14.2 / linux(Core i7-8700 3.2 GHz)

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
