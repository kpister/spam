package spamcore

// Source: http://divan.github.io/posts/go_concurrency_visualize/
// Other source: https://systembash.com/a-simple-go-tcp-server-and-tcp-client/

import (
    "fmt"
    "net"
    "bufio"
    "os"
    "strconv"

    "github.com/kpister/spam/e"
    "github.com/kpister/spam/parsecfg"
    "github.com/kpister/spam/peer"
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
            mpeer := peer.MakePeer(conn, "")
            if mpeer != nil {
                cfg.Peers = append(cfg.Peers, *mpeer)
            } else {
                fmt.Println("Could not connect to that peer")
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
    var listener net.Listener
    var reader *bufio.Reader
    for cfg.Port == 0 {
        reader = bufio.NewReader(os.Stdin)
        fmt.Print("Enter port: ")
        port, or := reader.ReadString('\n')
        e.Rr(or, false)
        cfg.Port, or = strconv.Atoi(port)
        e.Rr(or, false)

        listener, or = net.Listen("tcp", ":" + string(cfg.Port))
        if e.Rr(or, false) {
            cfg.Port = 0
        }
    }

    ch := make(chan string)
    go logger(ch)
    go server(listener, ch)
    beclient(reader, cfg)
}
