package bot

import (
    "bufio"
    "fmt"
    "os"
    "net"
    "strings"
    "time"
    "encoding/json"
    "io/ioutil"
)

type Bot struct {
    Conn net.Conn
    config Config
    privmsg string
    nick_cmd string
    user_cmd string
    eof string
}

type Config struct {
    Server string
    Nick string
    Ident string
    Hostmask string
    NSPassword string
    Channels []string
    Commands map[string]string
}

func InitBot(jsonFile string) Bot {
    bot := Bot{}
    var config Config
    file, _ := ioutil.ReadFile(jsonFile)
    err := json.Unmarshal(file, &config)
    if err != nil {
        fmt.Println("Configuration Error: ", err.Error())
    }
    bot.config = config
    bot.privmsg = "PRIVMSG"
    bot.nick_cmd = "NICK"
    bot.user_cmd = "USER"
    bot.eof = "EOF"
    return bot
}

func (b *Bot) Connect() {
    conn, err := net.Dial("tcp", b.config.Server)
    checkError(err)
    b.Conn = conn
}

func (b Bot) Auth() {
    user_cmd := fmt.Sprintf(
        "%s %s %s %s :%s",
        b.user_cmd,
        b.config.Ident,
        b.config.Ident,
        b.config.Hostmask,
        b.config.Nick)
    nick_cmd := fmt.Sprintf("%s %s", b.nick_cmd, b.config.Nick)
    id_cmd := fmt.Sprintf("%s nickserv identify %s", b.privmsg, b.config.NSPassword)
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
            if err.Error() == b.eof {
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
        fmt.Println("Fatal error: ", err.Error())
        os.Exit(1)
    }
}
