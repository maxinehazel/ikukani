package main

import (
	"fmt"
	"log"
	"os"

	"git.maxinekrebs.dev/softpunk/ikukani"
)

func main() {
	ikukani.Token = os.Getenv("WK_TOKEN")
	resp, err := ikukani.Summary()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
}
