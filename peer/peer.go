
package peer

import (
    "net"
    "fmt"
    "github.com/kpister/spam/e"
)

type Peer struct {
    Conn net.Conn
    Name string
    Status string
    Addr string
}

func MakePeer(addr, name string) *Peer {
    conn, or := net.Dial("tcp", addr)
    stop := e.Rr(or, false)

    if !stop {
        fmt.Println("Successfully connected to peer: " + conn.RemoteAddr().String())
        return &Peer{conn, name, "connected", addr}
    }

    return &Peer{conn, name, "offline", addr}
}

func Connect(peer *Peer) {
    conn, or := net.Dial("tcp", peer.Addr)

    if !e.Rr(or, false) {
        fmt.Println("Successfully connected to peer: " + conn.RemoteAddr().String())
        peer.Status = "connected"
        peer.Conn = conn
    }
}
