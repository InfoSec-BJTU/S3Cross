package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"os/exec"
	"strings"
	"time"

	"github.com/athanorlabs/go-dleq/types"
	ring "github.com/noot/ring-go"
	"golang.org/x/crypto/sha3"
)

func schnorrSig(pri types.Scalar, message []byte, curve ring.Curve) (types.Scalar, [32]byte) {

	s := time.Now()

	pub := curve.ScalarBaseMul(pri)
	msgHash := sha3.Sum256(message)
	rand := curve.NewRandomScalar()
	Rand := curve.ScalarBaseMul(rand)

	// c=H(X|R|m)
	var c_ []byte
	c_ = append(pub.Encode(), Rand.Encode()...)
	sli := msgHash[:]
	c_ = append(c_, sli...)
	c := sha3.Sum256(c_)

	// s=r-c*d，输出（c,s）
	sig := rand.Sub(curve.ScalarFromBytes(c).Mul(pri))

	e := time.Now()
	d := e.Sub(s)
	fmt.Printf("Schnorr signature generation time = %v\n", d)

	return sig, c

}

func schnorrSigVer(pub types.Point, message []byte, curve ring.Curve, sig types.Scalar, c [32]byte) bool {

	s := time.Now()

	// Q=sG+cP, calculate c' and compare
	Q := curve.ScalarBaseMul(sig).Add(pub.ScalarMul(curve.ScalarFromBytes(c)))
	msgHash := sha3.Sum256(message)
	var c__ []byte
	c__ = append(pub.Encode(), Q.Encode()...)
	c__ = append(c__, msgHash[:]...)
	cc := sha3.Sum256(c__)

	isok := bytes.Equal(c[:], cc[:])

	e := time.Now()
	d := e.Sub(s)
	fmt.Printf("Schnorr signature verification time = %v\n", d)

	if isok {
		return true
	}
	return false
}

// oepnssl prime -generate -bits 250 -hex
func getPrime(bits string) string {
	app := "openssl"
	arg0 := "prime"
	arg1 := "-generate"
	arg2 := "-bits"
	arg3 := bits
	arg4 := "-hex"
	cmd := exec.Command(app, arg0, arg1, arg2, arg3, arg4)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return "0"
	}

	return string(stdout)
}

func hexToBytes(hexString string) []byte {
	hexString = strings.TrimSpace(hexString)

	// add '\n'
	decodedByteArray, _ := hex.DecodeString(hexString)
	return decodedByteArray
}

func hexToBigInt(hex string) *big.Int {
	n := new(big.Int)
	n, _ = n.SetString(hex, 16)

	return n
}

