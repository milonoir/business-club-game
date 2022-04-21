package main

type player struct {
	cash    int
	stocks  map[string]int
	actions []card
}
