package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "discord bot token")
	flag.Parse()
	if Token == "" {
		flag.Usage()
		return
	}

}

func main() {

	session, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalln("Couldn't start new session")
	}

	err = session.Open()
	if err != nil {
		log.Fatalln("Session couldn't opened")
		return
	}

	fmt.Println("Discord Bot running ctrl + c for kill")

	session.AddHandler(messageHandler)
	session.AddHandler(followupHandler)

	registeredCmds := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, "", v)
		if err != nil {
			log.Fatalf("Couldn't recognized current command %v", err)
			return
		}

		registeredCmds[i] = cmd
	}

	defer session.Close()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	msgs := strings.Fields(m.Content)
	if len(msgs) <= 1 {
		s.ChannelMessageSend(m.ChannelID, "args must be more than one")
		return
	}

	for i, item := range msgs {
		if strings.Contains(item, "https") {
			// this is for eamil; handle for donwlaod image;r
			avatar, err := getImage(item)
			if err != nil {
				log.Fatalf("an error accured try again later %v", err)
			}

			name := msgs[len(msgs)-1-i]

			user, err := s.UserUpdate(name, avatar)
			if err != nil {
				log.Fatalf("Couldn't update user %v", err)
			}

			fmt.Printf("The user %s", user)

		}

	}

}

func getImage(url string) (string, error) {
	var base64img string
	var contentType string

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	img, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	contentType = http.DetectContentType(img)
	base64img = base64.StdEncoding.EncodeToString(img)

	return fmt.Sprintf("data:%s;base64,%s", contentType, base64img), nil

}

var (
	commands = []*discordgo.ApplicationCommand{{
		Name:        "followups",
		Description: "Followups messages",
	}}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"followups": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: "Suprise!",
				},
			})
			msg, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "Follwup message has been created, after 5 seconds it will be edited",
			})

			if err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went wrong",
				})
				return
			}

			time.Sleep(time.Second * 5)

			content := "Now the original message is gone and after 10 seconds this message will ~~self-destruct~~ be deleted."
			s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
				Content: &content,
			})

			time.Sleep(time.Second * 10)
			s.FollowupMessageDelete(i.Interaction, msg.ID)
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "For those, who didn't skip anything and followed tutorial along fairly, " +
					"take a unicorn :unicorn: as reward!\n" +
					"Also, as bonus... look at the original interaction response :D",
			})

		},
	}
)

func followupHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if hook, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
		hook(s, i)
	}

}
