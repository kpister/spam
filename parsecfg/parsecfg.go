
package parsecfg

import (
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
    Secret big.Int
}

func ParseCfg(filename string) *Cfg {
    cfg := Cfg{}

    bytecontents, or := ioutil.ReadFile(filename)
    e.Rr(or, true)

    contents := string(bytecontents)
    pieces := strings.Split(contents, "\n")

    readpeers := false
    for i, v := range pieces {
        if v == "peers" {
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
            cfg.Port, _ = strconv.Atoi(strings.Split(v, " ")[1])
        } else if strings.Contains(v, "secret") {
            var key big.Int
            _, suc := key.SetString(strings.Split(v, " ")[1], 10)
            if suc {
                cfg.Secret = key
            } else {
                fmt.Printf("Key not formatted correctly on line %d\n", i)
                continue
            }

        }

    }
    return &cfg
}

