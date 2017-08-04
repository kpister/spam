
package parsecfg

import (
    "io/ioutil"
    "net"
    "strings"
    "github.com/kpister/spam/e"
)

type Peer struct {
    Conn net.Conn
    Name string
}

type Cfg struct {
    Peers []Peer
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
        } else if readpeers {
            ppieces := strings.Split(v, " ")
            pname := ""
            pconn, or := net.Dial("tcp", ppieces[0])
            e.Rr(or, false)
            if len(ppieces) > 1 {
                pname = ppieces[1]
            }
            cfg.Peers = append(cfg.Peers, Peer{pconn, pname})
        } else if v == "end" {
            readpeers = false
        }

    }
    return cfg
}

