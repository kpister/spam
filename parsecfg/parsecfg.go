
package parsecfg

import (
    "fmt"
    "strings"
    "strconv"
    "io/ioutil"

    "github.com/kpister/spam/e"
    "github.com/kpister/spam/peer"
)

type Cfg struct {
    Peers []peer.Peer
    Port int
    Secret string
}

func ParseCfg(filename string) *Cfg {
    cfg := Cfg{}

    bytecontents, or := ioutil.ReadFile(filename)
    e.Rr(or, true)

    contents := string(bytecontents)
    pieces := strings.Split(contents, "\n")

    readpeers := false
    for _, v := range pieces {
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
            cfg.Secret = strings.Split(v, " ")[1]
        }

    }
    return &cfg
}

