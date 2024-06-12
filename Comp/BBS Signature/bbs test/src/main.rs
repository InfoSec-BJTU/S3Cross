extern crate mylib;

use mylib::*;

use std::time::Instant;

use std::mem;

fn main() {
    use rand::thread_rng;
    let mut rng = thread_rng();
    // 初始化群公私钥
    let SetUpResult { gsk, gpk } = setup(&mut rng);
    // 生成设备私钥
    let isk = issue(&gsk, &gpk, &mut rng);

    let start = Instant::now();
    let sig = sign(&isk, &gpk, &mut rng);
    let end = Instant::now();
    // 计算程序执行时间     
    let duration = end.duration_since(start);      
    println!("BBS sign time: {:?}", duration);

    let start = Instant::now();
    verify(&sig, &gpk).unwrap();
    let end = Instant::now();     
    // 计算程序执行时间          
    let duration = end.duration_since(start);
    println!("BBS verify time: {:?}", duration);

    assert!(is_signed_member(&isk, &sig, &gsk));
}
