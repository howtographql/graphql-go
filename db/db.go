package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/source/file"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	"github.com/satori/go.uuid"
	"time"
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
	if err != nil {
		log.Println("Error occurred during migration", err)
	}
	m.Steps(4)
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

type Vote struct {
	ID        string
	CreatedAt string
	UserID    string
	LinkID    string
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

func FindLinkByID(id string) *Link {
	link := &Link{}
	err := db.QueryRow("SELECT id, url, description, posted_by FROM link WHERE id = $1", id).Scan(&link.ID, &link.URL, &link.Description, &link.PostedBy)
	if err != nil {
		log.Println("error fetching link by id", err)
	}
	return link
}

func CreateVote(vote *Vote) *Vote {
	vote.ID = uuid.NewV4().String()
	vote.CreatedAt = time.Now().Format(time.RubyDate)
	_, err := db.Exec("INSERT INTO vote(id, created_at, user_id, link_id) VALUES($1, $2, $3, $4)", vote.ID, vote.CreatedAt, vote.UserID, vote.LinkID)
	if err != nil {
		log.Println("error inserting vote", err)
	}
	return vote
}

func FindVotesByLinkID(linkID string, wrapper func(vote *Vote)) {
	res, err := db.Query("SELECT id, created_at, user_id, link_id FROM vote WHERE link_id = $1", linkID)
	if err != nil {
		log.Println("error fetching votes", err)
	}
	defer res.Close()
	for res.Next() {
		vote := Vote{}
		if err := res.Scan(&vote.ID, &vote.CreatedAt, &vote.UserID, &vote.LinkID); err != nil {
			log.Fatal(err)
		}
		wrapper(&vote)
	}
}

func FindVotesByUserID(userID string, wrapper func(vote *Vote)) {
	res, err := db.Query("SELECT id, created_at, user_id, link_id FROM vote WHERE user_id = $1", userID)
	if err != nil {
		log.Println("error fetching votes", err)
	}
	defer res.Close()
	for res.Next() {
		vote := Vote{}
		if err := res.Scan(&vote.ID, &vote.CreatedAt, &vote.UserID, &vote.LinkID); err != nil {
			log.Fatal(err)
		}
		wrapper(&vote)
	}
}
