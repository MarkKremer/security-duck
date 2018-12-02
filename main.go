package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"encoding/json"
)

// https://www.sohamkamani.com/blog/2017/09/13/how-to-build-a-web-application-in-golang/

type ConversationResponse struct {
	Text           string `json:"text"`
	DuckSays       bool   `json:"duckSays"`
	MakeQuackSound bool   `json:"makeQuackSound"`
}

func newRouter() *mux.Router {
	r := mux.NewRouter()

	// Handle methods
	r.HandleFunc("/converse", conversationHandler).Methods("POST")

	// Serve /... to /static/...
	staticFileDirectory := http.Dir("./static/")
	staticFileHandler := http.FileServer(staticFileDirectory)
	r.PathPrefix("/").Handler(staticFileHandler).Methods("GET")

	return r
}

func main() {
	r := newRouter()
	err := http.ListenAndServe(":8080", r) // :80
	if err != nil {
		log.Fatal(err)
	}
}

func conversationHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, "Security duck says: <b>quack.</b>")

	response := ConversationResponse{
		Text:           "quack.",
		DuckSays:       true,
		MakeQuackSound: true,
	}

	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}