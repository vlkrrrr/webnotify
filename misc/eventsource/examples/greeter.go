package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/antage/eventsource.v1"
)

func main() {
	es := eventsource.New(nil, nil)
	defer es.Close()
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/events", es)
	i := 0
	go func() {
		for {
			es.SendEventMessage("hello"+strconv.Itoa(int(time.Now().Unix())), "", "")
			log.Printf("Hello has been sent (consumers: %d)", es.ConsumersCount())
			time.Sleep(2 * time.Second)
			i++
			if i == 5 {
				es.Close()
				break
			}
		}
	}()
	log.Print("Open URL http://localhost:8080/ in your browser.")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
