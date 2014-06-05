package games

import (
  "fmt"
  "net/http"
)

func GameHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "my games.")
}

