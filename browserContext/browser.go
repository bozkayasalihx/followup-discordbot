package browsercontext

import (
	"discord-bot/handler"
	"log"
	"sync"

	"github.com/playwright-community/playwright-go"
)

func IsLoggedin(page playwright.Page) bool {
	return page.URL() == "https://chat.openai.com/chat"
}

var chatgptAuth string = "https://chat.openapi.com/auth/login"
var homeChatgpt string = "https://chat.openapi.com/chat"

func Login(pw *playwright.Playwright) <-chan []*playwright.BrowserContextCookiesResult {
	var lock sync.Mutex
	r := make(chan []*playwright.BrowserContextCookiesResult)
	lock.Lock()
	go func() {
		defer close(r)
		defer lock.Unlock()
		browser, page := LaunchBrowser(pw, homeChatgpt, false)

		log.Println("Login with openapi account please!")

		page.On("framenavigated", func(frame playwright.Frame) {
			if frame.URL() != homeChatgpt {
				return
			}
			lock.Unlock()
		})

		lock.Lock()

		cookies, err := browser.Cookies(homeChatgpt)
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
