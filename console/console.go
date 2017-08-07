
package console

import (
    "fmt"
    "bufio"
    "os"
    "github.com/kpister/spam/e"
)

func Start() {
    fmt.Println("Starting console...")
    reader := bufio.NewReader(os.Stdin)

    for  {
        fmt.Print(">")
        _, or := reader.ReadString('\n')
        e.Rr(or, false)
    }
}
