const PI: f32 = std::f32::consts::PI;

const B: f32 = 4.0 / PI;

const A0: f32 = -(2.0 / PI) * (2.0 / PI);
const B0: f32 = B;
//const C0: f32 = 0.0;

const A1: f32 = (2.0 / PI) * (2.0 / PI);
const B1: f32 = B;
//const C1: f32 = 0.0;

const A: f32 = A1;
const ANII: f32 = -2.0 * A;

const I4TO_F5: f32 = 1.0 / 32768.0;

#[no_mangle]
pub extern "C" fn f32_sin_fast_0_180(x: f32) -> f32 {
    //A0 * x * x + B0 * x + C0
    //A0 * x * x + B0 * x
    f32_sin_fast_ab(x, A0, B0)
}

#[no_mangle]
pub extern "C" fn f32_sin_fast_m180_0(x: f32) -> f32 {
    //A1 * x * x + B1 * x + C1
    //A1 * x * x + B1 * x
    f32_sin_fast_ab(x, A1, B1)
}

#[no_mangle]
pub extern "C" fn f32_sin_fast_ab(x: f32, a: f32, b: f32) -> f32 {
    a * x * x + b * x
}

#[no_mangle]
pub extern "C" fn f32_sin_fast_a(x: f32, a: f32) -> f32 {
    f32_sin_fast_ab(x, a, B)
}

#[no_mangle]
pub extern "C" fn i4f5(i4: i16) -> f32 {
    let f: f32 = i4.into();
    f * I4TO_F5
}

#[no_mangle]
pub const extern "C" fn u4i4(u4: u16) -> i16 {
    u4 as i16
}

#[no_mangle]
pub const extern "C" fn u5u4(u5: u32) -> u16 {
    (u5 & 0xffff) as u16
}

#[no_mangle]
pub const extern "C" fn u6u4(u6: u64) -> u16 {
    (u6 & 0xffff) as u16
}

#[no_mangle]
pub extern "C" fn b2f5(b: bool) -> f32 {
    match b {
        true => 1.0,
        false => 0.0,
    }
}

#[no_mangle]
pub extern "C" fn f32_sin_fast_i32(i: i32) -> f32 {
    let i4: i16 = (i & 0xffff) as i16;
    let pos: bool = i4.is_positive();
    let pf: f32 = b2f5(pos);
    let a: f32 = A + ANII * pf;
    let b: f32 = B;
    let x: f32 = i4f5(i4); // -1.0 <= x < 1.0
    let xp: f32 = PI * x; // - pi <= xp < pi
    a * xp * xp + b * xp
}

#[no_mangle]
pub extern "C" fn f32_sin_fast_u64(u: u64) -> f32 {
    let i: i32 = (u & 0xffff) as i32;
    f32_sin_fast_i32(i)
}
