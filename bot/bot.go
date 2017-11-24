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
    Host string
    privmsg string
    nick_cmd string
    user_cmd string
}

func InitBot(host string) Bot {
    bot := Bot{}
    bot.Host = host
    bot.privmsg = "PRIVMSG"
    bot.nick_cmd = "NICK"
    bot.user_cmd = "USER"
    return bot
}

func (b *Bot) Connect() {
    conn, err := net.Dial("tcp", b.Host)
    checkError(err)
    b.Conn = conn
}

func (b Bot) Auth() {
    user_cmd := fmt.Sprintf("%s testbot testbot yourmom.com :testbot", b.user_cmd)
    nick_cmd := fmt.Sprintf("%s testbot", b.nick_cmd)
    id_cmd := fmt.Sprintf("%s nickserv identify somepassword", b.privmsg)
    b.Send(user_cmd)
    time.Sleep(2)
    b.Send(nick_cmd)
    time.Sleep(2)
    b.Send(id_cmd)
}

func (b Bot) Listen() {
    socketListener := bufio.NewReader(b.Conn) // *bufio.Reader
    for {
        msg, err := socketListener.ReadString('\n')
        if err != nil {
            if err.Error() == "EOF" {
                os.Exit(0)
            }
        }
        words := strings.Fields(strings.TrimSpace(msg))
        // cmd := words[1]
        // dest := words[2]
        src := words[0][1:]
        content := strings.Join(words[3:], " ")
        text := content[1:]
        line := fmt.Sprintf("<%s> %s", src, text)
        fmt.Println(line)
    }
}

func (b Bot) Send(msg string) {
    text := fmt.Sprintf("%s\n", msg)
    b.Conn.Write([]byte(text))
}

func checkError(err error) {
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
        os.Exit(1)
    }
}
