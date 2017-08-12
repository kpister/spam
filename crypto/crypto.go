
package crypto

import (
    "math/big"
)

var e = big.Int

func SetE() {
    e.SetUint64(65537)
}

func Encrypt(m, n big.Int) big.Int {
    var c big.Int
    c = c.Exp(&m, &e, &n)
    return c
}

func Decrypt(c, d, n big.Int) big.Int{
    var m big.Int
    m = m.Exp(&c, &d, &n)
    return m
}

func ConvertMessageToInt(m string) big.Int {
    var total big.Int
    // Find a good way to write this
    return total
}

func ConvertMessageFromInt(m big.Int) string {
    message := ""
    // Find a good way to write this.
    return message
}

