package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/richkarus/dynamik/internal"
	"os"
)

func main() {
	client, err := dynamik.NewClient()
	if err != nil {
		log.Fatal("unable to create a new Dynamik client", "fatal", err)
		os.Exit(1)
	}

	if !client.CheckIpMatches() {
		fmt.Println("► ✨ Dynamic record updated. ✨")
	} else {
		fmt.Println("► IPs match. Nothing to do! 😎")
	}
}
