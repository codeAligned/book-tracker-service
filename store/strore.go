package book_tracker

import (
	"github.com/jmoiron/sqlx"
	"fmt"
	"time"
	. "github.com/the-gigi/book-tracker-service/models"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
)

type Config struct {
	host string
	port int
	user string
	dbname string
    password string
	sslmode string
	filename string
}

type ConfigProvider interface {
	GetConfig() (Config, error)
}


type DB struct {
	configProvider     	 ConfigProvider
	conn               	 *sqlx.DB
	books_by_name      	 map[string]*Book
	categories_by_name 	 map[string]*Category
	books_by_id      	 map[int]*Book
	book_name_id_map     map[string]int
	category_name_id_map map[string]int
	categories_by_id 	 map[int]*Category

}

func NewDB(configProvider ConfigProvider) *DB {
	return &DB{
		configProvider: configProvider,
	}
}


func makePostgressSqlConnectionString(conf Config) (string, string) {
	t := "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"
	return "postgress", fmt.Sprintf(t, conf.host, conf.port, conf.user, conf.password, conf.dbname, conf.sslmode)
}

func makeSqliteConnectionString(conf Config) (string, string) {
	return "sqlite3", conf.filename
}

func getConnectionInfo(conf Config) (driver string, connectionString string, err error) {
	err = nil
	if conf.filename != "" {
		driver, connectionString = makeSqliteConnectionString(conf)
	} else {
		driver, connectionString = makePostgressSqlConnectionString(conf)
	}

	return
}


func (db *DB) init() error {
	err := db.connect()
	if err != nil {
		return err
	}

	_, err = db.GetBooks()
	if err != nil {
		return err
	}

	_, err = db.GetCategories()
	if err != nil {
		return err
	}

	return nil

}

func (db *DB) connect() error {
	conf, err := db.configProvider.GetConfig()
	driver, connectionString, err := getConnectionInfo(conf)
	if err != nil {
		return err
	}
	db.conn, err = sqlx.Open(driver, connectionString)
	if err != nil {
		return err
	}

	err = db.conn.Ping()
	if err != nil {
		db.conn = nil
		return err
	}

	return nil
}


// Implement the BookTracker interface
type dbBook struct {
	ID int `db:"id"`
	Name string `db:"name"`
	ISBN string `db:"isbn"`
	URL  string `db:"url"`
}

func (db *DB) GetBooks() ([]Book, error) {
	if db.conn == nil {
		err := db.connect()
		if err != nil {
			return nil, err
		}
	}

	dbBooks := []dbBook{}
	err := db.conn.Select(&dbBooks, "SELECT id, name, isbn, url FROM book")
	if err != nil {
		return nil, err
	}

	books := []Book{}
	db.books_by_id = map[int]*Book{}
	db.books_by_name = map[string]*Book{}
	db.book_name_id_map = map[string]int{}
	for _, b := range dbBooks {
		book := &Book{
			Name: b.Name,
			ISBN: b.ISBN,
			URL: b.URL,
		}
		db.books_by_id[b.ID] = book
		db.books_by_name[b.Name] = book
		db.book_name_id_map[b.Name] = b.ID
		books = append(books, *book)
	}

	return books, nil
}

type dbCategory struct {
	ID int `db:"id"`
	Name string `db:"name"`
}

func (db *DB) GetCategories() ([]Category, error) {
	if db.conn == nil {
		err := db.connect()
		if err != nil {
			return nil, err
		}
	}

	dbCategories := []dbCategory{}
	err := db.conn.Select(&dbCategories, "SELECT id, name FROM category")
	if err != nil {
		return nil, err
	}

	categories := []Category{}
	db.categories_by_id = map[int]*Category{}
	db.categories_by_name = map[string]*Category{}
	db.category_name_id_map = map[string]int{}
	for _, c := range dbCategories {
		category := &Category{
			Name: c.Name,
		}
		categories = append(categories, *category)
		db.categories_by_id[c.ID] = category
		db.categories_by_name[c.Name] = category
		db.category_name_id_map[c.Name] = c.ID
	}

	return categories, nil
}

func (db *DB) getBookByName(name string) (*Book, error) {
	book, ok := db.books_by_name[name]
	if ok {

		return book, nil
	}

	b := dbBook{}
	err := db.conn.Get(&b, "SELECT * FROM book WHERE NAME='?'", name)
	if err != nil {
		return nil, err
	}

	book = &Book{
		Name: b.Name,
		ISBN: b.ISBN,
		URL: b.URL,
	}
	db.books_by_id[b.ID] = book
	db.books_by_name[name] = book
 	
	return book, nil
}



func (db *DB) getCategoryByName(name string) (*Category, error) {
	category, ok := db.categories_by_name[name]
	if ok {
		return category, nil
	}

	c := dbCategory{}
	err := db.conn.Get(&c, "SELECT * FROM category WHERE name='?'", name)
	if err != nil {
		return nil, err
	}

	category = &Category{Name: c.Name,}
	db.categories_by_id[c.ID] = category
	db.categories_by_name[name] = category
	db.category_name_id_map[name] = c.ID

	return category, nil
}


type dbRank struct {
	ID int     			`db:"id"`
	BookID int     		`db:"book_id"`
	CategoryID int 		`db:"category_id"`
	Rank int       		`db:"rank"`
	Timestamp time.Time `db:"timestamp"`
	Change int          `db:"change"`
}

func (db *DB) GetSalesRanks(bookName string,
							categoryName string,
							start *time.Time,
							end *time.Time) ([]SalesRank, error) {
	if db.conn == nil {
		err := db.init()
		if err != nil {
			return nil, err
		}
	}

	dbRanks := []dbRank{}
	query := "SELECT * FROM rank"
	first := true
	if bookName != "" {
		query += " WHERE "
		first = false
		query += fmt.Sprintf("book_id = %d", db.book_name_id_map[bookName])
	}

	if categoryName != "" {
		if first {
			query += " WHERE "
			first = false
		} else {
			query += " AND "
		}
		query += fmt.Sprintf("category_id = %d", db.category_name_id_map[categoryName])
	}

	if start != nil {
		if first {
			query += " WHERE "
			first = false
		} else {
			query += " AND "
		}
		query += fmt.Sprintf("timestamp >= '%s'", start)
	}

	if end != nil {
		if first {
			query += " WHERE "
			first = false
		} else {
			query += " AND "
		}
		query += fmt.Sprintf("timestamp < '%s'", end)
	}

	err := db.conn.Select(&dbRanks, query)
	if err != nil {
		return nil, err
	}

	ranks := []SalesRank{}
	for _, r := range dbRanks {
		book := db.books_by_id[r.BookID]
		category := db.categories_by_id[r.CategoryID]
		ranks = append(ranks, SalesRank{
			Book: book,
			Category: category,
			Rank: r.Rank,
			Timstamp: r.Timestamp,
			Change: r.Change,
		})
	}

	return ranks, nil
}

func (db *DB) GetPotentialSales(bookName string, start, end *time.Time) ([]PotentialSale, error){
	return nil, nil
}
