package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    bot "./bot"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: ", os.Args[0], "host:port")
        os.Exit(1)
    }
    host := os.Args[1]
    conn, err := net.Dial("tcp", host)
    checkError(err)
    client := bot.Bot{conn}
    client.Auth()
    go client.Listen()
    str := make(chan string) // chan string
    for {
        go read_input(str)
        msg := <-str
        client.Send(msg)
    }
}

func read_input(str chan string) {
    consoleReader := bufio.NewReader(os.Stdin) // *bufio.Reader
    for {
        text, _ := consoleReader.ReadString('\n')
        str <- text
    }
}

func checkError(err error) {
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
        os.Exit(1)
    }
}
