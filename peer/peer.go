
package peer

import (
    "net"
    "fmt"
    "math/big"

    "github.com/kpister/spam/e"
    "github.com/kpister/spam/crypto"
)

/* A peer has the following statuses
 * 
 * offline: no connection
 * authsent: we have sent them a handshake (they know who we are)
 * authrec: we received a handshake (we know who they are)
 * authenticated: we have sent and received handshakes
 */

type Peer struct {
    Conn net.Conn
    Name string
    Status string
    Addr string
    PublicKey *big.Int
    SecretMessage string
    RemoteAddr string
}

var maddr string

func SetAddr(addr string) {
    maddr = addr
}

func MakePeer(addr, name string, public string) *Peer {
    var key big.Int
    pub, suc := key.SetString(public, 10)
    conn, or := net.Dial("tcp", addr)
    stop := e.Rr(or, false)
    m := maddr // TODO: Set a proper message here

    if !stop && suc {
        handshake(conn, pub, m)
        fmt.Println("Successfully connected to peer: " + conn.RemoteAddr().String())
        return &Peer{conn, name, "authsent", addr, pub, m, ""}
    }

    return &Peer{conn, name, "offline", addr, pub, m, ""}
}

func Connect(peer *Peer) {
    conn, or := net.Dial("tcp", peer.Addr)
    m := maddr // TODO: Set a proper message here

    if !e.Rr(or, false) {
        handshake(conn, peer.PublicKey, m)
        fmt.Println("Successfully connected to peer: " + conn.RemoteAddr().String())
        if peer.Status == "authrec" {
            peer.Status = "authenticated"
        } else {
            peer.Status = "authsent"
        }
        peer.Conn = conn
    }
}

func handshake(conn net.Conn, modulus *big.Int, m string) {
    // When we dial a peer, send an encrypted (signed) message
    kd := crypto.ConvertMessageToInt(m)
    message := (crypto.Encrypt(kd, modulus)).String()
    fmt.Fprintf(conn, "Handshake:" + message + "\n")
}

