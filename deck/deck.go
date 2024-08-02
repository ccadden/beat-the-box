package deck

import "math/rand"

type Deck struct {
	cards []int
}

func NewDeck() *Deck {
	d := Deck{}
	d.cards = make([]int, 52, 52)
	count := 0
	for i := range 13 {
		for range 4 {
			d.cards[count] = i + 2
			count++
		}
	}

	d.Shuffle()

	return &d
}

func (d *Deck) Shuffle() {
	for i := range d.cards {
		j := rand.Intn(i + 1)
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	}
}

func (d *Deck) Empty() bool {
	return len(d.cards) == 0
}

func (d *Deck) Deal() (int, bool) {
	if len(d.cards) == 0 {
		return 0, false
	}
	deadNum := d.cards[0]
	d.cards = d.cards[1:]

	return deadNum, true
}

func (d *Deck) DealMultiple(numCards int) ([]int, bool) {
	if numCards > len(d.cards) {
		return []int{}, false
	}

	deadCards := d.cards[:numCards]
	d.cards = d.cards[numCards:]

	return deadCards, true
}

func (d *Deck) AddCard(card int) {
	d.cards = append(d.cards, card)
}

func (d *Deck) AddCards(cards []int) {
	for _, card := range cards {
		d.AddCard(card)
	}
}

func (d *Deck) CardsRemaining() int {
	return len(d.cards)
}
