package book_tracker

import "time"

type Book struct {
	Name string
	ISBN string
	URL string
}


type Category struct {
	name string
}


type SalesRank struct {
	book *Book
	category *Category
	rank int
	timstamp time.Time
	change int
}

type PotentialSale struct {
	book *Book
	timstamp time.Time
	averageChange int
}
