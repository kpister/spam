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
    signal.Notify(exit, os.Interrupt)
    for range exit {
        os.Remove(".log")
        // Handle exit status
        fmt.Print("\n")
        os.Exit(1)
    }
}

func main(){

    // used for exit call
    exit := make(chan os.Signal, 1)
    go handleexit(exit)

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
                return
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
        } else if v == "--help" || v == "-h" {
            fmt.Println("Here is some help:\nTo run: spam\n" +
                        "To use a new cfg file: spam -i newcfg.cfg\n" +
                        "To run the console: spam -c\nTo generate a keypair: spam --gen-keypair\n" +
                        "For further documentation, visit https://github.com/kpister/spam")
            return
        } else {
            fmt.Println("That command doesn't exist. Try spam --help to see options")
            return
        }
    }

    // Parse the config file
    crypto.SetE()
    cfg := parsecfg.ParseCfg(configfile, true)

    // Start app
    go spamcore.StartServer(logfile, cfg)
    for { }
}
