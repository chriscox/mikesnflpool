package mikesnflpool

import (
  "net/http"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/cors"
  "server/mikesnflpool/user"
  "server/mikesnflpool/teams"
  "server/mikesnflpool/tournaments"
  "server/mikesnflpool/games"
)

func init() {
  m := martini.Classic()
  // Games
  m.Get("/api/season/:season/games", games.GameHandler)
  
  // Teams
  m.Get("/api/teams", teams.TeamHandler)

  // User
  // m.Get("/api/season/:season/tournament/:tournament/user/:user/userpicks", user.UserPickHandler)
  m.Post("/api/userpicks", user.UserPickHandler)
  m.Post("/api/makepicks", user.AddUserPickHandler)

  // Auth
  m.Post("/api/login", user.LoginHandler)
  m.Post("/api/auth", user.UserRegistrationHandler)

  // Tournament
  m.Get("/api/season/:season/tournaments", tournaments.TournamentHandler)
  m.Post("/api/tournaments", tournaments.AddTournamentHandler)

  // Admin
  m.Post("/api/games", games.AddGameHandler)
  m.Post("/api/teams", teams.AddTeamHandler)

  // CORS
  m.Use(cors.Allow(&cors.Options{
    AllowOrigins:     []string{"*"},
    AllowMethods:     []string{"POST", "GET"},
    AllowHeaders:     []string{"Origin", "0", "1", "2", "If-Modified-Since", "Content-Type"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
  }))

  http.Handle("/", m)

}


// func profile(w http.ResponseWriter, r *http.Request) {
//   params := r.URL.Query()
//   name := params.Get(":name")
//   w.Write([]byte("Hello " + name))
// }
