package main

import (
//    "os"
//    "os/signal"
//    "fmt"
    "github.com/kpister/spam/spamcore"
    "github.com/kpister/spam/parsecfg"
)
/*
func handleexit(exit chan os.Signal) {
    for signal := range exit {
        // Handle exit status
        fmt.Println(signal)
        os.Exit(1)
    }
}
*/
func main(){
    /*
    exitchannel := make(chan os.Signal, 1)
    signal.Notify(exitchannel, os.Interrupt)
    go handleexit(exitchannel)
    */

    cfg := parsecfg.ParseCfg("spam_core.cfg")

    spamcore.StartServer(&cfg)
}
