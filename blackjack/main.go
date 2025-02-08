package main

import (
	"fmt"

	"github.com/st0zy/gophercises/blackjack/blackjack"
)

func main() {
	game := blackjack.New(blackjack.Options{})
	fmt.Println(game.Play(blackjack.HumanAI()))
}
