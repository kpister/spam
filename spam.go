package main

import (
    "os"
    "fmt"
    "os/signal"

    "github.com/kpister/spam/crypto"
    "github.com/kpister/spam/keygen"
    "github.com/kpister/spam/console"
    "github.com/kpister/spam/spamcore"
    "github.com/kpister/spam/parsecfg"
)

func handleexit(exit chan os.Signal) {
    for range exit {
        os.Remove(".log")
        // Handle exit status
        fmt.Print("\n")
        os.Exit(1)
    }
}

func main(){

    exitchannel := make(chan os.Signal, 1)
    signal.Notify(exitchannel, os.Interrupt)
    go handleexit(exitchannel)

    configfile := "spam_core.cfg"
    logfile := ".log"

    // Search for command flags
    for i, v := range os.Args {
        // -i to set config file
        if v == "-i" {
            if len(os.Args) > i {
                configfile = os.Args[i + 1]
                logfile = ".log_" + configfile
            } else {
                fmt.Println("You failed to run this.")
                os.Exit(0)
            }
        } else if v == "-c" {
            if len(os.Args) > i + 1 {
                logfile = os.Args[i + 1]
            }
            defer console.Start(logfile)
            return
        } else if v == "--gen-keypair" {
            keygen.GenKeys()
            return
        }
    }

    crypto.SetE()
    cfg := parsecfg.ParseCfg(configfile)

    spamcore.StartServer(logfile, cfg)
}
