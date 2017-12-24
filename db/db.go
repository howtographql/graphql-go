package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/source/file"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
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
	PostedBy    string
}

type User struct {
	ID       string
	Name     string
	Email    string
	Password string
}

type AuthData struct {
	Email    string
	Password string
}

func CreateLink(link *Link) *Link {
	link.ID = uuid.NewV4().String()
	_, err := db.Exec("INSERT INTO link VALUES($1, $2, $3, $4)", link.ID, link.URL, link.Description, link.PostedBy)
	if err != nil {
		log.Println("error inserting link", err)
	}
	return link
}

// AllLinks takes a wrapper function that can be used
// to format each link in the result into the desired format
// for the caller
func AllLinks(wrapper func(link *Link)) {
	res, err := db.Query("SELECT id, url, description, posted_by FROM link")
	if err != nil {
		log.Println("error fetching all links", err)
	}
	defer res.Close()
	for res.Next() {
		link := Link{}
		if err := res.Scan(&link.ID, &link.URL, &link.Description, &link.PostedBy); err != nil {
			log.Fatal(err)
		}
		wrapper(&link)
	}
}

func FindUserByEmail(email string) *User {
	user := &User{}
	err := db.QueryRow("SELECT id, name, email, password FROM users WHERE email = $1", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		log.Println("error fetching user by email", err)
	}
	return user
}

func FindUserByID(id string) *User {
	user := &User{}
	err := db.QueryRow("SELECT id, name, email, password FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		log.Println("error fetching user by id", err)
	}
	return user
}

func CreateUser(user *User) *User {
	user.ID = uuid.NewV4().String()
	_, err := db.Exec("INSERT INTO users VALUES($1, $2, $3, $4)", user.ID, user.Name, user.Email, user.Password)
	if err != nil {
		log.Println("error inserting user", err)
	}
	return user
}
