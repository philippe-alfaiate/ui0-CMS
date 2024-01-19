package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"encoding/json"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/gofiber/fiber/v3"
)

type User struct {
	Id     int64
	Name   string
	Emails []string
}

func (u User) String() string {
	return fmt.Sprintf("User<%d %s %v>", u.Id, u.Name, u.Emails)
}

type Story struct {
	Id       int64
	Title    string
	AuthorId int64
	Author   *User `pg:"rel:has-one"`
}

func (s Story) String() string {
	return fmt.Sprintf("Story<%d %s %s>", s.Id, s.Title, s.Author)
}

func ExampleDB_Model(cfg pg.Options) {
	db := pg.Connect(&cfg)
	defer db.Close()

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	user1 := &User{
		Name:   "admin",
		Emails: []string{"admin1@admin", "admin2@admin"},
	}
	_, err = db.Model(user1).Insert()
	if err != nil {
		panic(err)
	}

	_, err = db.Model(&User{
		Name:   "root",
		Emails: []string{"root1@root", "root2@root"},
	}).Insert()
	if err != nil {
		panic(err)
	}

	story1 := &Story{
		Title:    "Cool story",
		AuthorId: user1.Id,
	}
	_, err = db.Model(story1).Insert()
	if err != nil {
		panic(err)
	}

	// Select user by primary key.
	user := &User{Id: user1.Id}
	err = db.Model(user).WherePK().Select()
	if err != nil {
		panic(err)
	}

	// Select all users.
	var users []User
	err = db.Model(&users).Select()
	if err != nil {
		panic(err)
	}

	// Select story and associated author in one query.
	story := new(Story)
	err = db.Model(story).
		Relation("Author").
		Where("story.id = ?", story1.Id).
		Select()
	if err != nil {
		panic(err)
	}

	fmt.Println(user)
	fmt.Println(users)
	fmt.Println(story)
	// Output: User<1 admin [admin1@admin admin2@admin]>
	// [User<1 admin [admin1@admin admin2@admin]> User<2 root [root1@root root2@root]>]
	// Story<1 Cool story User<1 admin [admin1@admin admin2@admin]>>
}

// createSchema creates database schema for User and Story models.
func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*User)(nil),
		(*Story)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			// Temp: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

type PathMessage struct {
	Path *string `json:"path"`
}
type HealthLinkMessage struct {
	Link *bool `json:"link"`
}

func pingOther(ip string, port string, link *bool, salt string) {
	var body []byte
	var msg PathMessage

	url := "http://" + ip + ":" + port + "/api/" + salt
	log.Println("url:", url)

	resp, err := http.Get(url)
	if err == nil {
		body, err = io.ReadAll(resp.Body)
	}

	if err == nil {
		err = json.Unmarshal(body, &msg)
	}

	if err == nil {
		log.Println(*msg.Path)
		if *msg.Path == salt {
			*link = true
		}
	}

	if err != nil {
		log.Println(err)
	}
}

func checkLink(link *bool) {
	ip := os.Getenv("ADMIN_IP")
	port := os.Getenv("ADMIN_PORT")
	port_number, err := strconv.Atoi(port)

	max := 2500 + 500
	min := 2500 - 500

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	if net.ParseIP(ip) == nil || err != nil || port_number > 65535 || port_number < 0 {
		log.Fatalln("ADMIN_IP or ADMIN_PORT number are wrong or not set in ENV")
	}

	for !*link {
		delay_ms := random.Intn(max-min) + min
		time.Sleep(time.Millisecond * time.Duration(delay_ms))
		log.Println("Link check... delay:", delay_ms)
		pingOther(ip, port, link, strconv.Itoa(rand.Intn(999999)))
	}
	log.Println("link up!")
}

func AnyPointer[A any](v A) *A {
	return &v
}

func getEnv(name string) (bool, string) {
	is_empty := false
	value := os.Getenv(name)
	if len(value) <= 0 {
		is_empty = true
	}
	return !is_empty, value
}

func getPgOption() (bool, pg.Options) {
	missing := false
	var o pg.Options
	var err_b bool
	err_b, o.Database = getEnv("DB_NAME")
	missing = missing || !err_b
	err_b, o.User = getEnv("DB_USER")
	missing = missing || !err_b
	err_b, o.Password = getEnv("DB_PASSWORD")
	missing = missing || !err_b
	err_b, o.Addr = getEnv("DB_ADDR")
	missing = missing || !err_b

	return !missing, o
}

func main() {
	db_enable, db_cfg := getPgOption()
	if db_enable {
		log.Println("DB Data found", db_cfg)
		ExampleDB_Model(db_cfg)
	} else {
		log.Println("No DB data")
	}

	host := os.Getenv("ADMIN_LISTEN_HOST")
	if len(host) <= 0 {
		host = "127.0.0.1:3000"
	}

	app := fiber.New()
	app.Use(func(c fiber.Ctx) error {
		log.Println("BaseURL", c.BaseURL(), " IP:", c.IP())
		return c.Next()
	})

	link := false

	// GET /api/*
	app.Get("/api/*", func(c fiber.Ctx) error {

		var p PathMessage
		p.Path = AnyPointer(c.Params("*"))
		msg, _ := json.Marshal(p)
		c.Set("content-type", "application/json; charset=utf-8")
		return c.SendString(string(msg))
	})

	// GET health-link
	app.Get("/health-link", func(c fiber.Ctx) error {

		var h HealthLinkMessage
		h.Link = AnyPointer(link)
		msg, _ := json.Marshal(h)
		c.Set("content-type", "application/json; charset=utf-8")
		return c.SendString(string(msg))
	})

	go checkLink(&link)

	log.Fatal(app.Listen(host))
}
