
package parsecfg

import (
    "net"
    "fmt"
    "strings"
    "strconv"
    "math/big"
    "io/ioutil"

    "github.com/kpister/spam/e"
    "github.com/kpister/spam/peer"
)

type Cfg struct {
    Peers []peer.Peer
    Port int
    SecretKey big.Int
    PublicKey big.Int
    MyIP string
}

func ParseCfg(filename string, testing bool) *Cfg {
    cfg := Cfg{}

    bytecontents, or := ioutil.ReadFile(filename)
    e.Rr(or, true)

    contents := string(bytecontents)
    pieces := strings.Split(contents, "\n")

    // Right now this is always true. We would change to false for non-localhost
    // TODO make this a parse flag?
    if !testing {
        cfg.MyIP = getMyIP()
    } else {
        cfg.MyIP = "127.0.0.1"
    }

    // Handle the file contents
    readpeers := false
    for i, v := range pieces {
        if len(v) > 0 && v[0] == byte('#') { // Comments
            continue
        }
        if v == "peers" { // Create peer list
            readpeers = true
        } else if v == "end" {
            readpeers = false
        } else if readpeers {
            ppieces := strings.Split(v, " ")
            if len(ppieces) != 3 {
                fmt.Println("Skipping peer: not formatted correctly: ip name public_key")
                continue
            }
            mpeer := peer.MakePeer(ppieces[0], ppieces[1], ppieces[2])
            cfg.Peers = append(cfg.Peers, *mpeer)
        } else if strings.Contains(v, "port") {
            port := strings.Split(v, " ")[1]
            cfg.Port, _ = strconv.Atoi(port)
            peer.SetAddr(cfg.MyIP + ":" + port)
        } else if strings.Contains(v, "secret") {
            var key big.Int
            _, suc := key.SetString(strings.Split(v, " ")[1], 10)
            if suc {
                cfg.SecretKey = key
            } else {
                fmt.Printf("Key not formatted correctly on line %d\n", i)
                continue
            }
        } else if strings.Contains(v, "public") {
            var key big.Int
            _, suc := key.SetString(strings.Split(v, " ")[1], 10)
            if suc {
                cfg.PublicKey = key
            } else {
                fmt.Printf("Key not formatted correctly on line %d\n", i)
                continue
            }
        }
    }
    return &cfg
}

// courtesy of jniltonho but this just seems really simple...
func getMyIP() string {
    addrs, or := net.InterfaceAddrs()
    e.Rr(or, true)

    for _, a := range addrs {
        if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
}
