# Trigomonetric functions

## Sine(x: 0 ... pi/2)

```
f(x) = ax^3 + bx^2 + cx + d
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
