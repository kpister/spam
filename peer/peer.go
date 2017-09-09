
package peer

import (
    "bytes"
	"fmt"
    "net"
	
	"math/big"

    "github.com/kpister/spam/e"
    "github.com/kpister/spam/crypto"
    //"github.com/kpister/spam/parsecfg"
)

/* A peer has the following statuses
 * 
 * offline: no connection
 * authsent: we have sent them a handshake (they know who we are)
 * authrec: we received a handshake (we know who they are)
 * authenticated: we have sent and received handshakes
 *
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
var mprivKey big.Int
var mpubKey big.Int
var mPrime1 big.Int
var mPrime2 big.Int

func SetAddr(addr string) {
    maddr = addr
}

func SetPrivKey(key big.Int) {
	mprivKey = key
}

func SetPubKey(key big.Int) {
	mpubKey = key
}

func SetPrime1(p big.Int) {
	mPrime1 = p
}

func SetPrime2(p big.Int) {
	mPrime2 = p
}

// Used when first created a peer (through console or through parsecfg)
func MakePeer(addr, name string, public string) *Peer {
    var key big.Int
    m := maddr // TODO: Set a proper message here
    pub, success := key.SetString(public, 10)
    conn, or := net.Dial("tcp", addr)

    if !e.Rr(or, false) && success {
        handshake(conn, pub, m)
        fmt.Println("Successfully connected to peer: " + conn.RemoteAddr().String())
        return &Peer{conn, name, "authsent", addr, pub, m, ""}
    }

    return &Peer{conn, name, "offline", addr, pub, m, ""}
}

// Used when trying to connect to a peer later (during the 5 second refreshes)
func Connect(peer *Peer) {
    conn, or := net.Dial("tcp", peer.Addr)
    m := maddr // TODO: Set a proper message here

    if !e.Rr(or, false) {
        handshake(conn, peer.PublicKey, m)
        fmt.Println("Successfully connected to peer: " + conn.RemoteAddr().String())
        // This is how we distinguish between if they have contacted us or not
        // Once we have connected to them, we can send them messages (they might ignore us though)
        if peer.Status == "authrec" {
            peer.Status = "authenticated"
        } else {
            peer.Status = "authsent"
        }
        peer.Conn = conn
    }
}

// Send the first shake
func handshake(conn net.Conn, modulus *big.Int, m string) {
    // When we dial a peer, send an encrypted (signed) message
	
	signature, err:= crypto.Sign(&mprivKey, &mpubKey, &mPrime1, &mPrime2, m)
	if err != nil {
    	fmt.Println("Signing failed.")
	} else {
		fmt.Println(signature)
		fmt.Println("Message signed successfully")
	}	
	var buffer bytes.Buffer
	buffer.WriteString(m)
	//buffer.Write(signature)
	//signedmessage := buffer.String()
	//intmessage := crypto.ConvertMessageToInt(signedmessage)
	intmessage := crypto.ConvertMessageToInt(m)
	// TODO sign message before encryption
	message := (crypto.Encrypt(intmessage, modulus)).String()
    fmt.Fprintf(conn, "Handshake:" + message + "\n")
}

