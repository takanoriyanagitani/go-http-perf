#[no_mangle]
pub extern "C" fn f64_sin(x: f64) -> f64 {
    x.sin()
}
