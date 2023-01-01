package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/playwright-community/playwright-go"
)

type DataShape struct {
	Name     *string `json:"name"`
	Domain   *string `json:"domain"`
	Path     *string `json:"path"`
	HttpOnly *bool   `json:"httponly"`
	Secure   *bool   `json:"secure"`
	Value    *string `json:"value"`
}

func WriteToFile(data []*playwright.BrowserContextCookiesResult) {
	for _, item := range data {
		s := DataShape{
			Value:    &item.Value,
			Name:     &item.Name,
			Secure:   &item.Secure,
			Domain:   &item.Domain,
			HttpOnly: &item.HttpOnly,
			Path:     &item.Path,
		}

		file, err := json.MarshalIndent(s, "", " ")
		if err != nil {
			log.Fatalf("Couldnt jsonify the data %v", err)
		}

		err = ioutil.WriteFile("sessin.json", file, 0644)
		if err != nil {
			log.Fatalf("Could save data to file %v", err)
		}

	}

}
