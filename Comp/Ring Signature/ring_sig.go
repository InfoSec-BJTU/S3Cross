package main

import (
	"fmt"
	"time"

	ring "github.com/noot/ring-go"
	"golang.org/x/crypto/sha3"
)

func signAndVerify(curve ring.Curve) {
	privkey := curve.NewRandomScalar()
	msgHash := sha3.Sum256([]byte("helloworld"))

	// size of the public key ring (anonymity set)
	const size = 160

	// our key's secret index within the set
	const idx = 50

	s := time.Now()

	keyring, err := ring.NewKeyRing(curve, size, privkey, idx)
	if err != nil {
		panic(err)
	}

	sig, err := keyring.Sign(msgHash, privkey)
	if err != nil {
		panic(err)
	}
	e := time.Now()
	d := e.Sub(s)
	fmt.Printf("Ring signature generation time = %v\n", d)

	ss := time.Now()
	ok := sig.Verify(msgHash)
	if !ok {
		fmt.Println("failed to verify :(")
		return
	}
	ee := time.Now()
	dd := ee.Sub(ss)
	fmt.Printf("Ring signature verification time = %v\n", dd)
}

func main() {
	fmt.Println("using secp256k1...")
	signAndVerify(ring.Secp256k1())
}
