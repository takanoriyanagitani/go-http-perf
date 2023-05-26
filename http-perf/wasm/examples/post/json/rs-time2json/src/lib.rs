#![deny(unsafe_code)]

use bytes::BufMut;
use prost::Message;

extern crate alloc;

use alloc::collections::btree_map::BTreeMap;

use crate::request::{HeaderContent, Request};

static mut _OUT_PB: [u8; 65536] = [0; 65536];

pub mod request {
    include!(concat!(env!("OUT_DIR"), "/request.rs"));
}

enum Error {
    InsufficientCapacity(usize, usize),
}

fn msg2buf<M, B>(msg: &M, buf: &mut B) -> Result<(), Error>
where
    M: Message,
    B: BufMut,
{
    msg.encode(buf)
        .map_err(|e| Error::InsufficientCapacity(e.required_capacity(), e.remaining()))
}

fn request2buf<B>(r: &Request, buf: &mut B) -> Result<(), Error>
where
    B: BufMut,
{
    msg2buf(r, buf)
}

fn request2bytes(r: &Request, mut bytes: &mut [u8]) -> Result<(), Error> {
    request2buf(r, &mut bytes)
}

fn serialize_request(r: &Request) -> Result<(), Error> {
    #[allow(unsafe_code)]
    let bytes: &mut [u8] = unsafe { &mut _OUT_PB };
    request2bytes(r, bytes)
}

fn micros2request(micros: i64, buf: &mut Request) {
    buf.clear();
    buf.method = "post".into();
    buf.url = "http://localhost:8080".into();
    buf.header.insert(
        "Content-Type".into(),
        HeaderContent {
            values: vec!["application/json".into()],
        },
    );
    buf.header.insert(
        "Custom-Header-Micros".into(),
        HeaderContent {
            values: vec![format!("{micros}")],
        },
    );
    let bytes: &[u8] = br#"{
        "helo": "wrld"
    }"#;
    buf.body = bytes.into();
}

fn _micros2serialized(micros: i64) -> Result<usize, Error> {
    let mut buf: Request = Request {
        method: "".into(),
        url: "".into(),
        header: BTreeMap::new(),
        body: vec![],
    };
    micros2request(micros, &mut buf);
    serialize_request(&buf)?;
    Ok(buf.encoded_len())
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn micros2serialized(micros: i64) -> i32 {
    _micros2serialized(micros)
        .map(|sz: usize| sz as i32)
        .unwrap_or(-1)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn time2req(micros: i64) -> i32 {
    micros2serialized(micros)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn addr() -> *mut u8 {
    unsafe { _OUT_PB.as_mut_ptr() }
}
