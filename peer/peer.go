
package peer

import (
    "net"
    "fmt"
    "math/big"

    "github.com/kpister/spam/e"
    "github.com/kpister/spam/crypto"
)

type Peer struct {
    Conn net.Conn
    Name string
    Status string
    Addr string
    PublicKey big.Int
}

func MakePeer(addr, name string, public *Int) *Peer {
    conn, or := net.Dial("tcp", addr)
    stop := e.Rr(or, false)
    handshake(conn, public)

    if !stop {
        fmt.Println("Successfully connected to peer: " + conn.RemoteAddr().String())
        return &Peer{conn, name, "connected", addr, public}
    }

    return &Peer{conn, name, "offline", addr, public}
}

func Connect(peer *Peer) {
    conn, or := net.Dial("tcp", peer.Addr)
    handshake(peer.Conn, peer.Modulus)

    if !e.Rr(or, false) {
        fmt.Println("Successfully connected to peer: " + conn.RemoteAddr().String())
        peer.Status = "connected"
        peer.Conn = conn
    }
}

func handshake(conn net.Conn, modulus big.Int) bool {
    //return true if it works
    // TODO When we dial a peer, send an encrypted (signed) message
    m := ""
    // Create message
    fmt.Fprintf(conn, crypto.Encrypt(crypto.ConvertMessageToInt(m), modulus))
    // Listen for response
        // Check response
        // nm := crypto.Decrypt(response, cfg.Private, modulus)
        // if nm == m {
            // Add peer
        // }
    return true
}

