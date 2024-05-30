package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"os/exec"
	"strings"
	"time"
	"unsafe"

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
	fmt.Printf("Size of pp: %d\n", unsafe.Sizeof(curve))

	// ==========================
	// == Batch key generation ==
	// ==========================
	// == Start ==
	s0 := time.Now()

	const size = 16
	const idx = 0
	const prime_len = "250"

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

	println(len(b.Encode()))
	fmt.Printf("dsk = %d\n", b.Encode())
	fmt.Printf("dpk = %d\n", B.Encode())

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

	// Pseudonym generation
	// == Start ==
	s0 = time.Now()
	p := curve.NewRandomScalar()
	P := curve.ScalarBaseMul(p)

	// Supervisory ciphertext c=(c1,c2)
	// c1=P
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
	A_ := D.Sub(curve.BasePoint()).ScalarMul(za).Add(curve.ScalarBaseMul(zb)).Sub(c2.Sub(P).ScalarMul(curve.ScalarFromBytes(h)))
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

}
