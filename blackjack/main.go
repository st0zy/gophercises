package main

import (
	"fmt"
	"strings"

	"github.com/st0zy/gophercises/deck"
)

type Hand []deck.Card

func (h Hand) String() string {
	strs := make([]string, len(h))

	for i := range h {
		strs[i] = h[i].String()
	}
	return strings.Join(strs, ", ")
}

func (h Hand) dealerString() string {
	strs := make([]string, len(h))
	strs[0] = h[0].String()
	for i := 1; i < len(h); i++ {
		strs[i] = "**HIDDEN**"
	}
	return strings.Join(strs, ", ")
}

func (h Hand) MinScore() int {
	var score int
	for _, c := range h {
		score += int(c.Rank)
	}
	return score
}

func (h Hand) Score() int {
	var score = h.MinScore()
	for _, c := range h {
		if c.Rank == deck.Ace {
			// Original value of ace is 1, we can promote it to 11, In such cases
			// we add the delta to the minScore
			return score + 10
		}
	}
	return score
}

func Shuffle(gs GameState) GameState {
	ret := clone(gs)
	ret.Deck = deck.New(deck.Deck(3), deck.Shuffle)
	return ret
}

func Deal(gs GameState) GameState {
	ret := clone(gs)
	ret.Player = make(Hand, 0, 5)
	ret.Dealer = make(Hand, 0, 5)
	var card deck.Card
	for range 2 {
		for _, h := range []*Hand{&ret.Player, &ret.Dealer} {
			card, ret.Deck = draw(ret.Deck)
			*h = append(*h, card)
		}
	}
	ret.State = StatePlayerTurn
	return ret
}

func Hit(gs GameState) GameState {
	var ret = clone(gs)
	hand := gs.CurrentPlayer()
	var card deck.Card
	card, ret.Deck = draw(ret.Deck)
	*hand = append(*hand, card)
	if hand.Score() > 21 {
		return Stand(gs)
	}
	return ret
}

func Stand(gs GameState) GameState {
	var ret = clone(gs)
	// iota is used in the definition of ths states.
	ret.State++
	return ret
}

func EndHand(gs GameState) GameState {

	var ret = clone(gs)
	pScore, dealerScore := gs.Player.Score(), gs.Dealer.Score()

	fmt.Println("=== FINAL HANDS ===")
	fmt.Println("Player:", gs.Player.String(), gs.Player.MinScore(), pScore)
	fmt.Println("Dealer:", gs.Dealer.String(), gs.Dealer.MinScore(), dealerScore)

	switch {
	case pScore > 21:
		fmt.Println("You Busted!")
	case dealerScore > 21:
		fmt.Println("Dealer busted!")
	case pScore > dealerScore:
		fmt.Println("You Win!")
	case dealerScore > pScore:
		fmt.Println("Dealer wins!")
	case dealerScore == pScore:
		fmt.Println("Tie")
	}
	ret.Player = nil
	ret.Dealer = nil
	fmt.Println()
	return ret
}

func main() {
	var gs GameState
	gs = Shuffle(gs)
	for range 10 {
		gs = Deal(gs)

		var input string

		for gs.State == StatePlayerTurn {
			fmt.Println(gs.Player)
			fmt.Println(gs.Dealer.dealerString())
			fmt.Println("Do you want to (h)it or (s)tand?")
			fmt.Scanf("%s\n", &input)
			switch input {
			case "h":
				gs = Hit(gs)
			case "s":
				gs = Stand(gs)
			}
		}

		for gs.State == StateDealerTurn {
			if gs.Dealer.Score() <= 16 || (gs.Dealer.Score() == 17 && gs.Dealer.Score() != 17) {
				gs = Hit(gs)
			} else {
				gs = Stand(gs)
			}
		}

		gs = EndHand(gs)
	}
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

type State int

const (
	StatePlayerTurn State = iota
	StateDealerTurn
	StateHandOver
)

type GameState struct {
	Deck   []deck.Card
	Player Hand
	Dealer Hand
	State  State
}

func (gs *GameState) CurrentPlayer() *Hand {
	switch gs.State {
	case StatePlayerTurn:
		return &gs.Player
	case StateDealerTurn:
		return &gs.Dealer
	default:
		panic("Isn't a current players turn")
	}
}

func clone(gs GameState) GameState {
	ret := GameState{
		State:  gs.State,
		Deck:   make([]deck.Card, len(gs.Deck)),
		Player: make(Hand, len(gs.Player)),
		Dealer: make(Hand, len(gs.Dealer)),
	}

	copy(ret.Deck, gs.Deck)
	copy(ret.Player, gs.Player)
	copy(ret.Dealer, gs.Dealer)
	return ret
}
