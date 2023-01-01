package handler

import (
	"io/ioutil"
	"log"
)

func WriteToFile(data []byte) {
	err := ioutil.WriteFile("data.json", data, 0)
	if err != nil {
		log.Fatalf("Couldn't write content to file %v", err)
	}

}
