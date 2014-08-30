package teamstandings

import (
  "appengine"
  "appengine/datastore"
  "github.com/go-martini/martini"
  "net/http"
  "server/mikesnflpool/teams"
  // "server/mikesnflpool/utils"
  // "strconv"
)

type TeamStandings struct {
  Season            int             `json:"season"`
  Week              int             `json:"week"`
  // TeamKey           *datastore.Key  `json:"teamKey"`
  Home              Totals          `json:"home"`         
  Away              Totals          `json:"away"` 
  Total             Totals          `json:"total"` 
}

type Totals struct {
  Wins    int         `json:"wins"`
  Losses  int         `json:"losses"`
  Ties    int         `json:"ties"`
}
  
// type GameType struct {
//   Home  Totals         `json:"home"`
//   Away  Totals         `json:"away"`
//   Total Totals         `json:"total"`
// }
// type TeamStandings      map[string]GameType

// type StatsMap struct {
//   Stats            TeamStandings      `json:"stats"`
// }


// func UpdateTeamStandings(week int, tournamentKey *datastore.Key, teams []Team) {

// }

func TeamStandingsHandler(params martini.Params, w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  // season,_ := strconv.Atoi(params["s"])
  // team,_ := params["t"]

  // if team == "ALL" {
    // var teams []teams.Team
    // q := datastore.NewQuery("Team")
    // if _, err := q.GetAll(c, &teams); err != nil {
    //   panic(err.Error)
    // }

    // Get teams
    teams, err := teams.GetTeams(c)
    if err != nil {
      panic(err.Error)
    }

    c.Infof("%v", teams)


    var standings TeamStandings
    standings.Season = 2014
    standings.Week = 1
    standings.Home.Wins = 2
    standings.Home.Losses = 1
    standings.Home.Ties = 3
    standings.Away.Wins = 2
    standings.Away.Losses = 1
    standings.Away.Ties = 3
    standings.Total.Wins = 2
    standings.Total.Losses = 1
    standings.Total.Ties = 3

    key := datastore.NewIncompleteKey(c, "TeamStandings", nil)
    if _, err := datastore.Put(c, key, &standings); err != nil {
      panic(err.Error)
    }

    c.Infof("%v", standings)


    // var teamStandings = make(map[string]GameType)
    // for _, t := range teams {
    //   var totals Totals

    //   totals.Wins = 1
    //   totals.Losses = 2
    //   totals.Ties = 3

    //   var gameType GameType
    //   gameType.Home = totals
    //   gameType.Away = totals
    //   gameType.Total = totals

    //   teamStandings[t.Abbr] = gameType
    // }


    // var statsMap StatsMap
    // statsMap.Stats = teamStandings
    // c.Infof("%v", statsMap)
    // utils.ServeJson(w, &statsMap)
}

// func calculateStandings(week int, team Team) {


// }

// func tallyScores() {

// }

