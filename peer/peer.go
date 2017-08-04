
package peer

import (
    "net"
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
        return &Peer{conn, name}
    }

    return nil
}
