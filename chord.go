package chord

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type responseMethodType func() string

// ResponsesMapType stores the prompts and methods to return a response
type ResponsesMapType map[string]responseMethodType

// BackgroundProcessType background process with ability to send messages
type BackgroundProcessType func(messages chan<- string)

// Chord a discord bot exposing backgroundProcess & responder
type Chord struct {
	authenticationToken string
	guild               string
	channel             string
	session             *discordgo.Session
	responses           ResponsesMapType
	process             BackgroundProcessType
}

var (
	chd *Chord
)

// NewChord creates a new golang chord
func NewChord(
	authenticationToken string,
	guild string,
	channel string,
	responses ResponsesMapType,
	process BackgroundProcessType,
) *Chord {
	chd = &Chord{
		authenticationToken: authenticationToken,
		guild:               guild,
		channel:             channel,
		responses:           responses,
		process:             process,
	}
	sesh, err := discordgo.New("Bot " + authenticationToken)
	if err != nil {
		log.Fatal(err)
	}
	chd.session = sesh

	chd.session.AddHandler(ready)
	chd.session.AddHandler(messageResponder)

	chd.session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	err = chd.session.Open()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Bot is running. ctrl-c to exit.")

	// create a channel to pass information from testBackgroundProcess to a process waiting to send messages to the guild
	messages := make(chan string)
	channelID, err := getChannelID(chd)
	if err != nil {
		log.Fatal(err)
	}
	// start processes
	go messageSenderUnprompted(chd.session, channelID, messages)
	go chd.process(messages)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	chd.session.Close()
	return chd
}

func messageResponder(s *discordgo.Session, m *discordgo.MessageCreate) {
	// don't do anything if the message is from this bot
	if m.Author.ID == s.State.User.ID {
		return
	}
	for k, f := range chd.responses {
		if k == strings.ToLower(m.Content) {
			s.ChannelMessageSend(m.ChannelID, f())
		}
	}
}

// background listener simply receives string messages and sends to discord channel
func messageSenderUnprompted(s *discordgo.Session, channelID string, messages chan string) {
	for {
		msg := <-messages
		s.ChannelMessageSend(channelID, msg)
	}
}

func ready(s *discordgo.Session, r *discordgo.Ready) {
	fmt.Println("Bot is Ready")
}

// might not be necessary
// for some reason s.State.Guilds[0].Name is unpopulated for me
// this is true even after discordgo.Ready event
func getChannelID(chd *Chord) (string, error) {
	for _, guild := range chd.session.State.Guilds {
		guild, err := chd.session.Guild(guild.ID)
		if err != nil {
			return "", err
		}
		if guild.Name == chd.guild {
			channels, err := chd.session.GuildChannels(guild.ID)
			if err != nil {
				return "", err
			}
			for _, channel := range channels {
				if channel.Name == chd.channel {
					return channel.ID, nil
				}
			}
		}
	}
	return "", errors.New("could not obtain channelID, ensure guild & channel names are correct")
}
