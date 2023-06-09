/*

f(x) = a(x-b)^2+c ~ sin(x)
b = pi/2
f(pi/2) = c = 1
f(x) = a(x-b)^2 + 1
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
*/

const PI: f32 = std::f32::consts::PI;

const A0: f32 = -(2.0 / PI) * (2.0 / PI);
const B0: f32 = 4.0 / PI;
const C0: f32 = 0.0;

#[no_mangle]
pub extern "C" fn f32_sin_fast_0_180(x: f32) -> f32 {
    A0 * x * x + B0 * x + C0
}
