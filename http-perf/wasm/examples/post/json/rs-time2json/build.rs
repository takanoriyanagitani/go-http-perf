use std::io;

fn main() -> Result<(), io::Error> {
    let mut cfg = prost_build::Config::new();
    cfg.btree_map(["."]);
    cfg.compile_protos(&["src/request.proto"], &["src/"])?;
    Ok(())
}
