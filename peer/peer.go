
package peer

import (
    "net"
    "fmt"
    "github.com/kpister/spam/e"
)

type Peer struct {
    Conn net.Conn
    Name string
}

func MakePeer(addr, name string) *Peer {
    conn, or := net.Dial("tcp", addr)
    stop := e.Rr(or, false)

    if !stop {
        fmt.Println("Successfully connected to peer: " + conn.RemoteAddr().String())
        return &Peer{conn, name}
    }

    return nil
}
