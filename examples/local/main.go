package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nakabonne/gosivy/agent"
)

func main() {
	err := agent.Listen(agent.Options{
		LogWriter: os.Stdout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer agent.Close()

	fmt.Println("Press Ctrl-C to quit.")
	time.Sleep(time.Hour)
}
