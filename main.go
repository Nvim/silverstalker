package main

import bot "github.com/Nvim/silverstalker/Bot"

import api "github.com/Nvim/silverstalker/Api"
import "fmt"

func main() {
	bot.BotToken = ""
	// bot.Run()
	fmt.Println(api.Api())
}
