package main

type player struct {
	Cash    int
	Stocks  map[string]int
	Actions []card
}