func main() {
	// Initialization
	curve := ring.Secp256k1()
	// fmt.Printf("Size of pp: %d\n", unsafe.Sizeof(curve))

	// ==========================
	// == Batch key generation ==
	// ==========================
	// == Start ==
	s0 := time.Now()

	const size = 16
	const idx = 0
	const prime_len = "250"
	const prime_len_2 = "1024"

	pris := make([]types.Scalar, size)
	pubs := make([]types.Point, size)
	for i := 0; i < size; i++ {
		idx := (i + idx) % size
		priv, _ := curve.DecodeToScalar(hexToBytes(getPrime(prime_len)))
		// priv := curve.NewRandomScalar()
		// fmt.Printf("Size of priv: %d bytes\n", int64(reflect.TypeOf(priv).Size()))
		pris[idx] = priv
		pubs[idx] = curve.ScalarBaseMul(priv)
	}

	e0 := time.Now()
	// == End ==
	d0 := e0.Sub(s0)
	fmt.Printf("Key batch generation time = %v\n", d0)

	// Specify a key pair
	b := pris[7]
	B := pubs[7]

	// println(len(b.Encode()))
	// fmt.Printf("dsk = %d\n", b.Encode())
	// fmt.Printf("dpk = %d\n", B.Encode())

	// // Encoding Test
	// ind := curve.NewRandomScalar()
	// fmt.Println(ind.Encode())                     // o
	// fmt.Println(hex.EncodeToString(ind.Encode())) // o
	// i := new(big.Int)
	// i.SetString(hex.EncodeToString(ind.Encode()), 16)
	// fmt.Println(i)

	// Supervisor key generation
	s0 = time.Now()
	d := curve.NewRandomScalar()
	D := curve.ScalarBaseMul(d)
	e0 = time.Now()
	d0 = e0.Sub(s0)
	fmt.Printf("Supervisor key generation time: %v\n", d0)

	// ==========================
	// == Pseudonym generation ==
	// ==========================
	// == Start ==
	s0 = time.Now()
	p := curve.NewRandomScalar()
	P := curve.ScalarBaseMul(p)

	// Supervisory ciphertext c=(c1,c2)
	c1 := P
	c2 := D.ScalarMul(p).Add(B)

	// GenPsdProof
	alpha := curve.NewRandomScalar()
	beta := curve.NewRandomScalar()
	A := D.Sub(curve.BasePoint()).ScalarMul(alpha).Add(curve.ScalarBaseMul(beta))

	var comb []byte
	comb = append(curve.BasePoint().Encode(), P.Encode()...)
	comb = append(comb, D.Encode()...)
	comb = append(comb, P.Encode()...)
	comb = append(comb, c2.Encode()...)
	comb = append(comb, A.Encode()...)
	h := sha3.Sum256(comb)

	za := alpha.Add(curve.ScalarFromBytes(h).Mul(p))
	zb := beta.Add(curve.ScalarFromBytes(h).Mul(b))

	e0 = time.Now()
	// == End ==

	d0 = e0.Sub(s0)
	fmt.Printf("Pseudonym generation time: %v\n", d0)

	// VerifyPsdProof
	// == Start
	s0 = time.Now()
	A_ := D.Sub(curve.BasePoint()).ScalarMul(za).Add(curve.ScalarBaseMul(zb)).Sub(c2.Sub(c1).ScalarMul(curve.ScalarFromBytes(h)))
	var comb2 []byte
	comb2 = append(curve.BasePoint().Encode(), P.Encode()...)
	comb2 = append(comb2, D.Encode()...)
	comb2 = append(comb2, P.Encode()...)
	comb2 = append(comb2, c2.Encode()...)
	comb2 = append(comb2, A_.Encode()...)
	h_ := sha3.Sum256(comb2)
	res := curve.ScalarFromBytes(h).Eq(curve.ScalarFromBytes(h_))
	e0 = time.Now()
	// == End

	if !res {
		fmt.Println("fail")
	} else {
		fmt.Println("VerpsdProof Success!")
	}

	d0 = e0.Sub(s0)
	fmt.Printf("Pseudonym verification time: %v\n", d0)

	// Decryption
	// == Start
	s0 = time.Now()
	pkd := c2.Sub(P.ScalarMul(d))
	res = B.Equals(pkd)
	e0 = time.Now()
	// == End

	d0 = e0.Sub(s0)
	fmt.Printf("Decrypt and verify time: %v\n", d0)
	println(res)

	s0 = time.Now()
	Np := getPrime(prime_len)
	e0 = time.Now()

	d0 = e0.Sub(s0)
	fmt.Printf("prime 250 generation time: %v\n", d0)
	fmt.Printf("%s\n", Np)

	for i := 0; i < 250; i++ {
		s0 = time.Now()
		Np = getPrime(prime_len)
		e0 = time.Now()
		d0 += e0.Sub(s0)
	}
	fmt.Printf("prime 250 generation time (avg): %v\n", d0/251)


	s0 = time.Now()
	Np = getPrime(prime_len_2)
	e0 = time.Now()

	d0 = e0.Sub(s0)
	fmt.Printf("prime 1024 generation time: %v\n", d0)
	fmt.Printf("%s\n", Np)

	for i := 0; i < 50; i++ {
		s0 = time.Now()
		Np = getPrime(prime_len_2)
		e0 = time.Now()
		d0 += e0.Sub(s0)
	}
	fmt.Printf("prime 1024 generation time (avg): %v\n", d0/51)

	// =========== S3Cross-SR: pi2 ===========
	// rand := curve.NewRandomScalar()
	// Rand := curve.ScalarBaseMul(rand)

	// Random generator
	h_sr := curve.ScalarBaseMul(curve.NewRandomScalar())

	a := curve.NewRandomScalar()
	jd := curve.ScalarFromInt(12)
	one := curve.ScalarFromInt(1)
	p_sr := a.Mul((b.Add(jd).Add(one)).Inverse())
	P_sr := curve.ScalarBaseMul(p_sr)

	//  Enc: Change the pseudonym private key and public key
	c1_sr := curve.ScalarBaseMul(p_sr)
	c2_sr := D.ScalarMul(p_sr).Add(B)

	s0 = time.Now()
	alpha_sr := curve.NewRandomScalar()
	beta_sr := curve.NewRandomScalar()
	gamma_sr := curve.NewRandomScalar()
	theta_sr := curve.NewRandomScalar()

	// A=P^{alpha+beta}
	A_sr := P_sr.ScalarMul(alpha_sr.Add(beta_sr))
	// B=(P/g)^{gamma}·g^{beta}
	B_sr := D.Sub(curve.BasePoint()).ScalarMul(gamma_sr).Add(curve.BasePoint().ScalarMul(beta_sr))
	// C=g^{alpha}·h^{theta_sr}
	k_sr := curve.NewRandomScalar()
	// Com=g^{jd}·h^{k}
	Com_sr := curve.ScalarBaseMul(jd).Add(h_sr.ScalarMul(k_sr))
	C_sr := curve.ScalarBaseMul(alpha_sr).Add(h_sr.ScalarMul(theta_sr))

	var comb_sr []byte
	comb_sr = append(A_sr.Encode(), B_sr.Encode()...)
	comb_sr = append(comb_sr, C_sr.Encode()...)
	comb_sr = append(comb_sr, P_sr.Encode()...)
	comb_sr = append(comb_sr, D.Encode()...)
	comb_sr = append(comb_sr, curve.BasePoint().Encode()...)
	comb_sr = append(comb_sr, h_sr.Encode()...)
	comb_sr = append(comb_sr, Com_sr.Encode()...)
	hash_sr := sha3.Sum256(comb_sr)

	z_alpha_sr := alpha_sr.Add(jd.Mul(curve.ScalarFromBytes(hash_sr)))
	z_beta_sr := beta_sr.Add(b.Mul(curve.ScalarFromBytes(hash_sr)))
	z_gamma_sr := gamma_sr.Add(p_sr.Mul(curve.ScalarFromBytes(hash_sr)))
	z_theta_sr := theta_sr.Add(k_sr.Mul(curve.ScalarFromBytes(hash_sr)))

	e0 = time.Now()
	d0 = e0.Sub(s0)
	fmt.Printf("pi2 generation time: %v\n", d0)

	// pi2=(hash_sr,z_alpha_sr,z_beta_sr,z_gamma_sr,z_theta_sr)

	// Verify
	s0 = time.Now()
	_A_sr := P_sr.ScalarMul(z_alpha_sr.Add(z_beta_sr)).Sub((curve.BasePoint().ScalarMul(a).Sub(P_sr)).ScalarMul(curve.ScalarFromBytes(hash_sr)))
	_B_sr := D.Sub(curve.BasePoint()).ScalarMul(z_gamma_sr).Add(curve.BasePoint().ScalarMul(z_beta_sr)).Sub((c2_sr.Sub(c1_sr)).ScalarMul(curve.ScalarFromBytes(hash_sr)))
	_C_sr := curve.ScalarBaseMul(z_alpha_sr).Add(h_sr.ScalarMul(z_theta_sr)).Sub(Com_sr.ScalarMul(curve.ScalarFromBytes(hash_sr)))

	var _comb_sr []byte
	_comb_sr = append(_A_sr.Encode(), _B_sr.Encode()...)
	_comb_sr = append(_comb_sr, _C_sr.Encode()...)
	_comb_sr = append(_comb_sr, P_sr.Encode()...)
	_comb_sr = append(_comb_sr, D.Encode()...)
	_comb_sr = append(_comb_sr, curve.BasePoint().Encode()...)
	_comb_sr = append(_comb_sr, h_sr.Encode()...)
	_comb_sr = append(_comb_sr, Com_sr.Encode()...)
	_hash_sr := sha3.Sum256(_comb_sr)
	res = curve.ScalarFromBytes(hash_sr).Eq(curve.ScalarFromBytes(_hash_sr))
	if res {
		fmt.Println("pi2 verification passed")
	}
	e0 = time.Now()
	d0 = e0.Sub(s0)
	fmt.Printf("pi2 verification time: %v\n", d0)

	// // p_sr Test
	// ta := curve.ScalarFromInt(24)
	// tb := curve.ScalarFromInt(2)
	// tjd := curve.ScalarFromInt(5)
	// tone := curve.ScalarFromInt(1)
	// t_p_sr := ta.Mul((tb.Add(tjd).Add(tone)).Inverse())

	// // p_sr := a.Mul((b.Add(jd).Add(one)).Inverse())

	// _t_p_sr := curve.ScalarFromInt(3)

	// println(t_p_sr)
	// println(_t_p_sr)
	// println(t_p_sr.Eq(_t_p_sr))

	// s0 = time.Now()
	// // the scalar we want to generate a range proof for
	// v := big.NewInt(12)
	// //
	// // gamma := big.NewInt(10)
	// gamma := new(big.Int)
	// gamma.SetBytes(k_sr.Encode())
	
	// prover := bp.NewProver(4)

	// // V = γH + vG.
	// V := bp.Commit(gamma, prover.BlindingGenerator, v, prover.ValueGenerator)

	// proof, err := prover.CreateRangeProof(V, v, gamma, [32]byte{}, [16]byte{})
	// if err != nil {
	// 	fmt.Println("failed to create range proof: ", err)
	// }
	// e0 = time.Now()
	// d0 = e0.Sub(s0)
	// fmt.Printf("pi2_range generation time: %v\n", d0)

	// s0 = time.Now()
	// if !prover.Verify(V, proof) {
	// 	fmt.Println("Expected valid proof")
	// } else {
	// 	fmt.Println("Valid bp proof")
	// }
	// e0 = time.Now()
	// d0 = e0.Sub(s0)
	// fmt.Printf("pi2_range verification time: %v\n", d0)





	//  Basic benchmark

	// s0 = time.Now()
	// T_mul := curve.ScalarBaseMul(pris[3]);
	// e0 = time.Now()

	// d0 = e0.Sub(s0)
	// fmt.Printf("T_mul time: %v\n", d0)
	// fmt.Printf("%s\n", T_mul.Encode())

	// for i := 0; i < 250; i++ {
	// 	s0 = time.Now()
	// 	T_mul = curve.ScalarBaseMul(pris[3]);
	// 	e0 = time.Now()
	// 	d0 += e0.Sub(s0)
	// }
	// fmt.Printf("T_mul time (avg): %v\n", d0/251)


	// T_add_1 := curve.ScalarBaseMul(pris[3]);
	// T_add_2 := curve.ScalarBaseMul(pris[4]);
	// s0 = time.Now()
	// T_add := T_add_1.Add(T_add_2);
	// e0 = time.Now()
	
	// d0 = e0.Sub(s0)
	// fmt.Printf("T_add time: %v\n", d0)
	// fmt.Printf("%s\n", T_add.Encode())

	// for i := 0; i < 250; i++ {
	// 	s0 = time.Now()
	// 	T_add = T_add_1.Add(T_add_2);
	// 	e0 = time.Now()
	// 	d0 += e0.Sub(s0)
	// }
	// fmt.Printf("T_add time (avg): %v\n", d0/251)
}
