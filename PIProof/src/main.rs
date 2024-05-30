use my_cp_snark::test_alg;

use std::fmt::Debug;
use serde::{Serialize, Deserialize};

// 定义一个泛型结构体，其中包含一个泛型 T 的数据成员
#[derive(Serialize, Deserialize, Debug)]
struct GenericData<T> {
    data: T
}

fn main() {
    let lop = 25;
    let mut ti = 0;
    let mut ti2 = 0;
    for i in 0..lop {
        let (p, v) = test_alg();
        ti = ti + p;
        ti2 = ti2 + v;
    }
    println!("cpsnark prove: {}", ti/lop);
    println!("cpsnark verify: {}", ti2/lop);

    // // 创建一个泛型结构体实例，其中 T 是 i32 类型
    // let int_data = GenericData { data: 123 };
    //
    // // 序列化
    // let serialized = serde_json::to_string(&int_data).unwrap();
    // println!("序列化后的 JSON 字符串: {}", serialized);
    //
    // // 反序列化
    // let deserialized: GenericData<i32> = serde_json::from_str(&serialized).unwrap();
    // println!("反序列化后的结构体: {:?}", deserialized);

}

