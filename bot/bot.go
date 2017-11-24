package bot

import (
    "bufio"
    "fmt"
    "os"
    "net"
    "strings"
    "time"
)

type Bot struct {
    Conn net.Conn
}

func (b *Bot) Connect(host string) {
    conn, err := net.Dial("tcp", host)
    checkError(err)
    b.Conn = conn
}

func (b Bot) Auth() {
    b.Conn.Write([]byte("USER testbot testbot yourmom.com :testbot\n"))
    time.Sleep(2)
    b.Conn.Write([]byte("NICK testbot\n"))
    time.Sleep(2)
    b.Conn.Write([]byte("PRIVMSG nickserv identify somepassword"))
}

func (b Bot) Listen() {
    socketListener := bufio.NewReader(b.Conn) // *bufio.Reader
    for {
        msg, err := socketListener.ReadString('\n')
        fmt.Println(strings.TrimSpace(msg))
        if err != nil {
            if err.Error() == "EOF" {
                os.Exit(0)
            }
        }
    }
}

func (b Bot) Send(msg string) {
    b.Conn.Write([]byte(msg))
}

func checkError(err error) {
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
        os.Exit(1)
    }
}
