use rand::thread_rng;

use curve25519_dalek::scalar::Scalar;

use merlin::Transcript;

use bulletproofs::{BulletproofGens, PedersenGens, RangeProof};

use std::time::Instant;

pub fn test_bp() {
    let mut now = Instant::now();
    // Generators for Pedersen commitments.  These can be selected
    // independently of the Bulletproofs generators.
    let pc_gens = PedersenGens::default();

    // Generators for Bulletproofs, valid for proofs up to bitsize 64
    // and aggregation size up to 1.
    let bp_gens = BulletproofGens::new(8, 1);

    // A secret value we want to prove lies in the range [0, 2^32)
    let secret_value = 1;

    // The API takes a blinding factor for the commitment.
    let blinding = Scalar::random(&mut thread_rng());

    // The proof can be chained to an existing transcript.
    // Here we create a transcript with a doctest domain separator.
    let mut prover_transcript = Transcript::new(b"doctest example");

    // Create a 32-bit rangeproof.
    let (proof, committed_value) = RangeProof::prove_single(
        &bp_gens,
        &pc_gens,
        &mut prover_transcript,
        secret_value,
        &blinding,
        8,
    ).expect("A real program could handle errors");
    println!("Time for prove : {}", now.elapsed().as_micros());

    now = Instant::now();
    // Verification requires a transcript with identical initial state:
    let mut verifier_transcript = Transcript::new(b"doctest example");
    assert!(
        proof
            .verify_single(&bp_gens, &pc_gens, &mut verifier_transcript, &committed_value, 8)
            .is_ok()
    );
    println!("Time for verify: {}", now.elapsed().as_micros());
}


