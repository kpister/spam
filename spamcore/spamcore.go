package spamcore

// Source: http://divan.github.io/posts/go_concurrency_visualize/
// Other source: https://systembash.com/a-simple-go-tcp-server-and-tcp-client/

import (
    "fmt"
    "net"
    "time"
    "bufio"
    "strings"
    "strconv"
    "io/ioutil"

    "github.com/kpister/spam/e"
    "github.com/kpister/spam/peer"
    "github.com/kpister/spam/parsecfg"
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

func send(cfg *parsecfg.Cfg) {
    for {
        for i, v := range cfg.Peers {
            if v.Status == "connected" {
                fmt.Fprintf(v.Conn, time.Now().String() + "\n")
            } else if v.Status == "offline" {
                peer.Connect(&cfg.Peers[i])
            }
        }
        time.Sleep(5000 * time.Millisecond)
    }
}

func handleconsole(logfile string, cfg *parsecfg.Cfg) {
    // Every 1 second read in .log
    // If top line is a command, execute that command
    for {
        time.Sleep(100 * time.Millisecond)
        cmd, or := ioutil.ReadFile(logfile)
        if or != nil {
            continue
        }

        // Information is transfered in the form: c (for console) command params
        // params is further split inside the handling of each command
        pieces := strings.SplitN(strings.TrimSpace(string(cmd)), " ", 3)

        // c: console
        if pieces[0] != "c" {
            continue
        }

        var filecmd string
        // handle commands
        if pieces[1] == "peers" {
            filecmd = ""
            for _, v := range cfg.Peers {
                filecmd += v.Addr + " " + v.Name + " " + v.Status + "?"
            }
        } else if pieces[1] == "broadcast" && len(pieces) == 3 {
            for _, v := range cfg.Peers {
                if v.Status == "connected" {
                    fmt.Fprintf(v.Conn, pieces[2])
                }
            }
            filecmd = "Message sent?"
        } else if pieces[1] == "add" && len(pieces) == 4 {
            ppieces := strings.Split(pieces[2], " ")
            if len(ppieces) != 3 {
                filecmd = "Please supply ip, name, public key"
            } else {
                mpeer := peer.MakePeer(ppieces[0], ppieces[1], ppieces[2])
                cfg.Peers = append(cfg.Peers, *mpeer)
                filecmd = "Peer added: " + ppieces[0] + "?"
            }
        } else if pieces[1] == "dropbyip" && len(pieces) == 3 {
            filecmd = removePeer(true, pieces[2], cfg)
        } else if pieces[1] == "dropbyname" && len(pieces) == 3 {
            filecmd = removePeer(false, pieces[2], cfg)
        }

        ioutil.WriteFile(logfile, []byte(filecmd), 0770)
    }
}

func removePeer(byip bool, text string, cfg *parsecfg.Cfg) string {
    response := "Successfully removed peer: " + text + "?"
    if len(cfg.Peers) == 0 {
        response = "Your peers list is empty. You cannot remove anyone?"
    } else {
        found := false
        for i, v := range cfg.Peers {
            if (byip && v.Addr == text) || (!byip && v.Name == text) {
                if i + 1 < len(cfg.Peers) {
                    cfg.Peers[i] = cfg.Peers[len(cfg.Peers) - 1]
                }
                cfg.Peers = cfg.Peers[0:len(cfg.Peers) -1]
                found = true
            }
        }
        if !found {
            response = "That peer does not exist in your peer list. Please use the peers command to see your peers?"
        }
    }
    return response
}



func StartServer(logfile string, cfg *parsecfg.Cfg) {
    listener, or := net.Listen("tcp", ":" + strconv.Itoa(cfg.Port))
    e.Rr(or, true)
    defer listener.Close()

    ch := make(chan string)
    go logger(ch)
    go server(listener, ch)
    go send(cfg)
    go handleconsole(logfile, cfg)
    for {}
}
