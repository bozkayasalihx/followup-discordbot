package handler

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

func Handler(signal os.Signal) {
	if signal == syscall.SIGTERM || signal == syscall.SIGINT {
		if err := os.RemoveAll("/tmp/twitter"); err != nil {
			log.Fatalf("Couldn't remove /tmp/twitter %v", err)
			os.Exit(0)
			return
		}

		os.Exit(1)
		fmt.Println("succesfullly removed /tmp/twitter")
	}

}
