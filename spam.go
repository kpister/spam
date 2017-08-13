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

// An async function to handle ^C input. This way we can cancel any import stuff
// TODO figure out what we actually need to close/delete
func handleexit(exit chan os.Signal) {
    for range exit {
        os.Remove(".log")
        // Handle exit status
        fmt.Print("\n")
        os.Exit(1)
    }
}

func main(){

    // used for exit call
    exitchannel := make(chan os.Signal, 1)
    signal.Notify(exitchannel, os.Interrupt)
    go handleexit(exitchannel)

    // default values, can be changed with flags
    configfile := "spam_core.cfg"
    logfile := ".log"

    // Search for command flags
    skip := true
    for i, v := range os.Args {
        if skip {
            skip = false
        } else if v == "-i" {
            if len(os.Args) > i {
                configfile = os.Args[i + 1]
                logfile = ".log_" + configfile
                skip = true
            } else {
                fmt.Println("You failed to run this.")
                os.Exit(0)
            }
        } else if v == "-c" {
            if len(os.Args) > i + 1 {
                logfile = os.Args[i + 1]
                skip = true
            }
            defer console.Start(logfile)
            return
        } else if v == "--gen-keypair" {
            keygen.GenKeys()
            return
        } else {
            fmt.Println("That command doesn't exist")
            return
        }
    }

    // Parse the config file
    crypto.SetE()
    cfg := parsecfg.ParseCfg(configfile, true)

    // Start app
    spamcore.StartServer(logfile, cfg)
}
