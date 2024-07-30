package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type Strategy int

const (
	Random = iota
	Smort
)

var schemes []Strategy = []Strategy{Random, Smort}

var results map[Strategy]map[rune]int = map[Strategy]map[rune]int{}

type Deck struct {
	cards []int
}

type Box struct {
	cards []int
}

func main() {
	for _, strat := range schemes {
		results[strat] = map[rune]int{'W': 0, 'L': 0}

		var m sync.Mutex
		var wg sync.WaitGroup
		for range 100000 {
			wg.Add(1)
			go async(strat, &m, &wg)
		}

		wg.Wait()
		fmt.Println("wins:", results[strat]['W'])
		fmt.Println("losses:", results[strat]['L'])
	}
}

func async(strat Strategy, m *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	defer m.Unlock()

	var result bool

	m.Lock()

	switch strat {
	case Random:
		result = RandomBeatTheBox()
	}

	if result {
		results[strat]['W']++
	} else {
		results[strat]['L']++

	}
}

func RandomBeatTheBox() bool {
	d := NewDeck()
	b := NewBox()

	// populate the box
	for range 9 {
		num, ok := d.Deal()

		if !ok {
			panic("Couldn't deal card")
		}

		b.Add(num)
	}

	for !d.Empty() {
		if b.Empty() {
			return false
		}
		// pick a card in the box to "beat"
		idx, val := b.RandomCard()

		// deal a card from the deck
		newCard, _ := d.Deal()

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

func randBool() bool {
	return rand.Float32() < 0.5
}

func NewBox() *Box {
	b := Box{}
	b.cards = []int{}
	return &b
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

func (d *Deck) Shuffle() {
	for i := range d.cards {
		j := rand.Intn(i + 1)
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	}
}

func (d *Deck) Empty() bool {
	return len(d.cards) == 0
}

func (b *Box) Empty() bool {
	return len(b.cards) == 0
}

// returns idx and value of random card in box
func (b *Box) RandomCard() (int, int) {
	idx := rand.Intn(len(b.cards))

	return idx, b.cards[idx]
}

func (d *Deck) Deal() (int, bool) {
	if len(d.cards) == 0 {
		return 0, false
	}
	deadNum := d.cards[0]
	d.cards = d.cards[1:]

	return deadNum, true
}
