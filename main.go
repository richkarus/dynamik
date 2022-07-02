package main

import (
	"fmt"
)

func main() {
	// Check to ensure IP matches dynamic record in Cloudflare
	// if record doesn't match, update with current public IP

	if !CheckIpMatches() {
		fmt.Println("► ✨ Dynamic record updated. ✨")
	} else {
		fmt.Println("► IPs match. Nothing to do! 😎")
	}
}
