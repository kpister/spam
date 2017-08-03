
package parsecfg

import (
    "net"
)

type Cfg struct {
    peers []net.Conn
}

func ParseCfg(filename string) Cfg {
    cfg := Cfg{}

    return cfg
}

