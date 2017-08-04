
package parsecfg

import (
    "io/ioutil"
    "strings"
    "strconv"

    "github.com/kpister/spam/e"
    "github.com/kpister/spam/peer"
)

type Cfg struct {
    Peers []peer.Peer
    Port int
}

func ParseCfg(filename string) Cfg {
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
            pname := ""
            e.Rr(or, false)
            if len(ppieces) > 1 {
                pname = ppieces[1]
            }
            mpeer := peer.MakePeer(ppieces[0], pname)
            if mpeer != nil {
                cfg.Peers = append(cfg.Peers, *mpeer)
            }
        } else if strings.Contains(v, "port") {
            cfg.Port, _ = strconv.Atoi(strings.Split(v, " ")[1])
        }

    }
    return cfg
}

