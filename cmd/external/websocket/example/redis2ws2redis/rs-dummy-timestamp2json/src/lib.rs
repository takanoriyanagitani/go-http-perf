#![deny(unsafe_code)]

use std::ops::{Deref, DerefMut};
use std::sync::Mutex;

static mut _OUTPUT: [u8; 65536] = [0; 65536];
static _BUF: Mutex<Option<TestData>> = Mutex::new(None);

const U4SCALE: f32 = 2.0 * std::f32::consts::PI / 65536.0;

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn init() {
    _BUF.lock()
        .map(|mut guard| {
            let rm: &mut Option<TestData> = guard.deref_mut();
            rm.replace(TestData::init());
        })
        .unwrap_or(())
}

fn buf2slice(target: &mut [u8]) -> Result<usize, String> {
    match _BUF.lock() {
        Ok(guard) => {
            let rm: &Option<TestData> = guard.deref();
            let ortd: Option<&TestData> = rm.into();
            TestData::opt2slice(ortd, target)
        }
        _ => Ok(0),
    }
}

fn output_ref() -> &'static [u8] {
    #[allow(unsafe_code)]
    unsafe {
        &_OUTPUT
    }
}

fn ref2ptr(r: &[u8]) -> *const u8 {
    r.as_ptr()
}

#[derive(serde::Serialize)]
struct TestWave {
    timestamp_ms: u16,
    current: f32,
    voltage: f32,
}

impl TestWave {
    fn from_seed(timestamp_us: i64) -> Self {
        let u4: u16 = (timestamp_us & 65535) as u16;
        let f5: f32 = u4.into();
        let x: f32 = U4SCALE * f5;
        Self {
            timestamp_ms: u4,
            current: x.sin(),
            voltage: x.cos(),
        }
    }
}

#[derive(serde::Serialize)]
struct TestData {
    timestamp_s: i64,
    version: i64,
    wave: Vec<TestWave>,
    name: String,
    status: i64,
}

impl TestData {
    fn to_slice(&self, buf: &mut [u8]) -> Result<usize, String> {
        serde_json_core::to_slice(self, buf).map_err(|e| format!("Unable to serialize: {e}"))
    }

    fn update(&mut self, timestamp_us: i64) {
        self.timestamp_s = timestamp_us / 1_000_000;
        self.version = 310;
        self.name = String::from("TEST1");
        self.status = timestamp_us & 0x3;

        self.wave.clear();
        for i in 0..256 {
            let seed: i64 = timestamp_us + i;
            let sample: TestWave = TestWave::from_seed(seed);
            self.wave.push(sample);
        }
    }

    fn init() -> Self {
        Self {
            timestamp_s: 0,
            version: -1,
            wave: Vec::with_capacity(256),
            name: "".into(),
            status: -1,
        }
    }

    fn update_opt(otd: Option<&mut Self>, timestamp_us: i64) {
        match otd {
            None => {}
            Some(td) => td.update(timestamp_us),
        }
    }

    fn opt2slice(otd: Option<&Self>, buf: &mut [u8]) -> Result<usize, String> {
        match otd {
            None => Ok(0),
            Some(td) => td.to_slice(buf),
        }
    }
}

fn update_buf(timestamp_us: i64) {
    _BUF.lock()
        .map(|mut guard| {
            let rm: &mut Option<TestData> = guard.deref_mut();
            let ortd: Option<&mut TestData> = rm.into();
            TestData::update_opt(ortd, timestamp_us)
        })
        .unwrap_or(())
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn output_ptr() -> *const u8 {
    let r: &[u8] = output_ref();
    ref2ptr(r)
}

fn timestamp2buf(timestamp_us: i64, buf: &mut [u8]) -> Result<usize, String> {
    update_buf(timestamp_us);
    buf2slice(buf)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn timestamp2json(timestamp_us: i64) -> usize {
    #[allow(unsafe_code)]
    let buf: &mut [u8] = unsafe { &mut _OUTPUT };
    match timestamp2buf(timestamp_us, buf) {
        Ok(usz) => usz,
        _ => 0,
    }
}
