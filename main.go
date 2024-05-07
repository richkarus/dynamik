package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"os"
)

func main() {
	client, err := NewDynamikClient()
	if err != nil {
		log.Fatal("unable to create Dynamik Client", "fatal", err)
		os.Exit(1)
	}

	if !client.CheckIpMatches() {
		fmt.Println("► ✨ Dynamic record updated. ✨")
	} else {
		fmt.Println("► IPs match. Nothing to do! 😎")
	}
}
