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
	Server         string
	Nick           string
	Ident          string
	Hostmask       string
	NSPassword     string
	Admins         []string
	Channels       []string
	CommandTrigger string
	Commands       map[string]string
	Courses        []map[string]string
}

type Event struct {
	src     string
	dest    string
	cmd     string
	content []string
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

func (b *Bot) Connect() {
	conn, err := net.Dial("tcp", b.config.Server)
	checkError(err)
	b.Conn = conn
}

func (b *Bot) Bootstrap() {
	b.Auth()
	b.JoinChannels()
}

func (b Bot) Auth() bool {
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
	time.Sleep(5 * time.Second)
	b.Send(nick_cmd)
	time.Sleep(5 * time.Second)
	b.Send(id_cmd)
	time.Sleep(5 * time.Second)
	return true
}

func (b Bot) JoinChannels() {
	for _, channel := range b.config.Channels {
		b.Join(channel)
	}
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
		b.ProcessEvent(event)
	}
}

func (b Bot) Pong(msg string) {
	response := fmt.Sprintf("PONG %s", msg)
	b.Send(response)
	fmt.Println("** Responded to Server Ping **")
}

func (b Bot) Join(channel string) {
	payload := fmt.Sprintf("JOIN %s", channel)
	b.Send(payload)
}

func (b Bot) Send(msg string) {
	text := fmt.Sprintf("%s\n", msg)
	b.Conn.Write([]byte(text))
}

func (b Bot) ProcessEvent(event Event) {
	if event.cmd == "PING" {
		b.Pong(event.content[0][1:])
	} else if event.cmd == "JOIN" {
		b.ProcessJoin(event)
	} else if len(event.content) > 0 {
		if len(event.content[0]) > 1 {
			debug_out := fmt.Sprintf("<%s/%s> %s",
				event.dest, event.src, event.content)
			status, cmd := b.CheckCommand(event.content)
			if status {
				fmt.Println(debug_out)
				response := fmt.Sprintf("%s %s %s",
					b.privmsg, event.dest, b.config.Commands[cmd])
				b.Send(response)
			} else {
				fmt.Println(debug_out)
			}
		}
	}
}

func (b Bot) ProcessJoin(event Event) {
	if b.CheckAdmin(event.src) {
		nick, _, _ := parseSource(event.src)
		response := fmt.Sprintf("%s %s %s the badass has entered the room",
			b.privmsg, event.dest[1:], nick)
		b.Send(response)
	}
}

func (b Bot) CheckAdmin(event_src string) bool {
	// returns bool
	nick, _, _ := parseSource(event_src)
	for _, item := range b.config.Admins {
		if nick == item {
			return true
		}
	}
	return false
}

func (b Bot) CheckCommand(event_content []string) (bool, string) {
	// returns bool, command
	if string(event_content[0][1]) == string(b.config.CommandTrigger) {
		for k, _ := range b.config.Commands {
			if k == strings.Split(event_content[0], b.config.CommandTrigger)[1] {
				return true, k
			}
		}
	}
	return false, ""
}

func parseSource(event_src string) (string, string, string) {
	// returns nick, ident, vhost
	split := strings.Split(event_src, "!")
	nick := split[0]
	host_split := strings.Split(split[1], "@")
	return nick, host_split[0], host_split[1]
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error: ", err.Error())
		os.Exit(1)
	}
}

func Extend(slice []string, element string) []string {
	n := len(slice)
	if n == cap(slice) {
		newSlice := make([]string, len(slice), 2*len(slice)+1)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0 : n+1]
	slice[n] = element
	return slice
}

func ParseEvent(msg string) Event {
	var event Event
	words := strings.Fields(strings.TrimSpace(msg))
	content := make([]string, 0)
	if words[0] == "PING" {
		event.src = ""
		event.dest = ""
		event.cmd = words[0]
		content = Extend(content, strings.Join(words[1:], " "))
		event.content = content
	} else {
		event.src = words[0][1:]
		event.cmd = words[1]
		event.dest = words[2]
		if len(words) >= 4 {
			for _, item := range words[3:] {
				content = Extend(content, item)
			}
		}
		event.content = content
	}
	return event
}
