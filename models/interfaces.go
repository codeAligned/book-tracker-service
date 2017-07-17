package book_tracker

import "time"


type  BookTracker interface {
	GetBooks() ([]Book, error)
	GetCategories() ([]Category, error)
	GetSalesRanks(bookName, categoryName string, start, end *time.Time) ([]SalesRank, error)
	GetPotentialSales(bookName string, start, end *time.Time) ([]PotentialSale, error)
}
