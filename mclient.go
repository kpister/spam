package main

import "net"
import "fmt"
import "bufio"
import "os"
import "reflect"

func Client() {

    reader := bufio.NewReader(os.Stdin)
    fmt.Println(reflect.TypeOf(reader))
    // connect to this socket
    conn, _ := net.Dial("tcp", "127.0.0.1:8080")
    for {
        // read in input from stdin
        fmt.Print("Text to send: ")
        text, _ := reader.ReadString('\n')
        // send to socket
        fmt.Fprintf(conn, text + "\n")
        // listen for reply
        message, _ := bufio.NewReader(conn).ReadString('\n')
        fmt.Println("Message from server: "+message)
    }
}
