
package e

import "fmt"

func Rr(err error, pan bool) bool {
    if err != nil {
        fmt.Print("Error: ")
        fmt.Println(err)
        if pan {
            panic(err)
        }
        return true
    }
    return false
}

