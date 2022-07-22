package main

import (
	"homework-1/internal/commander"
	"homework-1/internal/handlers"
	"log"
	"os"
)

func main() {
	tgApiKey := os.Getenv("TG_API_KEY")
	if tgApiKey == "" {
		log.Fatal("TG_API_KEY env variable does not exist")
	}

	cmd, err := commander.Init(tgApiKey)
	if err != nil {
		log.Fatal(err)
	}

	handlers.AddHandlers(cmd)

	if err = cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
