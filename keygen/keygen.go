
package keygen

import (
    "fmt"
    "crypto/rsa"
    "crypto/rand"

    "github.com/kpister/spam/e"
)

// Stand alone code to generate public and private keys.
// Keep both for yourself and share the public with your friends!
// This is run with spam -gen-keypair
func GenKeys() {
    reader := rand.Reader
    bitsize := 2048

    key, or := rsa.GenerateKey(reader, bitsize)
    e.Rr(or, true)

    fmt.Println("===========Public Key==========")
    fmt.Println(key.PublicKey.N)
    fmt.Println("===============================")
    fmt.Println("===========Private Key=========")
    fmt.Println(key.D)
    fmt.Println("===============================")
}
