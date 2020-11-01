package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nakabonne/gosivy/agent"
)

func main() {
	err := agent.Listen(agent.Options{
		Addr: "127.0.0.1:9090",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer agent.Close()

	fmt.Println("Press Ctrl-C to quit.")
	time.Sleep(time.Hour)
}
