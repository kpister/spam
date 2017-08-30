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
    "math/big"
    "io/ioutil"

    "github.com/kpister/spam/e"
    "github.com/kpister/spam/peer"
    "github.com/kpister/spam/crypto"
    "github.com/kpister/spam/parsecfg"
)

// Here we handle incoming messages
func handler(conn net.Conn, cfg *parsecfg.Cfg, reader *bufio.Reader) {
    remoteAddr := conn.RemoteAddr().String()
    var p peer.Peer

    for {
        message, or := reader.ReadString('\n')
        if p.RemoteAddr == "" {
            for i, v := range cfg.Peers {
                if v.RemoteAddr == remoteAddr {
                    p = cfg.Peers[i]
                }
            }
        }
        if or != nil && or.Error() == "EOF" {
            fmt.Println("Disconnected from " + remoteAddr)
            for i, v := range cfg.Peers {
                if v.RemoteAddr == remoteAddr {
                    cfg.Peers[i].Status = "offline"
                    cfg.Peers[i].RemoteAddr = ""
                }
            }
            break // This is when they close their node
        } else if !e.Rr(or, false) {
            // output message received
            // We only acknowledge messages from peers we have authenticated (authrec at least)
            if p.RemoteAddr != "" {
                fmt.Print("Message Received from " + p.Name +":" + string(message))
            }
        } else {
            break
        }
    }
    conn.Close()
}

// Decrypt message and figure out who sent it
// TODO add digital signiture 
func handleshake(keystring, remoteaddr string, cfg *parsecfg.Cfg){
    var k big.Int
    key, suc := k.SetString(strings.TrimSpace(keystring), 10)
    fmt.Println("Trying to handshake with ", remoteaddr, "...")
    if suc {
        decrypted := crypto.Decrypt(key, &(cfg.SecretKey), &(cfg.PublicKey))
        payload := crypto.ConvertMessageFromInt(decrypted)
		// TODO Obtain signature from payload
		message := payload[0:len(payload)-256]
		signature := []byte(payload[len(payload)-256:])

		if ! crypto.Verify(key, message, signature) {
			// TODO: throw exception or return failure
			fmt.Println("failed to verify signature.")
		}
		for i, v := range cfg.Peers {
            if message == v.Addr {
                fmt.Println("Decrypted message: " + message + " Peer addr: " + v.Addr)
                cfg.Peers[i].RemoteAddr = remoteaddr
                if v.Status == "authsent" {
                    cfg.Peers[i].Status = "authenticated"
               } else {
                    cfg.Peers[i].Status = "authrec"
                }
            }
        }
    }
}

// Start the actual server which spawns off all the little threads
func server(listener net.Listener, cfg *parsecfg.Cfg) {
    for {
        conn, or := listener.Accept()
        defer conn.Close()
        if e.Rr(or, false) {
            continue
        }
        reader := bufio.NewReader(conn)
        message, or := reader.ReadString('\n')
        if e.Rr(or, false){
            continue
        }
        go handleshake(strings.Split(string(message), ":")[1], conn.RemoteAddr().String(), cfg) // Handle handshake
        go handler(conn, cfg, reader)
    }
}

// Here is where we actually broadcast messages to peers we have authsent at least
// Alternatively we try and connect to our offline and authrec peers
func send(cfg *parsecfg.Cfg) {
    for {
        for i, v := range cfg.Peers {
            if v.Status == "authenticated" || v.Status == "authsent" {
                fmt.Fprintf(v.Conn, time.Now().String() + "\n")
            } else if v.Status == "offline" || v.Status == "authrec" {
                peer.Connect(&cfg.Peers[i])
            }
        }
        time.Sleep(5000 * time.Millisecond)
    }
}

// Deal with the console file io. Check console/console.go for more infomation
func handleconsole(logfile string, cfg *parsecfg.Cfg) {
    // Every .1 second read in .log
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
                if v.Status == "authenticated" || v.Status == "authsent" {
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
    listener, or := net.Listen("tcp", cfg.MyIP + ":" + strconv.Itoa(cfg.Port))
    e.Rr(or, true)
    defer listener.Close()

    go server(listener, cfg)
    go send(cfg)
    go handleconsole(logfile, cfg)
    for {}
}
