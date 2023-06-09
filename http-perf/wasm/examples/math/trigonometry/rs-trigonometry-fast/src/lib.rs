const PI: f32 = std::f32::consts::PI;

const A: f32 = 4.0 * (PI - 4.0) / PI / PI / PI;
const B: f32 = 4.0 * (3.0 - PI) / PI / PI;
const C: f32 = 1.0;

const M16: f32 = PI / 32768.0;

const A0: f32 = (4.0 * PI - 16.0) / PI / PI / PI;
const A1: f32 = -(4.0 * PI - 16.0) / PI / PI / PI;
const A2: f32 = -(4.0 * PI - 16.0) / PI / PI / PI;
const A3: f32 = (4.0 * PI - 16.0) / PI / PI / PI;

const B0: f32 = (12.0 - 4.0 * PI) / PI / PI;
const B1: f32 = (-36.0 + 8.0 * PI) / PI / PI;
const B2: f32 = (-60.0 + 16.0 * PI) / PI / PI;
const B3: f32 = (84.0 - 20.0 * PI) / PI / PI;

const C0: f32 = 1.0;
const C1: f32 = (24.0 - 5.0 * PI) / PI;
const C2: f32 = (72.0 - 21.0 * PI) / PI;
const C3: f32 = (-144.0 + 33.0 * PI) / PI;

const D0: f32 = 0.0;
const D1: f32 = -4.0 + 1.0 * PI;
const D2: f32 = -28.0 + 9.0 * PI;
const D3: f32 = 80.0 - 18.0 * PI;

// f(x) = 4(pi-4)(x/pi)^3 + 4(3-pi)(x/pi)^2 + x ~ sin(x)|x: 0 ... pi/2
#[no_mangle]
pub extern "C" fn f32_sin_0_90(x: f32) -> f32 {
    A * x * x * x + B * x * x + C * x
}

#[no_mangle]
pub extern "C" fn u5u4(u5: u32) -> u16 {
    u5 as u16
}

#[no_mangle]
pub extern "C" fn u4x(u4: u16) -> f32 {
    M16 * (u4 as f32)
}

#[no_mangle]
pub extern "C" fn u5x(u5: u32) -> f32 {
    M16 * ((u5 as u16) as f32)
}

#[no_mangle]
pub extern "C" fn u5u3(u5: u32) -> u8 {
    let u4: u16 = u5 as u16;
    let b2: u8 = (u4 >> 14) as u8;
    b2
}

#[no_mangle]
pub extern "C" fn is0(u5: u32) -> bool {
    0 == u5u3(u5)
}

#[no_mangle]
pub extern "C" fn is1(u5: u32) -> bool {
    1 == u5u3(u5)
}

#[no_mangle]
pub extern "C" fn is2(u5: u32) -> bool {
    2 == u5u3(u5)
}

#[no_mangle]
pub extern "C" fn is3(u5: u32) -> bool {
    3 == u5u3(u5)
}

#[no_mangle]
pub extern "C" fn is0f(u5: u32) -> f32 {
    is0(u5).into()
}

#[no_mangle]
pub extern "C" fn is1f(u5: u32) -> f32 {
    is1(u5).into()
}

#[no_mangle]
pub extern "C" fn is2f(u5: u32) -> f32 {
    is2(u5).into()
}

#[no_mangle]
pub extern "C" fn is3f(u5: u32) -> f32 {
    is3(u5).into()
}

#[no_mangle]
pub extern "C" fn au5(u5: u32) -> f32 {
    let a0: f32 = A0 * is0f(u5);
    let a1: f32 = A1 * is1f(u5);
    let a2: f32 = A2 * is2f(u5);
    let a3: f32 = A3 * is3f(u5);

    a0 + a1 + a2 + a3
}

#[no_mangle]
pub extern "C" fn bu5(u5: u32) -> f32 {
    let b0: f32 = B0 * is0f(u5);
    let b1: f32 = B1 * is1f(u5);
    let b2: f32 = B2 * is2f(u5);
    let b3: f32 = B3 * is3f(u5);

    b0 + b1 + b2 + b3
}

#[no_mangle]
pub extern "C" fn cu5(u5: u32) -> f32 {
    let c0: f32 = C0 * is0f(u5);
    let c1: f32 = C1 * is1f(u5);
    let c2: f32 = C2 * is2f(u5);
    let c3: f32 = C3 * is3f(u5);

    c0 + c1 + c2 + c3
}

#[no_mangle]
pub extern "C" fn du5(u5: u32) -> f32 {
    let d0: f32 = D0 * is0f(u5);
    let d1: f32 = D1 * is1f(u5);
    let d2: f32 = D2 * is2f(u5);
    let d3: f32 = D3 * is3f(u5);

    d0 + d1 + d2 + d3
}

#[no_mangle]
pub extern "C" fn du5_select(u5: u32) -> f32 {
    let d0: f32 = is0(u5).then_some(D0).unwrap_or_default();
    let d1: f32 = is1(u5).then_some(D1).unwrap_or_default();
    let d2: f32 = is2(u5).then_some(D2).unwrap_or_default();
    let d3: f32 = is3(u5).then_some(D3).unwrap_or_default();

    d0 + d1 + d2 + d3
}

#[no_mangle]
pub extern "C" fn u32_sin(u: u32) -> f32 {
    let a: f32 = au5(u);
    let b: f32 = bu5(u);
    let c: f32 = cu5(u);
    let d: f32 = du5(u);

    let x: f32 = u5x(u);

    a * x * x * x + b * x * x + c * x + d
}
