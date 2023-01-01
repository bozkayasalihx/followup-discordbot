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
		fmt.Println("succesfullly removed /")
		os.Exit(1)
	}

}
