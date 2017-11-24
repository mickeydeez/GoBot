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

func (b Bot) Auth() int {
    b.Conn.Write([]byte("USER testbot testbot yourmom.com :testbot\n"))
    time.Sleep(2)
    b.Conn.Write([]byte("NICK testbot\n"))
    time.Sleep(2)
    b.Conn.Write([]byte("PRIVMSG nickserv identify somepassword"))
    return 0
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
