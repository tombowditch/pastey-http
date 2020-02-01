package main

import (
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func indexPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(`ig.lc - commandline pastebin
	
	pipe to 'nc ig.lc 3333'

	open ports: 1111, 2222, 3333, 4444, 5555, 6666, 7777, 8888, 9999

	- pastes are stored for 72 hours, after which they are automatically deleted
	- you can use ig.lc or bind.sh, both domains are setup for the paste service
	
	example
	=======

	~> echo "hello" | nc ig.lc 9999
	https://ig.lc/yourpaste

	~> cat /etc/nginx/nginx.conf | nc ig.lc 3333
	https://ig.lc/yourpaste

	~> cat 100mb.bin | nc ig.lc 9999
	too much data


	`))
}

func getIdentifier(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	identifier := ps.ByName("identifier")

	val, _ := client.Get("pastey_" + identifier).Result()

	if val != "" {
		// yea
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(val))
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("not found or expired"))
	}

}

var client *redis.Client

func main() {
	client = redis.NewClient(&redis.Options{
		Addr:     "pastey-redis:6379",
		Password: "",
		DB:       0,
	})

	logrus.Info("starting...")

	r := httprouter.New()

	r.GET("/", indexPage)
	r.GET("/:identifier", getIdentifier)

	if err := http.ListenAndServe("0.0.0.0:3334", r); err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}

}
