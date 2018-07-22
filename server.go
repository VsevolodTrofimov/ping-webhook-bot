package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"time"
)

func makeMessageHandler(db *gorm.DB, c chan Ping) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var ping Ping
		var proj Project
		q := r.URL.Query()
		uuid := mux.Vars(r)["uuid"]

		pingVal := q["m"]
		if len(pingVal) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Missing Ping parameter use /:id?m=\"your message\"")
			return
		}
		ping.Val = pingVal[0]

		err := db.First(&proj, "uuid = ?", uuid).RecordNotFound()
		fmt.Println(proj, &proj)
		if err {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Proejct with id %s not found, try to use /list command on the bot", uuid)
			return
		}
		ping.Proj = proj

		if len(q["t"]) == 0 {
			ping.Kind = "info"
		} else {
			ping.Kind = q["t"][0]
		}

		ping.Time = time.Now()

		c <- ping
		fmt.Fprintf(w, "Recived, forwarded to TG \n\n[%s]\n%s", ping.Kind, ping.Val)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Pinging yourself as effortless as curl, check @Ping_webhook_bot telegram")
}

// Server is http reciever for the incoming pings
func Server(c chan Ping) {
	fmt.Println("Server starting")

	db := DBConnect()

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/{uuid:[-a-zA-Z0-9]+}", makeMessageHandler(db, c))

	http.Handle("/", r)

	fmt.Println("Server Running")
	log.Fatal(http.ListenAndServe(":8080", nil)) // WTF
}
