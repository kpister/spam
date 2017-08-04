package spamcore

// Source: http://divan.github.io/posts/go_concurrency_visualize/
// Other source: https://systembash.com/a-simple-go-tcp-server-and-tcp-client/

import (
    "fmt"
    "net"
    "bufio"
    "os"

    "github.com/kpister/spam/parsecfg"
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

func beclient(reader *bufio.Reader, cfg *parsecfg.Cfg) {
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
                cfg.Peers = append(cfg.Peers, parsecfg.Peer{c, ""})
            }
        } else if cmd == "broadcast" {
            if len(cfg.Peers) == 0 {
                fmt.Println("You have no connections")
                continue
            }
            fmt.Print("Enter message: ")
            text, _ := reader.ReadString('\n')
            for _, v := range cfg.Peers {
                fmt.Fprintf(v.Conn, text + "\n")
            }
        } else if cmd == "exit" {
            return
        } else if cmd == "list" {
            for _, v := range cfg.Peers {
                fmt.Print(v.Conn.RemoteAddr())
                fmt.Println(v.Name)
            }
        } else if cmd == "remove" {
            if len(cfg.Peers) == 0 {
                fmt.Println("Your peers list is empty. You cannot remove anyone")
                continue
            } else if len(cfg.Peers) == 1 {
                fmt.Println("You only have one peer: "+ cfg.Peers[0].Conn.RemoteAddr().String() + ". They have now been removed")
                cfg.Peers = cfg.Peers[:0]
            } else {
                fmt.Print("Enter the address of the peer you want to remove: ")
                p, _ := reader.ReadString('\n')
                p = p[0:len(p)-1]
                found := false
                for i, v := range cfg.Peers {
                    if v.Conn.RemoteAddr().String() == p {
                        if i + 1 < len(cfg.Peers) {
                            cfg.Peers[i] = cfg.Peers[len(cfg.Peers) - 1]
                        }
                        cfg.Peers = cfg.Peers[0:len(cfg.Peers) -1]
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
func StartServer(cfg *parsecfg.Cfg) {
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
    beclient(reader, cfg)
}
