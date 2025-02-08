package blackjack

import (
	"fmt"
	"strconv"

	"github.com/st0zy/gophercises/deck"
)

type AI interface {
	Play(hand []deck.Card, dealer deck.Card) Move
	Results(hand []deck.Card, dealer []deck.Card)
	Bet(shuffled bool) int
}

type humanAI struct{}

func HumanAI() AI {
	return humanAI{}
}

func (ai humanAI) Play(hand []deck.Card, dealer deck.Card) Move {
	fmt.Println("Player::", hand)
	fmt.Println("Dealer::", dealer)
	fmt.Println("Do you want to (h)it or (s)tand or (d)ouble?")
	for {
		var input string
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			return MoveHit
		case "s":
			return MoveStand
		case "d":
			return MoveDouble
		default:
			fmt.Println("Invalid option")
		}
	}
}

func (ai humanAI) Results(hand []deck.Card, dealer []deck.Card) {
	fmt.Println("=== FINAL HANDS ===")
	fmt.Println("Player:", hand)
	fmt.Println("Dealer:", dealer)

}

func (ai humanAI) Bet(shuffled bool) int {
	if shuffled {
		fmt.Println("The deck was just shuffled")
	}
	fmt.Println("What would you like to Bet?")
	var input string
	fmt.Scanf("%s\n", &input)
	ret, _ := strconv.Atoi(input)
	return ret
}

type DealerAI struct{}

func (ai DealerAI) Bet(shuffled bool) int {
	return -1
}

func (ai DealerAI) Play(hand []deck.Card, dealer deck.Card) Move {
	dScore := Score(hand...)
	// dealer's rules of when to hit or stand
	if dScore <= 16 || (dScore == 17 && Soft(hand...)) {
		return MoveHit
	} else {
		return MoveStand
	}
}

func (ai DealerAI) Results(hand []deck.Card, dealer []deck.Card) {
}
