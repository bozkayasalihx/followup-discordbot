package main

import (
	"bufio"
	browsercontext "discord-bot/browserContext"
	"discord-bot/handler"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/playwright-community/playwright-go"
)

// var (
// 	Token   string
// 	remove  string
// 	addcmd  string
// 	cmdDesc string
// 	disable string
// )

// func init() {
// 	flag.StringVar(&Token, "t", "", "discord bot token")
// 	flag.StringVar(&remove, "r", "", "remove command you want")
// 	flag.StringVar(&disable, "d", "", "disable all comands")
// 	flag.StringVar(&addcmd, "e", "", `add new command to system`)
// 	flag.StringVar(&cmdDesc, "cd", "", "add description for that you newly created command")
// 	flag.Parse()

// 	if Token == "" {
// 		flag.Usage()
// 		return
// 	}

// }

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSEGV)

	go func() {
		for {
			s := <-c
			handler.Handler(s)
		}
	}()

	err := playwright.Install(&playwright.RunOptions{
		Verbose:  true,
		Browsers: []string{"chromium"},
	})
	if err != nil {
		log.Fatalf("Couldn't install browser %v", err)
	}

	pw, err := playwright.Run(&playwright.RunOptions{
		Browsers: []string{"chromium"},
		Verbose:  true,
	})

	if err != nil {
		log.Fatalf("Couldn't Run playwight %v", err)
	}

	browser, page := browsercontext.LaunchBrowser(pw, browsercontext.AuthTwitter, true)
	for !browsercontext.IsLoggedin(page) {
		cookies := <-browsercontext.Login(pw)

		for _, cookie := range cookies {
			if err := browser.AddCookies(playwright.BrowserContextAddCookiesOptionsCookies{
				Name:     &cookie.Name,
				Domain:   &cookie.Domain,
				Value:    &cookie.Value,
				Path:     &cookie.Path,
				Expires:  &cookie.Expires,
				HttpOnly: &cookie.HttpOnly,
				Secure:   &cookie.Secure,
				SameSite: &cookie.SameSite,
			}); err != nil {
				log.Fatalf("Couldn't set cookie to current context %v", err)
			}

		}

		_, err = page.Goto("https://chat.openai.com/chat")
		if err != nil {
			log.Fatalf("Couldnt go to %v err: %v", browsercontext.Twitter, err)
		}

		fmt.Printf("current path is: %v", page.URL())
		image, err := page.Screenshot()
		if err != nil {
			log.Fatalf("Couldt take screenshot form twitter %v", err)
		}
		ss, err := os.Create("ss.png")
		if err != nil {
			log.Fatalf("Couldn't create screenshot %v", err)
		}
		fw := bufio.NewWriter(ss)
		fw.Write(image)
		fmt.Println("all done !!")

	}

	if err := browser.ClearCookies(); err != nil {
		log.Printf("Couldn't clear out the cookies %v", err)
	}

	if err := browser.Close(); err != nil {
		log.Printf("Couldn't Shutdown the browser %v", err)
	}

	if err := pw.Stop(); err != nil {
		log.Printf("pw Couldn't stopped %v", err)
	}

	// session, err := discordgo.New("Bot " + Token)
	// if err != nil {
	// 	log.Fatalln("Couldn't start new session"
	// }

	// err = session.Open()
	// if err != nil {
	// 	log.Fatalln("Session couldn't opened")
	// 	return
	// }

	// fmt.Println("Discord Bot running ctrl + c for kill")

	// session.AddHandler(messageHandler)
	// session.AddHandler(followupHandler)

	// registeredCmds := make(map[string]*discordgo.ApplicationCommand, len(commands))

	// for _, v := range commands {
	// 	cmd, err := session.ApplicationCommandCreate(session.State.User.ID, "", v)
	// 	if err != nil {
	// 		log.Fatalf("Couldn't recognized current command %v", err)
	// 		return
	// 	}
	// 	registeredCmds[v.Name] = cmd
	// }
	// removeCommand(session, registeredCmds)
	// defer session.Close()

	// sc := make(chan os.Signal, 1)
	// signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	// <-sc

}

// func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	if m.Author.ID == s.State.User.ID {
// 		return
// 	}

// 	msgs := strings.Fields(m.Content)
// 	if len(msgs) <= 1 {
// 		fmt.Println("must be bigger than one")
// 		return
// 	}

// 	for i, item := range msgs {
// 		if strings.Contains(item, "https") {
// 			// this is for eamil; handle for donwlaod image;r
// 			avatar, err := getImage(item)
// 			if err != nil {
// 				log.Fatalf("an error accured try again later %v", err)
// 			}

// 			name := msgs[len(msgs)-1-i]

// 			user, err := s.UserUpdate(name, avatar)
// 			if err != nil {
// 				log.Fatalf("Couldn't update user %v", err)
// 			}

// 			fmt.Printf("The user %s", user)

// 		}

// 	}

// }

// func getImage(url string) (string, error) {
// 	var base64img string
// 	var contentType string

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return "", err
// 	}

// 	defer func() {
// 		_ = resp.Body.Close()
// 	}()

// 	img, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", err
// 	}

// 	contentType = http.DetectContentType(img)
// 	base64img = base64.StdEncoding.EncodeToString(img)

// 	return fmt.Sprintf("data:%s;base64,%s", contentType, base64img), nil

// }

// var (
// 	commands = []*discordgo.ApplicationCommand{{
// 		Name:        "followups",
// 		Description: "Followups messages",
// 	}}

// 	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
// 		"maker": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 			fmt.Println("make malive")
// 		},
// 	}
// )

// func followupHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

// 	if hook, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
// 		hook(s, i)
// 	}

// }

// func removeCommand(s *discordgo.Session, appCmds map[string]*discordgo.ApplicationCommand) {

// 	if remove == "" {
// 		return
// 	}

// 	if cmd, ok := appCmds[remove]; ok {
// 		err := s.ApplicationCommandDelete(s.State.User.ID, "", cmd.ID)
// 		if err != nil {
// 			fmt.Printf("Couldn't Remove app cmd that id %v", err)
// 			return
// 		}

// 		fmt.Printf("Successfully command removed %v", cmd.Name)
// 	}

// }

// func disableCommands(s *discordgo.Session, appCmds map[string]*discordgo.ApplicationCommand) {
// 	if disable == "" {
// 		log.Fatalf("Must be different than empty string")
// 		return
// 	}
// 	test := "followups"
// 	if cmd, ok := appCmds[test]; ok {
// 		fmt.Printf("You're Valid to go %v", cmd.ApplicationID)
// 	}
// }

// func blocker(s *discordgo.Session, m *discordgo.MessageCreate) bool {
// 	if s.State.User.ID == m.Author.ID {
// 		return false
// 	}
// 	return true
// }
