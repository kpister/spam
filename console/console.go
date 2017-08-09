
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

    log, or := os.OpenFile(logfile, os.O_RDWR, 0777)
    e.Rr(or, true)
    readwrite := bufio.NewReadWriter(bufio.NewReader(log), bufio.NewWriter(log))

    for  {
        fmt.Print(">")
        cmd, or := reader.ReadString('\n')
        e.Rr(or, false)

        log.Seek(0,0)

        if cmd == "peers\n" {
            _, or := readwrite.WriteString("c peers\n")
            e.Rr(or, false)
        }
        or = readwrite.Flush()
        e.Rr(or, false)
        time.Sleep(1200 * time.Millisecond)

        log.Seek(0,0)

        line, or := readwrite.ReadString('\n')
        output := strings.Replace(line, "?", "\n", -1)
        fmt.Print(output)

    }
}
