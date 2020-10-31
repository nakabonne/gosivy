package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nakabonne/gosivy/agent"
)

func main() {
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal(err)
	}
	defer agent.Close()
	fmt.Println("Pause this process for one hour...")
	fmt.Println("Press Ctrl-C to quit.")
	time.Sleep(time.Hour)
}
