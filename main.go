package main

import (
    "bufio"
    "fmt"
    "os"
    bot "github.com/mickeydeez/GoBot/bot"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: ", os.Args[0], "host:port")
        os.Exit(1)
    }
    host := os.Args[1]
    var client bot.Bot
    client.Connect(host)
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

