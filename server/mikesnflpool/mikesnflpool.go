package mikesnflpool

import (
  "net/http"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/cors"
  "server/mikesnflpool/games"
  "server/mikesnflpool/teams"
  "server/mikesnflpool/teamstandings"
  "server/mikesnflpool/tournaments"
  "server/mikesnflpool/user"
  "server/mikesnflpool/userpicks"
)

func init() {
  m := martini.Classic()
  // Games
  m.Get("/api/season/:s/games", games.GameHandler)
  // TODO: Convert this to a DELETE instead of POST
  //m.Delete("/api/games/:g", games.DeleteGameHandler)
  m.Post("/api/season/:s/week/:w/deletegame/:g", games.DeleteGameHandler)
  m.Post("/api/games", games.AddOrUpdateGameHandler)
  
  // Teams
  m.Get("/api/teams", teams.TeamHandler)
  m.Get("/api/season/:s/teams/:t/standings", teamstandings.TeamStandingsHandler)

  // User
  m.Get("/api/tournament/:t/users", user.UserHandler)
  m.Get("/api/tournament/:t/season/:s/userpicks", userpicks.AllUserPickHandler)
  m.Get("/api/tournament/:t/season/:s/user/:u/userpicks", userpicks.UserPickHandler)
  m.Post("/api/userpicks", userpicks.AddUserPickHandler)
  // m.Post("/api/botpicks", userpicks.UpdateBotPicksHandler)
  m.Get("/api/tournament/:t/season/:s/userstats", userpicks.UserStatsHandler)

  // Auth  
  m.Post("/api/login", user.LoginHandler)
  m.Post("/api/auth", user.UserRegistrationHandler)
  m.Post("/api/passwordforgot", user.PasswordForgot)
  m.Post("/api/passwordreset", user.PasswordReset)

  // Tournament
  m.Get("/api/season/:season/tournaments", tournaments.TournamentHandler)
  m.Post("/api/tournaments", tournaments.AddTournamentHandler)

  // Admin
  
  m.Post("/api/teams", teams.AddTeamHandler)

  // CORS
  m.Use(cors.Allow(&cors.Options{
    AllowOrigins:     []string{"*"},
    AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
    AllowHeaders:     []string{"*", "Origin", "0", "1", "2", "If-Modified-Since", "Content-Type"},
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
