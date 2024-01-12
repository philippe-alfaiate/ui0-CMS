package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"encoding/json"

	"github.com/gofiber/fiber/v3"
)

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
		body, err = ioutil.ReadAll(resp.Body)
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

func main() {
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
