
package crypto

import (
    "strconv"
    "math/big"
)

var e big.Int

func SetE() {
    e.SetUint64(65537)
}

func Encrypt(m, n *big.Int) *big.Int {
    var c big.Int
    ret := c.Exp(m, &e, n)
    return ret
}

func Decrypt(c, d, n *big.Int) *big.Int{
    var m big.Int
    ret := m.Exp(c, d, n)
    return ret
}

func ConvertMessageToInt(m string) *big.Int {
    // TODO: Find a good way to write this
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

