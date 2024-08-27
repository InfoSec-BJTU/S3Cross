use my_cp_snark::test_alg;

fn main() {
    let lop = 25;
    let mut ti = 0;
    let mut ti2 = 0;
    for i in 0..lop {
        let (p, v) = test_alg();
        ti = ti + p;
        ti2 = ti2 + v;
    }
    println!("Time for proving 'PIProof' : {}", ti/lop);
    println!("Time for verifying 'PIProof' : {}", ti2/lop);

}

