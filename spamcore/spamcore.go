package spamcore

// Source: http://divan.github.io/posts/go_concurrency_visualize/
// Other source: https://systembash.com/a-simple-go-tcp-server-and-tcp-client/

import (
    "fmt"
    "net"
    "bufio"
    "time"
    "strconv"

    "github.com/kpister/spam/e"
    "github.com/kpister/spam/parsecfg"
    "github.com/kpister/spam/peer"
)

func handler(conn net.Conn, ch chan string) {
    reader := bufio.NewReader(conn)
    remoteAddr := conn.RemoteAddr().String()

    for {
        message, or := reader.ReadString('\n')
        if or != nil && or.Error() == "EOF" {
            ch <- "Disconnected from " + remoteAddr + "\n"
            break
        } else if !e.Rr(or, true) {
            // output message received
            ch <- "Message Received from " + remoteAddr +":" + string(message)
        }
    }
    conn.Close()
}

func logger(ch chan string) {
    for {
        fmt.Print(<-ch)
    }
}

func server(listener net.Listener, ch chan string) {
    for {
        conn, or := listener.Accept()
        defer conn.Close()
        if e.Rr(or, false) {
            continue
        }
        go handler(conn, ch)
    }
}
/*
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
*/

func send(cfg *parsecfg.Cfg) {
    for {
        if len(cfg.Peers) > 0 {
            for i, v := range cfg.Peers {
                if v.Status == "connected" {
                    fmt.Fprintf(v.Conn, time.Now().String() + "\n")
                } else if v.Status == "offline" {
                    peer.Connect(&cfg.Peers[i])
                }
            }
        }
        time.Sleep(5000 * time.Millisecond)
    }
}
func StartServer(cfg *parsecfg.Cfg) {
    listener, or := net.Listen("tcp", ":" + strconv.Itoa(cfg.Port))
    e.Rr(or, true)
    defer listener.Close()

    ch := make(chan string)
    go logger(ch)
    go server(listener, ch)
    go send(cfg)
    for {}
}
