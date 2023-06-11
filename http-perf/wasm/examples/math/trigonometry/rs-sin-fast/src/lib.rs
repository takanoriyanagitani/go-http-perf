const PI: f32 = std::f32::consts::PI;

const A0: f32 = -(2.0 / PI) * (2.0 / PI);
const B0: f32 = 4.0 / PI;
//const C0: f32 = 0.0;

const A1: f32 = (2.0 / PI) * (2.0 / PI);
const B1: f32 = 4.0 / PI;
//const C1: f32 = 0.0;

#[no_mangle]
pub extern "C" fn f32_sin_fast_0_180(x: f32) -> f32 {
    //A0 * x * x + B0 * x + C0
    A0 * x * x + B0 * x
}

#[no_mangle]
pub extern "C" fn f32_sin_fast_m180_0(x: f32) -> f32 {
    //A1 * x * x + B1 * x + C1
    A1 * x * x + B1 * x
}
