package bot

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

type Bot struct {
	Conn         net.Conn
	config       Config
	privmsg      string
	nick_cmd     string
	user_cmd     string
	nickserv_cmd string
	eof          string
}

type Config struct {
	Server     string
	Nick       string
	Ident      string
	Hostmask   string
	NSPassword string
	Channels   []string
	Commands   map[string]string
	Courses    []map[string]string
}

type Event struct {
	src     string
	dest    string
	cmd     string
	content string
}

func InitBot(jsonFile string) Bot {
	bot := Bot{}
	var config Config
	file, _ := ioutil.ReadFile(jsonFile)
	err := json.Unmarshal(file, &config)
	checkError(err)
	bot.config = config
	bot.privmsg = "PRIVMSG"
	bot.nick_cmd = "NICK"
	bot.user_cmd = "USER"
	bot.nickserv_cmd = "nickserv identify"
	bot.eof = "EOF"
	return bot
}

func ParseEvent(msg string) Event {
	var event Event
	words := strings.Fields(strings.TrimSpace(msg))
	if words[0] == "PING" {
		event.src = ""
		event.dest = ""
		event.cmd = words[0]
		event.content = strings.Join(words[1:], " ")
	} else {
		event.src = words[0][1:]
		event.cmd = words[1]
		event.dest = words[2]
		if len(words) >= 4 {
			event.content = strings.Join(words[3:], " ")
		} else {
			event.content = ""
		}
	}
	return event
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
	nick_cmd := fmt.Sprintf(
		"%s %s",
		b.nick_cmd,
		b.config.Nick)
	id_cmd := fmt.Sprintf(
		"%s %s %s",
		b.privmsg,
		b.nickserv_cmd,
		b.config.NSPassword)
	b.Send(user_cmd)
	time.Sleep(2)
	b.Send(nick_cmd)
	time.Sleep(2)
	b.Send(id_cmd)
}

func (b Bot) Listen() {
	socketListener := bufio.NewReader(b.Conn)
	for {
		msg, err := socketListener.ReadString('\n')
		if err != nil {
			if err.Error() == b.eof {
				os.Exit(0)
			}
		}
		event := ParseEvent(msg)
		if event.cmd == "PING" {
			response := fmt.Sprintf("PONG %s", event.content)
			b.Send(response)
			fmt.Println("Responded to Ping")
		} else {
			line := fmt.Sprintf("Src: %s\nDest: %s\nCmd: %s\nContent: %s\n\n",
				event.src, event.dest, event.cmd, event.content)
			fmt.Printf(line)
		}
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
