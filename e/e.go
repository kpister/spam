
package e

import "fmt"

// Our cute little error package.
// TODO add our own errors

// return True if error. When pan is true - the program will panic on an error and exit
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

