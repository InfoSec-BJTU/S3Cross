# S$^3$Cross
Rust and go implementation of the S3Cross cross-domain authentication scheme

## Paper
S$^3$Cross: Blockchain-Based Cross-Domain Authentication with Self-Sovereign and Supervised Identity Management

## Usage

All the codes are tested on both Mac and Ubuntu platforms.

Mac Studio (M1 Max, 64GB RAM, macOS 14.4.1):

|Language|Version|Architecture|
|:------:|-------|------------|
|Rust|rustc 1.73.0-nightly (500647fd8 2023-07-27)|aarch64-apple-darwin|
|Go|go1.20.1|darwin/arm64|

Ubuntu Server (Intel(R) Xeon(R) Gold 6230, 125GB RAM, Ubuntu 20.04.6 LTS):

|Language|Version|Architecture|
|:------:|-------|------------|
|Rust|cargo 1.72.0-nightly (5b377cece 2023-06-30)|x86_64-unknown-linux-gnu|
|Go|go1.20.4|linux/amd64|

To run the pseudonym management and cross-domain authentication request verification/generation process, run `go run psd.go` in the `Psd` directory.

To run the PIProof process, run `cargo run` in the `PIProof` directory.

Programs written using Rust can achieve better performance in the following mode, where "xxxx" is the name of the project (can be found in Cargo.toml):

- cargo build --release
- target/release/xxxx

## Warning

**This library is not ready for production use!**

Theoretically, the specific values of `Com_sr` (line 274) and `V` (line 346) in `Psd/psd.go` should be the same to prove the knowledge and the range of `jd` simultaneously. But the curves used by the other parts and [\[bulletproofs\]](https://github.com/neucc1997/bulletproofs) of `pi2` in `Psd` are different. We will address the issue of inconsistent curves in later versions.


