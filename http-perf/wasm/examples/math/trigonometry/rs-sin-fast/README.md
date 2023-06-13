# Sine-like function

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

| i16    | f32       |
|:------:|:---------:|
|      0 |  0.000    |
|  32767 |  0.999... |
| -32768 | -1.000    |

f32(i) = i/32768.0
