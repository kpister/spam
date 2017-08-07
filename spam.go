package main

import (
    "os"
    "os/signal"
    "fmt"
    "github.com/kpister/spam/spamcore"
    "github.com/kpister/spam/parsecfg"
    "github.com/kpister/spam/console"
)

func handleexit(exit chan os.Signal) {
    for signal := range exit {
        os.Remove(".log")
        // Handle exit status
        os.Exit(1)
    }
}

func main(){

    exitchannel := make(chan os.Signal, 1)
    signal.Notify(exitchannel, os.Interrupt)
    go handleexit(exitchannel)

    configfile := "spam_core.cfg"

    // Search for command flags
    for i, v := range os.Args {
        // -i to set config file
        if v == "-i" {
            if len(os.Args) > i {
                configfile = os.Args[i + 1]
            } else {
                fmt.Println("You failed to run this.")
                os.Exit(0)
            }
        } else if v == "-c" {
            defer console.Start()
            return
        }
    }

    os.Create(".log")

    cfg := parsecfg.ParseCfg(configfile)

    spamcore.StartServer(cfg)
}
