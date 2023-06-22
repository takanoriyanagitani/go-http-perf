(module
  (type (;0;) (func (param i32) (result f32)))
  (type (;1;) (func (param i64) (result f32)))
  (func (;0;) (type 0) (param i32) (result f32)
    local.get 0
    f32.convert_i32_s
    f32.const 0x1p-15 (;=3.05176e-05;)
    f32.mul)
  (func (;1;) (type 0) (param i32) (result f32)
    f32.const 0x1p+0 (;=1;)
    f32.const 0x0p+0 (;=0;)
    local.get 0
    select)
  (func (;2;) (type 0) (param i32) (result f32)
    (local f32)
    local.get 0
    i32.extend16_s
    local.tee 0
    f32.convert_i32_s
    f32.const 0x1p-15 (;=3.05176e-05;)
    f32.mul
    f32.const 0x1.921fb6p+1 (;=3.14159;)
    f32.mul
    local.tee 1
    f32.const 0x1.45f306p+0 (;=1.27324;)
    f32.mul
    local.get 1
    local.get 1
    f32.const -0x1.9f02f4p-2 (;=-0.405285;)
    f32.const 0x1.9f02f4p-2 (;=0.405285;)
    local.get 0
    i32.const 0
    i32.gt_s
    select
    f32.mul
    f32.mul
    f32.add)
  (func (;3;) (type 1) (param i64) (result f32)
    (local i32 f32)
    local.get 0
    i32.wrap_i64
    i32.extend16_s
    local.tee 1
    f32.convert_i32_s
    f32.const 0x1p-15 (;=3.05176e-05;)
    f32.mul
    f32.const 0x1.921fb6p+1 (;=3.14159;)
    f32.mul
    local.tee 2
    f32.const 0x1.45f306p+0 (;=1.27324;)
    f32.mul
    local.get 2
    local.get 2
    f32.const -0x1.9f02f4p-2 (;=-0.405285;)
    f32.const 0x1.9f02f4p-2 (;=0.405285;)
    local.get 1
    i32.const 0
    i32.gt_s
    select
    f32.mul
    f32.mul
    f32.add)
  (memory (;0;) 16)
  (global (;0;) (mut i32) (i32.const 1048576))
  (global (;1;) i32 (i32.const 1048576))
  (global (;2;) i32 (i32.const 1048576))
  (export "memory" (memory 0))
  (export "i4f5" (func 0))
  (export "b2f5" (func 1))
  (export "f32_sin_fast_i32" (func 2))
  (export "f32_sin_fast_u64" (func 3))
  (export "__data_end" (global 1))
  (export "__heap_base" (global 2)))
