
package console

import (
    "os"
    "fmt"
    "time"
    "bufio"
    "strings"

    "github.com/kpister/spam/e"
)

func Start(logfile string) {
    fmt.Println("Starting console...")
    reader := bufio.NewReader(os.Stdin)

    log, _ := os.Open(logfile)
    readwrite := bufio.NewReadWriter(bufio.NewReader(log), bufio.NewWriter(log))

    for  {
        fmt.Print(">")
        cmd, or := reader.ReadString('\n')
        e.Rr(or, false)

        log.Seek(0,0)

        if cmd == "peers\n" {
            readwrite.WriteString("c peers\n")
        }
        readwrite.Flush()

        time.Sleep(1200 * time.Millisecond)

        log.Seek(0,0)

        line, or := readwrite.ReadString('\n')
        output := strings.Replace(line, "?", "\n", -1)
        fmt.Println(output)


    }
}
