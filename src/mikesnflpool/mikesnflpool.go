package mikesnflpool

import (
  "fmt"
  "net/http"
  "github.com/gorilla/mux"
  "games"
  "teams"
)

func init() {
  r := mux.NewRouter()
  r.HandleFunc("/", handler)
  r.HandleFunc("/games", games.GameHandler)
  r.HandleFunc("/teams", teams.TeamHandler)
  r.HandleFunc("/conferences", teams.ConferenceHandler)
  r.HandleFunc("/root", teams.Root)
  r.HandleFunc("/sign", teams.Sign)
  http.Handle("/", r)
}

func handler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Hello, world! It's working.")
}

func fooHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Thanks for the %s!", r.Method)
}

