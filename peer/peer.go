
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
    PublicKey *big.Int
}

func MakePeer(addr, name string, public string) *Peer {
    var key big.Int
    pub, suc := key.SetString(public, 10)
    conn, or := net.Dial("tcp", addr)
    stop := e.Rr(or, false)

    if !stop && suc {
        handshake(conn, pub)
        fmt.Println("Successfully connected to peer: " + conn.RemoteAddr().String())
        return &Peer{conn, name, "connected", addr, pub}
    }

    return &Peer{conn, name, "offline", addr, pub}
}

func Connect(peer *Peer) {
    conn, or := net.Dial("tcp", peer.Addr)

    if !e.Rr(or, false) {
        handshake(conn, peer.PublicKey)
        fmt.Println("Successfully connected to peer: " + conn.RemoteAddr().String())
        peer.Status = "connected"
        peer.Conn = conn
    }
}

func handshake(conn net.Conn, modulus *big.Int) bool {
    //return true if it works
    // When we dial a peer, send an encrypted (signed) message
    m := "asdfasdfd" // TODO: Set a proper message here
    // Create message
    kd := crypto.ConvertMessageToInt(m)
    message := (crypto.Encrypt(kd, modulus)).String()
    fmt.Println("Message: " + kd.String())
    fmt.Println("Encrypted: " + message)
    fmt.Fprintf(conn, "Handshake:" + message)
    // Listen for response
        // Check response
        // nm := crypto.Decrypt(response, cfg.Private, modulus)
        // if nm == m {
            // Add peer
        // }
    return true
}

