package db

import (
	"github.com/mattes/migrate"
	_ "github.com/lib/pq"
	"database/sql"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
	"log"
	"github.com/satori/go.uuid"
)

const dbURL = "postgres://postgres@localhost:5432/hackernews?sslmode=disable"

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"hackernews", driver)
	m.Steps(2)
}

type Link struct {
	ID          string
	URL         string
	Description string
}

func CreateLink(link *Link) *Link {
	link.ID = uuid.NewV4().String()
	_, err := db.Exec("INSERT INTO link VALUES($1, $2, $3)", link.ID, link.URL, link.Description)
	if err != nil {
		log.Println("error inserting link", err)
	}
	return link
}

// AllLinks takes a wrapper function that can be used
// to format each link in the result into the desired format
// for the caller
func AllLinks(wrapper func(link *Link)) {
	res, err := db.Query("SELECT id, url, description from link")
	if err != nil {
		log.Println("error fetching all links", err)
	}
	defer res.Close()
	for res.Next() {
		link := Link{}
		if err := res.Scan(&link.ID, &link.URL, &link.Description); err != nil {
			log.Fatal(err)
		}
		wrapper(&link)
	}
}
