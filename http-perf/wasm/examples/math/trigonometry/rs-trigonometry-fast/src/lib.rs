const PI: f32 = std::f32::consts::PI;

const A: f32 = 4.0 * (PI - 4.0) / PI / PI / PI;
const B: f32 = 4.0 * (3.0 - PI) / PI / PI;
const C: f32 = 1.0;

// f(x) = 4(pi-4)(x/pi)^3 + 4(3-pi)(x/pi)^2 + x ~ sin(x)|x: 0 ... pi/2
#[no_mangle]
pub extern "C" fn f32_sin_0_90(x: f32) -> f32 {
    A * x * x * x + B * x * x + C * x
}
