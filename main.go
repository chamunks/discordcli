// This file provides a basic "quick start" example of using the Discordgo
// package to connect to Discord using the New() helper function.
package main

import (
	"log"
	"regexp"

	"github.com/chzyer/readline"
	"github.com/theboxmage/DiscordCli/DiscordState"
)

//Global Message Types
const (
	ErrorMsg  = "Error"
	InfoMsg   = "Info"
	HeaderMsg = "Head"
	TextMsg   = "Text"
)

//Version is current version const
const Version = "v1.4.2 - Box Develop~"

//Session is global Session
var Session *DiscordState.Session

//State is global State
var State *DiscordState.State

//UserChannels is global User Channels

//MsgType is a string containing global message type
type MsgType string

func main() {
	//Initialize Config
	GetConfig()
	CheckState()
	Clear()
	Msg(HeaderMsg, "discord-cli - version: %s\n\n", Version)

	//NewSession
	Session = DiscordState.NewSession(Config.Username, Config.Password) //Please don't abuse
	err := Session.Start()
	if err != nil {
		log.Println("Session Failed")
		log.Fatalln(err)
	}
	//Attach New Window
	InitWindow()

	//Attach Even Handlers
	State.Session.DiscordGo.AddHandler(newMessage)
	//State.Session.DiscordGo.AddHandler(newReaction)
	//Setup Readline
	rl, err := readline.NewEx(&readline.Config{
		Prompt:         "> ",
		UniqueEditLine: true,
	})

	defer rl.Close()
	log.SetOutput(rl.Stderr()) // let "log" write to l.Stderr instead of os.Stderr
	State.Session.DiscordGo.UpdateStatus(0, "discord-cli")

	//Start Listening
	for {
		line, _ := rl.Readline()

		//QUIT
		if line == ":q" {
			break
		}

		//Parse Commands
		line = ParseForCommands(line)

		line = ParseForMentions(line)

		if line != "" {
			State.Session.DiscordGo.ChannelMessageSend(State.Channel.ID, line)
		}
	}

	return
}

//InitWindow creates a New CLI Window
func InitWindow() {
	SelectGuildMenu()
	if State.Channel == nil {
		SelectChannelMenu()
	}
	State.Enabled = true
	ShowContent()
}

//ShowContent shows defaulth Channel content
func ShowContent() {
	Clear()
	Header()
	if Config.MessageDefault {
		State.RetrieveMessages(Config.Messages)
		PrintMessages(Config.Messages)
	}
}

//ShowEmptyContent shows an empty channel
func ShowEmptyContent() {
	Clear()
	Header()
}

//ParseForMentions parses input string for mentions
func ParseForMentions(line string) string {
	r, err := regexp.Compile("\\@\\w+")
	if err != nil {
		Msg(ErrorMsg, "Regex Error: ", err)
	}

	lineByte := r.ReplaceAllFunc([]byte(line), ReplaceMentions)

	return string(lineByte[:])
}

//ReplaceMentions replaces mentions to ID
func ReplaceMentions(input []byte) []byte {
	var OutputString string

	SizeByte := len(input)
	InputString := string(input[1:SizeByte])

	if Member, ok := State.Members[InputString]; ok {
		OutputString = "<@" + Member.User.ID + ">"
	} else {
		OutputString = "@" + InputString
	}
	return []byte(OutputString)
}
