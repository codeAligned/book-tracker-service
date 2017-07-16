package book_tracker

import (
	"github.com/jmoiron/sqlx"
	"database/sql"
	"fmt"
	"time"
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
	configProvider ConfigProvider
	conn *sql.DB
}

func NewDB(configProvider ConfigProvider) *DB {
	return &DB{
		configProvider: configProvider,
	}
}


func makePostgressSqlConnectionString(conf Config) string {
	t := "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"
	return fmt.Sprintf(t, conf.host, conf.port, conf.user, conf.password, conf.dbname, conf.sslmode)
}

func makeSqliteConnectionString(conf Config) string {
	return conf.filename
}

func getConnectionString(conf Config) string {
	if conf.filename != "" {
		return makeSqliteConnectionString(conf)
	} else {
		return makePostgressSqlConnectionString(conf)
	}
}


func connect(db *DB) error {
	conf, err := db.configProvider.GetConfig()
	connectionString := getConnectionString(conf)
	db.conn, err = sql.Open("postgres", connectionString)
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

func (db *DB) GetBooks() ([]Book, error) {

}

func (db *DB) GetCategories ([]Category, error) {
}

func (db *DB) GetSalesRanks(bookName string,
							categoryName string,
							start time.Time,
							end time.Time) ([]SalesRank, error) {

}

func (db *DB) GetPotentialSales(bookName string, start, end time.Time) ([]PotentialSale, error){

}
