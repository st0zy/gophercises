package deck

import (
	"fmt"
	"testing"
)

func ExampleCard() {
	fmt.Println(Card{Suit: Heart, Rank: Ace})
	fmt.Println(Card{Suit: Spade, Rank: Jack})
	fmt.Println(Card{Suit: Diamond, Rank: Queen})
	fmt.Println(Card{Suit: Club, Rank: Two})
	fmt.Println(Card{Suit: Heart, Rank: Three})
	fmt.Println(Card{Suit: Heart, Rank: Four})
	fmt.Println(Card{Suit: Joker})

	//Output:
	//Ace of Hearts
	//Jack of Spades
	//Queen of Diamonds
	//Two of Clubs
	//Three of Hearts
	//Four of Hearts
	//Joker

}

func TestNew(t *testing.T) {
	cards := New()
	if len(cards) != 13*4 {
		t.Error("incorrect number of cards in a new deck")
	}
}

func TestDefaultSort(t *testing.T) {
	cards := New(DefaultSort)
	want := Card{Suit: Spade, Rank: Ace}
	if cards[0] != want {
		t.Error("expected first card to be Ace of Spades")
	}
}

func TestSort(t *testing.T) {
	cards := New(Sort(Less))
	want := Card{Suit: Spade, Rank: Ace}
	if cards[0] != want {
		t.Error("expected first card to be Ace of Spades")
	}
}

func TestJokers(t *testing.T) {
	cards := New(Jokers(3))
	count := 0
	for _, c := range cards {
		if c.Suit == Joker {
			count++
		}
	}

	if count != 3 {
		t.Error("jokers are not equal to 3")
	}
}

func TestFilter(t *testing.T) {
	filter := func(c Card) bool {
		return c.Rank == 2 || c.Rank == 3
	}
	cards := New(Filter(filter))

	for _, c := range cards {
		if c.Rank == 2 || c.Rank == 3 {
			t.Error("found cards which were to be removed.")
		}
	}
}

func TestDeck(t *testing.T) {
	cards := New(Deck(3))
	if len(cards) != 13*4*3 {
		t.Error("incorrect number of cards in a new deck")
	}
}
