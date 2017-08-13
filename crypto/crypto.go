
package crypto

import (
    "math/big"
)

var e big.Int

func SetE() {
    e.SetUint64(65537)
}

func Encrypt(m, n *big.Int) *big.Int {
    var c big.Int
    c.Exp(m, &e, n)
    return &c
}

func Decrypt(c, d, n *big.Int) *big.Int{
    var m big.Int
    m.Exp(c, d, n)
    return &m
}

func ConvertMessageToInt(m string) *big.Int {
    var total big.Int
    total.SetInt64(0)
    var length big.Int
    length.SetInt64(int64(len(m)))
    // Find a good way to write this
    total.Add(&total, &length)
    return &total
}

func ConvertMessageFromInt(m big.Int) string {
    message := string(m.Uint64())
    // Find a good way to write this.
    return message
}

