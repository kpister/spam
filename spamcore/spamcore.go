package spamcore

// Source: http://divan.github.io/posts/go_concurrency_visualize/
// Other source: https://systembash.com/a-simple-go-tcp-server-and-tcp-client/

import (
    "fmt"
    "net"
    "bufio"
    "os"
)

func handler(c net.Conn, ch chan string) {
    ch <- c.RemoteAddr().String()
    for {
        message, _ := bufio.NewReader(c).ReadString('\n')
        if message == "" {
            break
        }
        // output message received
        fmt.Print("Message Received:", string(message))
        // send new string back to client
        c.Write([]byte(message + "\n"))

    }
    c.Close()
}

func logger(ch chan string) {
    for {
        fmt.Println(<-ch)
    }
}

func server(l net.Listener, ch chan string) {
    for {
        c, err := l.Accept()
        if err != nil {
            continue
        }
        go handler(c, ch)
    }
}

func beclient(reader *bufio.Reader) {
    var connections []net.Conn
    for {
        fmt.Print("Enter command: ")
        cmd, _ := reader.ReadString('\n')
        cmd = cmd[:len(cmd)-1]
        if cmd == "connect" {
            fmt.Print("Enter address (eg: 127.0.0.1:8080): ")
            conn, _ := reader.ReadString('\n')
            conn =  conn[:len(conn)-1]
            c, err := net.Dial("tcp", conn)
            if err != nil {
                fmt.Println("That peer does not exist")
            } else {
                connections = append(connections, c)
            }
        } else if cmd == "broadcast" {
            if len(connections) == 0 {
                fmt.Println("You have no connections")
                continue
            }
            fmt.Print("Enter message: ")
            text, _ := reader.ReadString('\n')
            for _, v := range connections {
                fmt.Fprintf(v, text + "\n")
            }
        } else if cmd == "exit" {
            return
        } else if cmd == "list" {
            for _, v := range connections {
                fmt.Println(v.RemoteAddr())
            }
        } else if cmd == "remove" {
            if len(connections) == 0 {
                fmt.Println("Your peers list is empty. You cannot remove anyone")
                continue
            } else if len(connections) == 1 {
                fmt.Println("You only have one peer: "+ connections[0].RemoteAddr().String() + ". They have now been removed")
                connections = connections[:0]
            } else {
                fmt.Print("Enter the address of the peer you want to remove: ")
                p, _ := reader.ReadString('\n')
                p = p[0:len(p)-1]
                found := false
                for i, v := range connections {
                    if v.RemoteAddr().String() == p {
                        if i + 1 < len(connections) {
                            connections[i] = connections[len(connections) - 1]
                        }
                        connections = connections[0:len(connections) -1]
                        found = true
                    }
                }
                if !found {
                    fmt.Println("That peer does not exist in your peer list. Please use the list command to see your peers")
                }
            }
        } else {
            fmt.Println("You didn't enter a registered command. Try:\nconnect\nbroadcast\nlist\nexit")
        }
    }
}
func StartServer() {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter port: ")
    port, _ := reader.ReadString('\n')
    l, err := net.Listen("tcp", ":" + port[0: len(port)-1])
    if err != nil {
        panic(err)
    }
    ch := make(chan string)
    go logger(ch)
    go server(l, ch)
    beclient(reader)
}
