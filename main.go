package main

import (
	"fmt"
	"math/rand"
)

type Deck struct {
	cards []int
}

type Box [][]int

func main() {
	fmt.Println("Hello, world!")
	d := NewDeck()
	fmt.Println(d)
	d.Shuffle()
	fmt.Println(d)

	for range 52 {
		num, ok := d.Deal()
		fmt.Println(num, ok)
	}

	fmt.Println(d.cards)
}

func NewDeck() Deck {
	d := Deck{}
	d.cards = make([]int, 52, 52)
	count := 0
	for i := range 13 {
		for range 4 {
			d.cards[count] = i + 2
			count++
		}
	}

	return d
}

func (d *Deck) Shuffle() {
	for i := range d.cards {
		j := rand.Intn(i + 1)
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	}
}

func (d *Deck) Deal() (int, bool) {
	if len(d.cards) == 0 {
		return 0, false
	}
	deadNum := d.cards[0]
	d.cards = d.cards[1:]

	return deadNum, true
}
