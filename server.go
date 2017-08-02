package main

// Source: http://divan.github.io/posts/go_concurrency_visualize/

import (
    "fmt"
    "net"
    "time"
    "bufio"
)

func handler(c net.Conn, ch chan string) {
    ch <- c.RemoteAddr().String()
    for {
        message, _ := bufio.NewReader(c).ReadString('\n')
        // output message received
        fmt.Print("Message Received:", string(message))
        // send new string back to client
        c.Write([]byte(message + "\n"))

    }
    c.Close()
}

func logger(ch chan string) {
    for {
        fmt.Println(<-ch)
    }
}

func server(l net.Listener, ch chan string) {
    for {
        c, err := l.Accept()
        if err != nil {
            continue
        }
        go handler(c, ch)
    }
}

func main() {
    l, err := net.Listen("tcp", ":8080")
    if err != nil {
        panic(err)
    }
    ch := make(chan string)
    go logger(ch)
    go server(l, ch)
    for {
        fmt.Print(".")
        time.Sleep(time.Millisecond * 1000)
    }
}
