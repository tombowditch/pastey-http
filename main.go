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
	w.Write([]byte(`bind.sh - commandline pastebin`))
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

	if err := http.ListenAndServe("0.0.0.0:3333", r); err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}

}
