package main

import (
	"fmt"
	"math/rand"
	"sync"
)

const (
	Random = iota
	Smort
	Count
)

const MIDPOINT int = 8
const CARD_VIBES int = 10
const BOX_SIZE int = 9

var schemes map[string]Strategy = map[string]Strategy{"Random": Random, "Smort": Smort, "Count": Count}
var results map[Strategy]map[rune]int = map[Strategy]map[rune]int{}

type Strategy int

// Aces considered high always -- for now
type Deck struct {
	cards []int
}

type Box struct {
	cards []int
}

func randBool() bool {
	return rand.Float32() < 0.5
}

func getHotness(card int) int {
	switch {
	case card < 7:
		return 1
	case card > 9:
		return -1
	default:
		return 0
	}
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
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

func NewBox() *Box {
	b := Box{}
	b.cards = []int{}
	return &b
}

// Populate a Box from a Deck
func (b *Box) PopulateFrom(d *Deck) {
	for range BOX_SIZE {
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
			if Abs(val-MIDPOINT) > extreme {
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

func async(strat Strategy, m *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	defer m.Unlock()

	var result bool

	m.Lock()

	switch strat {
	case Random:
		result = RandomBeatTheBox()
	case Smort:
		result = SmortBeatTheBox()
	case Count:
		result = CardCountingBeatTheBox()
	}

	if result {
		results[strat]['W']++
	} else {
		results[strat]['L']++

	}
}

func CardCountingBeatTheBox() bool {
	hotness := 0

	d := NewDeck()
	b := NewBox()

	// populate the box
	b.PopulateFrom(d)

	for !d.Empty() {
		if b.Empty() {
			return false
		}

		idx, val := b.MostExtremeCard()

		// deal a card from the deck
		newCard, ok := d.Deal()

		if !ok {
			panic("Couldn't deal card")
		}

		if Abs(hotness) > CARD_VIBES { // hotness exceeds the vibes limit
			// Need to fix this block, not calculating whether it should be higher or lower correctly
			if hotness < 0 { // lots of high cards have been played, low card more likely
				if newCard < val {
					b.Replace(idx, newCard)
				} else {
					b.Shrink(idx)
				}
			} else { // lots of low cards have been played, high card more likely
				if newCard > val {
					b.Replace(idx, newCard)
				} else {
					b.Shrink(idx)
				}
			}
		} else { // use the default method
			if val < MIDPOINT {
				if newCard > val {
					b.Replace(idx, newCard)
				} else {
					b.Shrink(idx)
				}
			} else { // card is "high" guess lower
				if newCard < val {
					b.Replace(idx, newCard)
				} else {
					b.Shrink(idx)
				}
			}
		}

		hotness += getHotness(newCard)
	}

	return true
}

func SmortBeatTheBox() bool {
	d := NewDeck()
	b := NewBox()

	// populate the box
	b.PopulateFrom(d)

	for !d.Empty() {
		if b.Empty() {
			return false
		}

		idx, val := b.MostExtremeCard()

		// deal a card from the deck
		newCard, ok := d.Deal()

		if !ok {
			panic("Couldn't deal card")
		}

		// card is "low" guess higher
		if val < MIDPOINT {
			if newCard > val {
				b.Replace(idx, newCard)
			} else {
				b.Shrink(idx)
			}
		} else { // card is "high" guess lower
			if newCard < val {
				b.Replace(idx, newCard)
			} else {
				b.Shrink(idx)
			}
		}
	}

	return true
}

func RandomBeatTheBox() bool {
	d := NewDeck()
	b := NewBox()

	// populate the box
	b.PopulateFrom(d)

	for !d.Empty() {
		if b.Empty() {
			return false
		}
		// pick a card in the box to "beat"
		idx, val := b.RandomCard()

		// deal a card from the deck
		newCard, ok := d.Deal()

		if !ok {
			panic("Couldn't deal card")
		}

		// choose "strategy" - higher/lower
		if randBool() { // higher
			// check if successful, if not "close" that index of the box
			if newCard > val {
				b.Replace(idx, newCard)
			} else {
				b.Shrink(idx)
			}
		} else { // lower
			// check if successful, if not "close" that index of the box
			if newCard < val {
				b.Replace(idx, newCard)
			} else {
				b.Shrink(idx)
			}
		}

	}

	return true
}

func main() {
	for label, strat := range schemes {
		results[strat] = map[rune]int{'W': 0, 'L': 0}

		var m sync.Mutex
		var wg sync.WaitGroup
		for range 100000 {
			wg.Add(1)
			go async(strat, &m, &wg)
		}

		wg.Wait()
		fmt.Println(label)
		fmt.Println("wins:", results[strat]['W'])
		fmt.Println("losses:", results[strat]['L'])
		fmt.Printf("win %%: %0.2f\n\n", float64(results[strat]['W'])/float64(results[strat]['L'])*100)
	}
}
