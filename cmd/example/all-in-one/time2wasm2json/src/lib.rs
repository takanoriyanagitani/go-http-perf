#![deny(unsafe_code)]

use std::ops::{Deref, DerefMut};
use std::sync::Mutex;

extern "C" {
    fn sin_f64(x: f64) -> f64;
    fn sin_f32(x: f32) -> f32;
}

fn _sin_f64(x: f64) -> f64 {
    #[allow(unsafe_code)]
    unsafe {
        sin_f64(x)
    }
}

fn _sin_f32(x: f32) -> f32 {
    #[allow(unsafe_code)]
    unsafe {
        sin_f32(x)
    }
}

static mut _OUT: [u8; 65536] = [0; 65536];
static _INTERNAL: Mutex<Option<Data>> = Mutex::new(None);

fn output_ref() -> &'static [u8] {
    #[allow(unsafe_code)]
    unsafe {
        &_OUT
    }
}

fn ref2ptr(r: &[u8]) -> *const u8 {
    r.as_ptr()
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn output2ptr() -> *const u8 {
    let r: &[u8] = output_ref();
    ref2ptr(r)
}

#[derive(serde::Serialize)]
struct Meta {
    base_timestamp_s: i64,
    user_name: String,
    base_current: f64,
}

impl Meta {
    fn new(seed_unixtime_us: i64) -> Self {
        Self {
            base_timestamp_s: seed_unixtime_us / 1_000_000,
            user_name: String::from("USR_0001"),
            base_current: _sin_f64(seed_unixtime_us as f64),
        }
    }
}

#[derive(serde::Serialize)]
struct Detail {
    timestamp_ms: u16,
    current: f32,
}

impl Detail {
    fn new(seed_unixtime_us: i64) -> Self {
        let timestamp_ms: u16 = (seed_unixtime_us & 0xffff) as u16;
        let ms: f32 = timestamp_ms.into();
        Self {
            timestamp_ms,
            current: _sin_f32(ms * 1e-3),
        }
    }
}

#[derive(serde::Serialize)]
struct Data {
    m: Meta,
    d: Vec<Detail>,
}

impl Data {
    fn update_by_seed(&mut self, seed_unixtime_us: i64) {
        self.m = Meta::new(seed_unixtime_us);
        self.d.clear();
        let cap: usize = self.d.capacity();
        for i in 0..cap {
            let ix: i64 = i as i64;
            let us: i64 = seed_unixtime_us + ix;
            let neo: Detail = Detail::new(us);
            self.d.push(neo);
        }
    }

    fn update_opt(o: Option<&mut Self>, seed_unixtime_us: i64) -> i32 {
        match o {
            None => -1,
            Some(ms) => {
                ms.update_by_seed(seed_unixtime_us);
                1
            }
        }
    }

    fn update_lock(l: &Mutex<Option<Self>>, seed_unixtime_us: i64) -> i32 {
        match l.lock() {
            Ok(mut guard) => {
                let md: &mut Option<Self> = guard.deref_mut();
                let om: Option<&mut Self> = md.into();
                Self::update_opt(om, seed_unixtime_us)
            }
            Err(_) => -1,
        }
    }

    fn update_static(seed_unixtime_us: i64) -> i32 {
        let mr: &Mutex<Option<Self>> = &_INTERNAL;
        Self::update_lock(mr, seed_unixtime_us)
    }

    fn to_slice(&self, buf: &mut [u8]) -> i32 {
        match serde_json_core::to_slice(self, buf) {
            Ok(sz) => sz as i32,
            Err(_) => -1,
        }
    }

    fn to_slice_opt(o: Option<&Self>, buf: &mut [u8]) -> i32 {
        match o {
            Some(s) => s.to_slice(buf),
            None => -1,
        }
    }

    fn to_slice_lock(l: &Mutex<Option<Self>>, buf: &mut [u8]) -> i32 {
        match l.lock() {
            Ok(guard) => {
                let ro: &Option<Self> = guard.deref();
                let or: Option<&Self> = ro.into();
                Self::to_slice_opt(or, buf)
            }
            Err(_) => -1,
        }
    }

    fn to_slice_static() -> i32 {
        let l: &Mutex<Option<Data>> = &_INTERNAL;
        #[allow(unsafe_code)]
        let buf: &mut [u8] = unsafe { &mut _OUT };
        Self::to_slice_lock(l, buf)
    }

    fn init() -> Self {
        Self {
            m: Meta::new(0),
            d: Vec::with_capacity(256),
        }
    }
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn init_internal() -> i32 {
    match _INTERNAL.lock() {
        Ok(mut guard) => {
            let mo: &mut Option<Data> = guard.deref_mut();
            let neo: Data = Data::init();
            mo.replace(neo);
            1
        }
        Err(_) => -1,
    }
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn unixtime2json(seed_unixtime_us: i64) -> i32 {
    Data::update_static(seed_unixtime_us);
    Data::to_slice_static()
}
