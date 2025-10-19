package main

import (
	"fmt"
	"log"

	"github.com/snowlynxsoftware/parallax-game/cmd"
)

func main() {
	handler := cmd.NewHandler()
	err := handler.ExecuteCommand()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("done")
}
