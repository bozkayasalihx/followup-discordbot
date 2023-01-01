package browsercontext

import (
	"discord-bot/handler"
	"log"
	"sync"

	"github.com/playwright-community/playwright-go"
)

// const Twitter = "https://twitter.com/home"
const Twitter = "https://twitter.com/home"
const AuthTwitter = "https://chat.openai.com/auth/login"

func IsLoggedin(page playwright.Page) bool {
	return page.URL() == "https://chat.openai.com/chat"
}

func Login(pw *playwright.Playwright) <-chan []*playwright.BrowserContextCookiesResult {
	var lock sync.Mutex
	r := make(chan []*playwright.BrowserContextCookiesResult)
	lock.Lock()
	go func() {
		defer close(r)
		defer lock.Unlock()
		browser, page := LaunchBrowser(pw, AuthTwitter, false)

		log.Println("Login with twitter account please!")

		page.On("framenavigated", func(frame playwright.Frame) {
			if frame.URL() != "https://chat.openapi.com/chat" {
				return
			}
			lock.Unlock()
		})

		lock.Lock()

		cookies, err := browser.Cookies(Twitter)
		if err != nil {
			log.Fatalf("Couldn't get cookie from home %v", err)
		}

		// browser context  now writes content to file ;
		handler.WriteToFile(cookies)

		if err := browser.Close(); err != nil {
			log.Fatalf("Couldn't Close the browser %v", err)
		}

		r <- cookies

	}()
	return r
}

func LaunchBrowser(pw *playwright.Playwright, url string, headless bool) (playwright.BrowserContext, playwright.Page) {

	browser, err := pw.Chromium.LaunchPersistentContext("/tmp/twitter", playwright.BrowserTypeLaunchPersistentContextOptions{
		Headless: playwright.Bool(headless),
	})

	if err != nil {
		log.Fatalf("Couldnt start new instace of chromium %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		log.Printf("Couldnt start new page %v", err)
	}

	if _, err := page.Goto(url); err != nil {
		log.Fatalf("Couldn't go %v: %v", url, err)
	}

	return browser, page

}
