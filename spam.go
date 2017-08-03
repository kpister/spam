package main

import (
    "github.com/kpister/spam/spamcore"
    "github.com/kpister/spam/parsecfg"
)

func main(){
    cfg := parsecfg.ParseCfg("spam_core.cfg")

    spamcore.StartServer(&cfg)
}
