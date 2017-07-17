package book_tracker

import "time"

type Book struct {
	Name string
	ISBN string
	URL string
}


type Category struct {
	Name string
}


type SalesRank struct {
	Book *Book
	Category *Category
	Rank int
	Timstamp time.Time
	Change int
}

type PotentialSale struct {
	Book *Book
	Timstamp time.Time
	AverageChange int
}
