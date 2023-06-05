#[no_mangle]
pub extern "C" fn f64_sin(x: f64) -> f64 {
    x.sin()
}

#[no_mangle]
pub extern "C" fn f32_sin(x: f32) -> f32 {
    x.sin()
}

#[no_mangle]
pub extern "C" fn f64_atan(x: f64) -> f64 {
    x.atan()
}
