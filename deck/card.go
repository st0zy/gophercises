//go:generate stringer -type=Suit,Rank

package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Suit uint8

const (
	Spade Suit = iota
	Diamond
	Club
	Heart
	Joker // this is a special case
)

var suits = [...]Suit{Spade, Diamond, Club, Heart}

type Rank uint8

const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

const (
	minRank = Ace
	maxRank = King
)

type Card struct {
	Suit
	Rank
}

func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}
	return fmt.Sprintf("%s of %ss", c.Rank, c.Suit)

}

// Reason for returning slice instead of modifying inline is to support adding new cards
type CardOptions func([]Card) []Card

func New(opts ...CardOptions) []Card {
	var cards []Card
	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			cards = append(cards, Card{Suit: suit, Rank: rank})
		}
	}
	for _, opt := range opts {
		cards = opt(cards)
	}
	return cards
}

func absoluteWeight(c Card) int {
	return int(c.Suit)*int(maxRank) + int(c.Rank)
}

func Less(cards []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return absoluteWeight(cards[i]) < absoluteWeight(cards[j])
	}
}

func Sort(less func(cards []Card) func(i, j int) bool) CardOptions {
	return func(cards []Card) []Card {
		sort.Slice(cards, less(cards))
		return cards
	}
}

func DefaultSort(cards []Card) []Card {
	sort.Slice(cards, Less(cards))
	return cards
}

func Shuffle(cards []Card) []Card {
	ret := make([]Card, len(cards))
	r := rand.New(rand.NewSource(time.Now().Unix()))
	// perm = [0,1,6,4,5,...]
	perm := r.Perm(len(cards))
	for i, j := range perm {
		ret[i] = cards[j]
	}
	return ret
}

func Jokers(n int) CardOptions {
	return func(cards []Card) []Card {
		for i := range n {
			cards = append(cards, Card{Rank: Rank(i), Suit: Joker})
		}
		return cards
	}
}

func Filter(f func(c Card) bool) CardOptions {
	return func(cards []Card) []Card {
		var ret []Card
		for _, c := range cards {
			if !f(c) {
				ret = append(ret, c)
			}
		}

		return ret
	}
}

func Deck(n int) CardOptions {
	return func(cards []Card) []Card {
		var ret []Card
		for i := 0; i < n; i++ {
			//Shallow copy only, but should work for this case
			ret = append(ret, cards...)
		}
		return ret
	}
}
