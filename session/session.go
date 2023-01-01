package session

import (
	"log"

	"github.com/playwright-community/playwright-go"
)

var loginCookieName string = "__Secure-next-auth.session-token"

func AddToSession(pw playwright.BrowserContext, url string, cookieValue string) playwright.BrowserContext {

	cookies, err := pw.Cookies(url)
	if err != nil {
		log.Fatalf("Couldn't get cookie from url %v", err)
	}

	for _, cookie := range cookies {

		convCookie := playwright.BrowserContextAddCookiesOptionsCookies{
			Name:     &cookie.Name,
			Value:    &cookie.Value,
			Domain:   &cookie.Domain,
			HttpOnly: &cookie.HttpOnly,
			SameSite: &cookie.SameSite,
			Path:     &cookie.Path,
			Expires:  &cookie.Expires,
		}

		if err := pw.AddCookies(convCookie); err != nil {
			log.Fatalf("Couldn't set cookies %v", err)
		}

		// final cookie and most important one ;

		if err := pw.AddCookies(playwright.BrowserContextAddCookiesOptionsCookies{
			Name:     &loginCookieName,
			Value:    &cookieValue,
			Domain:   &cookie.Domain,
			Path:     &cookie.Path,
			SameSite: playwright.SameSiteAttributeLax,
			Secure:   playwright.Bool(true),
			Expires:  &cookie.Expires,
		}); err != nil {
			log.Fatalf("Could set main session cookie %v", err)
		}

	}

	return pw

}
