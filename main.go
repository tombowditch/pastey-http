package main

import (
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
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
	err := godotenv.Load()
	if err != nil {
		panic("Cannot load .env file, please make sure it is created.")
	}

	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	logrus.Info("starting...")

	r := httprouter.New()

	r.GET("/", indexPage)
	r.GET("/:identifier", getIdentifier)

	if err := http.ListenAndServe(os.Getenv("BIND_HOST")+":"+os.Getenv("BIND_PORT"), r); err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}

}
