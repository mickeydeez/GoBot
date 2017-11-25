package main

import (
	"bufio"
	"fmt"
	bot "github.com/mickeydeez/GoBot/bot"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "<JSON Config>")
		os.Exit(1)
	}
	jsonFile := os.Args[1]
	client := bot.InitBot(jsonFile)
	client.Connect()
	client.Auth()
	go client.Listen()
	str := make(chan string)
	for {
		go read_input(str)
		msg := <-str
		client.Send(msg)
	}
}

func read_input(str chan string) {
	consoleReader := bufio.NewReader(os.Stdin)
	for {
		text, _ := consoleReader.ReadString('\n')
		str <- text
	}
}
