package blackjack

import (
	"errors"
	"fmt"

	"github.com/st0zy/gophercises/deck"
)

type Move func(*Game) error

var errBust = errors.New("Went bust...")

func MoveHit(gs *Game) error {
	hand := gs.currentHand()
	var card deck.Card
	card, gs.deck = draw(gs.deck)
	*hand = append(*hand, card)
	if Score(*hand...) > 21 {
		return errBust
	}
	return nil
}

func MoveStand(gs *Game) error {
	gs.state++
	return nil
}

func MoveDouble(gs *Game) error {
	if len(gs.player) != 2 {
		return errors.New("can't double down now")
	}
	gs.playerBet *= 2
	MoveHit(gs)
	MoveStand(gs)
	return nil
}

type state int

const (
	statePlayerTurn state = iota
	stateDealerTurn
	stateHandOver
)

type Game struct {
	//unexported fields
	decks           int
	hands           int
	blackjackPayout float64

	deck  []deck.Card
	state state

	player    []deck.Card
	balance   int
	playerBet int

	dealer   []deck.Card
	dealerAI AI
}

type Options struct {
	Decks           int
	Hands           int
	BlackjackPayout float64
}

func New(opt Options) Game {
	ret := Game{
		state:    statePlayerTurn,
		dealerAI: DealerAI{},
	}

	if opt.Decks == 0 {
		opt.Decks = 3
	}
	if opt.Hands == 0 {
		opt.Hands = 1
	}
	if opt.BlackjackPayout <= 0 {
		opt.BlackjackPayout = 1.5
	}
	ret.decks = opt.Decks
	ret.hands = opt.Hands
	ret.blackjackPayout = opt.BlackjackPayout

	return ret
}

func (gs *Game) currentHand() *[]deck.Card {
	switch gs.state {
	case statePlayerTurn:
		return &gs.player
	case stateDealerTurn:
		return &gs.dealer
	default:
		panic("Isn't a current players turn")
	}
}

func (g *Game) Play(ai AI) int {
	// creating a new shuffled deck
	g.deck = nil
	for range g.hands {
		shuffled := false
		if len(g.deck) < 52*g.decks/3 {
			g.deck = deck.New(deck.Deck(3), deck.Shuffle)
			shuffled = true
		}
		bet(g, ai, shuffled)
		deal(g)

		for g.state == statePlayerTurn {
			hand := make([]deck.Card, len(g.player))
			copy(hand, g.player)
			move := ai.Play(hand, g.dealer[0])
			err := move(g)
			switch err {
			case errBust:
				MoveStand(g)
			case nil:
				continue
			default:
				panic(err)
			}
		}

		for g.state == stateDealerTurn {
			dealerHand := make([]deck.Card, len(g.dealer))
			copy(dealerHand, g.dealer)
			move := g.dealerAI.Play(dealerHand, g.dealer[0])
			move(g)
		}
		endHand(g, ai)
	}
	return g.balance
}

func deal(g *Game) {
	g.player = make([]deck.Card, 0, 5)
	g.dealer = make([]deck.Card, 0, 5)
	var card deck.Card
	for range 2 {
		card, g.deck = draw(g.deck)
		g.player = append(g.player, card)
		card, g.deck = draw(g.deck)
		g.dealer = append(g.dealer, card)
	}
	g.state = statePlayerTurn
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

func bet(g *Game, ai AI, shuffled bool) {
	bet := ai.Bet(shuffled)
	g.playerBet = bet
}

func minScore(hand ...deck.Card) int {
	var score int
	for _, c := range hand {
		score += min(int(c.Rank), 10)
	}
	return score
}

func min(score1 int, score2 int) int {
	if score1 < score2 {
		return score1
	}
	return score2
}

func Soft(hand ...deck.Card) bool {
	minScore := minScore(hand...)
	score := Score(hand...)

	return minScore != score
}

func Score(hand ...deck.Card) int {
	var score = minScore(hand...)
	if score > 11 {
		return score
	}
	for _, c := range hand {
		if c.Rank == deck.Ace {
			// Original value of ace is 1, we can promote it to 11, In such cases
			// we add the delta to the minScore
			return score + 10
		}
	}
	return score
}

func Blackjack(hand ...deck.Card) bool {
	if len(hand) == 2 && Score(hand...) == 21 {
		return true
	}
	return false
}

func endHand(g *Game, ai AI) {
	pScore, dealerScore := Score(g.player...), Score(g.dealer...)
	playerBlackjack, dealerBlackjack := Blackjack(g.player...), Blackjack(g.dealer...)
	winnings := g.playerBet
	//long term keep track of game results
	switch {
	case dealerBlackjack:
		winnings = -winnings
	case playerBlackjack:
		winnings = int(float64(winnings) * (g.blackjackPayout))
	case pScore > 21:
		fmt.Println("You Busted!")
		winnings = winnings * -1
	case dealerScore > 21:
		fmt.Println("Dealer busted!")
	case pScore > dealerScore:
		fmt.Println("You Win!")
	case dealerScore > pScore:
		fmt.Println("Dealer wins!")
		winnings = winnings * -1
	case dealerScore == pScore:
		fmt.Println("Tie")
		winnings = 0
	}
	fmt.Println()
	ai.Results(g.player, g.dealer)
	g.player = nil
	g.dealer = nil
	g.balance += winnings
	fmt.Println()
}
