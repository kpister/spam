
package crypto

import (
    "strconv"
    "math/big"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/rand"
)

var e big.Int

/* provides basic encryption and decryption
 */

// This is the public key exponent. It is shared among all keys. If we use a different bit value for keygen, this must change
func SetE() {
    e.SetUint64(65537)
}

// TODO: concatenate message + hash
func Encrypt(m, n *big.Int) *big.Int {
    // EncryptedMessage = Message ^ PublicExponent mod (peer's) PublicModulus
    var c big.Int
    c.Exp(m, &e, n)
    return &c
}

// TODO: split message from hash and return both
func Decrypt(c, d, n *big.Int) *big.Int{
    // DecryptedMessage = EncryptedMessage ^ (your) SecretKey mod (your) PublicModulus
    var m big.Int
    m.Exp(c, d, n)
    return &m
}

// sign then encrypt
func Sign(privkey *big.Int, m string) ([]byte, error) {
	// Get message hash
	digest := hash(m)
	// Create rsa.PrivateKey object 
	priv := new(rsa.PrivateKey)  
	priv.N = privkey

	return rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, digest[:])
}

func Verify(pubKey *big.Int, m string , s []byte) bool{
	pub := new(rsa.PublicKey)
	pub.N = pubKey
	digest := hash(m)	
	err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, digest, s)
	if err == nil {
		return true
	} 
	return false
}

func hash(m string) []byte {
	h := sha256.New()
	h.Write([]byte(m))
	digest := h.Sum(nil)
	return digest
}

func ConvertMessageToInt(m string) *big.Int {
    // TODO: Find a good way to write this
    // We currently turn a string into an ascii representation. Each char is 3 digits.
    var total, expon, it, sol, zero, temp big.Int
    total.SetInt64(0)
    bytes := []byte(m)
    it.SetInt64(1000)
    zero.SetInt64(0)

    for i, v := range bytes {
        expon.SetInt64(int64((len(bytes) - i - 1)))
        temp.SetInt64(int64(v))
        sol.Exp(&it, &expon, &zero)
        sol.Mul(&sol, &temp)
        total.Add(&total, &sol)
    }
    return &total
}

func ConvertMessageFromInt(m *big.Int) string {
    // TODO: Find a good way to write this.
    // Reverse ConvertMessageToInt
    intmessage := m.String()
    message := ""

    for {
        if len(intmessage) % 3 != 0 {
            intmessage = "0" + intmessage
        } else {
            break
        }
    }
    for i := 0; i < len(intmessage) - 2; i+=3 {
        piece, _ := strconv.Atoi(intmessage[i:i+3])
        message += string(piece)
    }
    return message
}

