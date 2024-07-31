package box

import (
	"github.com/ccadden/beat-the-box/constants"
	"github.com/ccadden/beat-the-box/deck"
	"github.com/ccadden/beat-the-box/helpers"
	"math/rand"
)

type Box struct {
	cards []int
}

func NewBox() *Box {
	b := Box{}
	b.cards = []int{}
	return &b
}

// Populate a Box from a Deck
func (b *Box) PopulateFrom(d *deck.Deck) {
	for range constants.BOX_SIZE {
		num, ok := d.Deal()

		if !ok {
			panic("Couldn't deal card")
		}

		b.Add(num)
	}
}

func (b *Box) Add(card int) {
	b.cards = append(b.cards, card)
}

func (b *Box) Replace(idx, value int) {
	b.cards[idx] = value
}

func (b *Box) Shrink(idx int) {
	b.cards[idx] = b.cards[len(b.cards)-1]
	b.cards = b.cards[:len(b.cards)-1]
}

// returns the index and value of the card which has the highest or lowest value
func (b *Box) MostExtremeCard() (int, int) {
	var extremeIdx int
	var extreme int

	for idx, val := range b.cards {
		if idx == 0 {
			extremeIdx = idx
			extreme = val
		} else {
			if val == 2 || val == 14 { // pErFoRmAnCe
				return idx, val
			}
			// furtherst distance from midpoint would be most extreme
			if helpers.Abs(val-constants.MIDPOINT) > extreme {
				extreme = val
				extremeIdx = idx
			}
		}
	}

	return extremeIdx, extreme
}

func (b *Box) Empty() bool {
	return len(b.cards) == 0
}

// returns idx and value of random card in box
func (b *Box) RandomCard() (int, int) {
	idx := rand.Intn(len(b.cards))

	return idx, b.cards[idx]
}
