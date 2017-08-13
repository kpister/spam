
package crypto

import (
    "strconv"
    "math/big"
)

var e big.Int

/* provides basic encryption and decryption
 */
func SetE() {
    e.SetUint64(65537)
}

func Encrypt(m, n *big.Int) *big.Int {
    // EncryptedMessage = Message ^ PublicExponent mod (peer's) PublicModulus
    var c big.Int
    ret := c.Exp(m, &e, n)
    return ret
}

func Decrypt(c, d, n *big.Int) *big.Int{
    // DecryptedMessage = EncryptedMessage ^ (your) SecretKey mod (your) PublicModulus
    var m big.Int
    ret := m.Exp(c, d, n)
    return ret
}

func ConvertMessageToInt(m string) *big.Int {
    // TODO: Find a good way to write this
    // We currently turn a string into an ascii representation. Each char is 3 digits.
    var total, expon, it, sol, zero, temp big.Int
    total.SetInt64(0)
    bytes := []byte(m)

    for i, v := range bytes {
        expon.SetInt64(int64((len(bytes) - i - 1)))
        it.SetInt64(1000)
        temp.SetInt64(int64(v))
        zero.SetInt64(0)
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

