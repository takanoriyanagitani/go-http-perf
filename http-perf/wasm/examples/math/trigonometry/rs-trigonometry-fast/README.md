# Trigomonetric functions

## Sine(x: 0 ... 2 pi)

| x range         | f(x) ~ sin(x)                                             |
|:---------------:|:---------------------------------------------------------:|
| 0     ... 0.5pi | 4(pi-4)(  x     /pi)^3 + 4(3-pi)(  x     /pi)^2 + x       |
| 0.5pi ...    pi | 4(pi-4)((-x+ pi)/pi)^3 + 4(3-pi)((-x+ pi)/pi)^2 - x +  pi |
|    pi ... 1.5pi | 4(pi-4)((-x+ pi)/pi)^3 - 4(3-pi)((-x+ pi)/pi)^2 - x +  pi |
| 1.5pi ...   2pi | 4(pi-4)(( x-2pi)/pi)^3 - 4(3-pi)((-x+2pi)/pi)^2 + x - 2pi |

f(x) = ax^3 + bx^2 + cx + d ~ sin(x)

| x range         | a              | b              | c             | d         |
|:---------------:|:--------------:|:--------------:|:-------------:|:---------:|
| 0     ... 0.5pi |  (4pi-16)/pi^3 | (12-4pi)/pi^2  | 1             | 0         |
| 0.5pi ...    pi | -(4pi-16)/pi^3 | (8pi-36)/pi^2  | (24-5pi)/pi   | pi - 4    |
|    pi ... 1.5pi | -(4pi-16)/pi^3 | (16pi-60)/pi^2 | (72-21pi)/pi  | 9pi - 28  |
| 1.5pi ...   2pi |  (4pi-16)/pi^3 | (84-20pi)/pi^2 | (33pi-144)/pi | 80 - 18pi |

## Sine(x: 0 ... pi/2)

```
f(x) = ax^3 + bx^2 + cx + d ~ sin(x)
f(0) = sin(0) = 0 = d

f(x) = ax^3 + bx^2 + cx

f'(x) = 3ax^2 + 2bx + c
f'(0) = cos(0) = 1 = c

f'(x) = 3ax^2 + 2bx + 1
f(x)  = ax^3 + bx^2 + x

f(pi/2) = sin(pi/2) = 1 = a pi^3/8 + b pi^2/4 + pi/2
8 = a pi^3 + 2b pi^2 + 4 pi
16 = 2a pi^3 + 4b pi^2 + 8 pi
16 - 2a pi^3 - 8 pi = 4b pi^2

f'(pi/2) = cos(pi/2) = 0 = 3a pi^2 / 4 + 2b pi/2 + 1
0 = 3a pi^2 + 4b pi + 4
0 = 3a pi^3 + 4b pi^2 + 4pi
  = 3a pi^3 + 4pi + (16 - 2a pi^3 - 8 pi)
  = 3a pi^3 + 4 pi + 16 - 2a pi^3 - 8 pi
  = a pi^3 - 4 pi + 16
4 pi - 16 = a pi^3
a = (4 pi - 16)/pi^3
  = 4(pi-4)/pi^3
4b pi^2 = 16 - 8pi - 2a pi^3
2b pi^2 = 8 - 4 pi - a pi^3
 = 8 - 4 pi - 4(pi-4) pi^3/pi^3
 = 8 - 4 pi - 4(pi-4)
 = 8 - 4 pi - 4 pi + 16
 = 24 - 8 pi
 = 8(3 - pi)
b pi^2 = 4(3-pi)
b = 4(3-pi)/pi^2

f(x) = 4(pi-4)(x/pi)^3 + 4(3-pi)(x/pi)^2 + x

```
