package book_tracker


import (
	. "github.com/onsi/gomega"
	"testing"
	"fmt"
	"os"
)


type SqliteConfigProvider struct {
}

func (p *SqliteConfigProvider) GetConfig() (Config, error) {
	config := &Config{
		//filename: "/Users/gigi.sayfan/git/book-tracker/book-tracker.db",
		filename: "test.db",
	}

	return *config, nil
}


func resetDB() *DB {
	p := &SqliteConfigProvider{}
	conf, _ := p.GetConfig()
	os.Remove(conf.filename)
	schema := `
		CREATE TABLE book (
			id INTEGER NOT NULL,
			name VARCHAR(256) NOT NULL,
			isbn VARCHAR(13) NOT NULL,
			url VARCHAR(1024) NOT NULL,
			track BOOLEAN,
			PRIMARY KEY (id),
			UNIQUE (name),
			CHECK (track IN (0, 1))
		);
		CREATE TABLE category (
			id INTEGER NOT NULL,
			name VARCHAR(1024),
			PRIMARY KEY (id)
		);
		CREATE TABLE rank (
			id INTEGER NOT NULL,
			book_id INTEGER NOT NULL,
			category_id INTEGER NOT NULL,
			rank INTEGER NOT NULL,
			timestamp DATETIME NOT NULL, change INTEGER,
			PRIMARY KEY (id),
			FOREIGN KEY(book_id) REFERENCES book (id),
			FOREIGN KEY(category_id) REFERENCES category (id)
		);

	`
	db := NewDB(&SqliteConfigProvider{})
	err := db.connect()
	Ω(err).Should(BeNil())

	c := db.conn
	_, err = c.Exec(schema)
	Ω(err).Should(BeNil())

	// Add some data
	data := `
		INSERT INTO book VALUES(1, 'Book 1','123','url-1',1);
		INSERT INTO book VALUES(2, 'Book 2','456','url-2',1);
		INSERT INTO category VALUES(1, 'Category 1');
		INSERT INTO category VALUES(2, 'Category 2');
		INSERT INTO category VALUES(3, 'Category 3');
		INSERT INTO rank VALUES(1,1,1,10,'2017-07-17 00:00:00.000000',2);
		INSERT INTO rank VALUES(2,1,2,20,'2017-07-17 00:00:00.000000',2);
		INSERT INTO rank VALUES(3,2,1,80,'2017-07-17 00:00:00.000000',5);
		INSERT INTO rank VALUES(4,2,3,90,'2017-07-17 00:00:00.000000',5);
		INSERT INTO rank VALUES(5,1,1,12,'2017-07-17 01:00:00.000000',2);
		INSERT INTO rank VALUES(6,1,2,22,'2017-07-17 01:00:00.000000',2);
		INSERT INTO rank VALUES(7,2,1,75,'2017-07-17 01:00:00.000000',-5);
		INSERT INTO rank VALUES(8,2,3,85,'2017-07-17 01:00:00.000000',-5);
		INSERT INTO rank VALUES(9,1,1,14,'2017-07-17 02:00:00.000000',2);
		INSERT INTO rank VALUES(10,1,2,24,'2017-07-17 02:00:00.000000',2);
		INSERT INTO rank VALUES(11,2,1,100,'2017-07-17 02:00:00.000000',25);
		INSERT INTO rank VALUES(12,2,3,100,'2017-07-17 02:00:00.000000',15);
	`
	_, err = c.Exec( data)
	Ω(err).Should(BeNil())

	return db
}


func TestGetBooks(t *testing.T) {
	RegisterTestingT(t)

	db := resetDB()

	books, err := db.GetBooks()
	Ω(err).Should(BeNil())
	Ω(books).Should(HaveLen(2))
}

func TestGetCategories(t *testing.T) {
	RegisterTestingT(t)

	db := NewDB(&SqliteConfigProvider{})
	categories, err := db.GetCategories()
	Ω(err).Should(BeNil())
	Ω(categories).Should(HaveLen(3))
}

func TestGetSalesRanks	(t *testing.T) {
	RegisterTestingT(t)

	db := NewDB(&SqliteConfigProvider{})
	ranks, err := db.GetSalesRanks("", "", nil, nil)
	Ω(err).Should(BeNil())
	Ω(ranks).Should(HaveLen(12))
	for _, r := range ranks {
		fmt.Println(r)

	}
}
