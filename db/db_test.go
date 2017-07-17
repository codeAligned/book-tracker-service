package book_tracker


import (
	. "github.com/onsi/gomega"
	"testing"
	"fmt"
)


type SqliteConfigProvider struct {

}

func (p *SqliteConfigProvider) GetConfig() (Config, error) {
	config := &Config{
		filename: "/Users/gigi.sayfan/git/book-tracker/book-tracker.db",
	}

	return *config, nil
}


func TestGetBooks(t *testing.T) {
	RegisterTestingT(t)

	db := NewDB(&SqliteConfigProvider{})
	books, err := db.GetBooks()
	Ω(err).Should(BeNil())
	Ω(books).Should(HaveLen(2))
}

func TestGetCategories(t *testing.T) {
	RegisterTestingT(t)

	db := NewDB(&SqliteConfigProvider{})
	categories, err := db.GetCategories()
	Ω(err).Should(BeNil())
	Ω(categories).Should(HaveLen(6))
}

func TestGetSalesRanks	(t *testing.T) {
	RegisterTestingT(t)

	db := NewDB(&SqliteConfigProvider{})
	ranks, err := db.GetSalesRanks("", "", nil, nil)
	Ω(err).Should(BeNil())
	for _, r := range ranks {
		fmt.Println(r)

	}
}
