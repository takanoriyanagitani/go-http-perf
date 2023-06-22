const PI: f32 = std::f32::consts::PI;

const B: f32 = 4.0 / PI;

const A1: f32 = (2.0 / PI) * (2.0 / PI);

const A: f32 = A1;
const ANII: f32 = -2.0 * A;

const I4TO_F5: f32 = 1.0 / 32768.0;

#[no_mangle]
pub extern "C" fn i4f5(i4: i16) -> f32 {
    let f: f32 = i4.into();
    f * I4TO_F5
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
