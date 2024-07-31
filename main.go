package main

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/ccadden/beat-the-box/box"
	"github.com/ccadden/beat-the-box/constants"
	"github.com/ccadden/beat-the-box/deck"
	"github.com/ccadden/beat-the-box/helpers"
)

const (
	Random = iota
	Smort
	Count
)

var schemes map[string]Strategy = map[string]Strategy{"Random": Random, "Smort": Smort, "Count": Count}
var results map[Strategy]map[rune]int = map[Strategy]map[rune]int{}

type Strategy int

// Aces considered high always -- for now

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

	d := deck.NewDeck()
	b := box.NewBox()

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

		if helpers.Abs(hotness) > constants.CARD_VIBES { // hotness exceeds the vibes limit
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
			if val < constants.MIDPOINT {
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
	d := deck.NewDeck()
	b := box.NewBox()

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
		if val < constants.MIDPOINT {
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
	d := deck.NewDeck()
	b := box.NewBox()

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
