package mikesnflpool

import (
  // "fmt"
  "net/http"
  // "github.com/routes"
  // "github.com/gorilla/mux"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/cors"
  // "github.com/codegangsta/inject"
  // "server/games"
  "server/teams"

)

func init() {

  // mux := routes.New()
  // // mux.Get("/user", profile)
  // // mux.Options("/", optionsHandler)
  // mux.Get("/api/teams", teams.TeamHandler)
  // mux.Post("/api/teams", teams.AddTeamHandler)
  // http.Handle("/", &RouteMux{mux})
  // http.Handle("/", mux)


  // r := mux.NewRouter()
  // r.HandleFunc("/api/games", games.GameHandler).Methods("GET")
  // r.HandleFunc("/api/games", games.AddGameHandler).Methods("POST")
  // r.HandleFunc("/api/teams", teams.TeamHandler).Methods("GET")
  // r.HandleFunc("/api/teams", teams.AddTeamHandler).Methods("POST")
  // http.Handle("/", &MyServer{r})

  m := martini.Classic()
  m.Get("/api/teams", teams.TeamHandler)
  m.Post("/api/teams", teams.AddTeamHandler)
  // m.Get("/api/teams", func() string {
  //   return "Hello world!"
  // })

  m.Use(cors.Allow(&cors.Options{
    AllowOrigins:     []string{"*"},
    AllowMethods:     []string{"POST", "PUT", "PATCH", "GET"},
    AllowHeaders:     []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
  }))

  http.Handle("/", m)

}

// type MyServer struct {
//     r *mux.Router
// }

// func (s *MyServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
//     if origin := req.Header.Get("Origin"); origin != "" {
//         rw.Header().Set("Access-Control-Allow-Origin", origin)
//         rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
//         rw.Header().Set("Access-Control-Allow-Headers",
//             "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
//     }
//     // Stop here if its Preflighted OPTIONS request
//     if req.Method == "OPTIONS" {
//         return
//     }
//     // Lets Gorilla work
//     s.r.ServeHTTP(rw, req)
// }

// func profile(w http.ResponseWriter, r *http.Request) {
//   params := r.URL.Query()
//   name := params.Get(":name")
//   w.Write([]byte("Hello " + name))
// }

// func handler(w http.ResponseWriter, r *http.Request) {
//   fmt.Fprint(w, "Hello, world! It's working.")
// }

// func fooHandler(w http.ResponseWriter, r *http.Request) {
//   fmt.Fprintf(w, "Thanks for the %s!", r.Method)
// }

