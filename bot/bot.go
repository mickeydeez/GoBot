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
    conn net.Conn
}

func (b Bot) Connect(host string) {
    conn, err := net.Dial("tcp", host)
    checkError(err)
    b.conn = conn
}

func (b Bot) Auth() int {
    b.conn.Write([]byte("USER testbot testbot yourmom.com :testbot\n"))
    time.Sleep(2)
    b.conn.Write([]byte("NICK testbot\n"))
    time.Sleep(2)
    b.conn.Write([]byte("PRIVMSG nickserv identify somepassword"))
    return 0
}

func (b Bot) Listen() {
    socketListener := bufio.NewReader(b.conn) // *bufio.Reader
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
    b.conn.Write([]byte(msg))
}

func checkError(err error) {
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
        os.Exit(1)
    }
}
