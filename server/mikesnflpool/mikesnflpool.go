package mikesnflpool

import (
  "net/http"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/cors"
  "server/auth"
  "server/games"
  "server/teams"

)

func init() {
  m := martini.Classic()
  // Games.
  // GET            /api/v1/season/:s/games      controllers.api.Games.list(s:Integer, week:Integer ?= null)

  m.Get("/api/season/:s/games", games.GameHandler)
  m.Post("/api/games", games.AddGameHandler)

  // Teams.
  m.Get("/api/teams", teams.TeamHandler)
  m.Post("/api/teams", teams.AddTeamHandler)

  // Auth.
  m.Post("/api/auth/register", auth.RegisterHandler)

  // CORS
  m.Use(cors.Allow(&cors.Options{
    AllowOrigins:     []string{"*"},
    AllowMethods:     []string{"POST", "GET"},
    AllowHeaders:     []string{"Content-Type"},
    AllowCredentials: true,
  }))

  http.Handle("/", m)

}


// func profile(w http.ResponseWriter, r *http.Request) {
//   params := r.URL.Query()
//   name := params.Get(":name")
//   w.Write([]byte("Hello " + name))
// }
