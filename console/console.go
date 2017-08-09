
package console

import (
    "os"
    "fmt"
    "time"
    "bufio"
    "strings"
    "io/ioutil"

    "github.com/kpister/spam/e"
)

func Start(logfile string) {
    fmt.Println("Starting console...")
    reader := bufio.NewReader(os.Stdin)

    for  {
        fmt.Print(">")
        cmd, or := reader.ReadString('\n')
        e.Rr(or, false)

        cmd = strings.TrimSpace(cmd)

        var filecmd string

        if cmd == "peers" {
            filecmd = "c peers\n"
        } else if cmd == "broadcast" {
            fmt.Print("Enter text to send: ")
            text, or := reader.ReadString('\n')
            if e.Rr(or, false) {
                continue
            }

            filecmd = "c broadcast " + text
        } else if cmd == "add peer" {
            fmt.Print("Enter ip:port of peer you wish to add: ")
            text, or := reader.ReadString('\n')
            if e.Rr(or, false) {
                continue
            }

            filecmd = "c add " + text
        } else if cmd == "drop peer by ip" {
            fmt.Print("Enter ip:port of peer you wish to drop: ")
            text, or := reader.ReadString('\n')
            if e.Rr(or, false) {
                continue
            }

            filecmd = "c dropbyip " + text
        } else if cmd == "drop peer by name" {
            fmt.Print("Enter the name of peer you wish to drop (you must drop unnamed peers by ip): ")
            text, or := reader.ReadString('\n')
            if e.Rr(or, false) {
                continue
            }

            filecmd = "c dropbyname " + text
        } else if cmd == "drop peer" {
            fmt.Println("To drop a peer, either use `drop peer by ip` or `drop peer by name`")
            continue
        } else {
            fmt.Println("That command didn't make sense. Try again")
            continue
        }

        or = ioutil.WriteFile(logfile, []byte(filecmd), 0770)
        e.Rr(or, false)

        time.Sleep(200 * time.Millisecond)

        line, or := ioutil.ReadFile(logfile)
        output := strings.Replace(string(line), "?", "\n", -1)
        fmt.Print(output)

    }
}
