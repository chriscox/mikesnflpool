package teamstandings

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"github.com/go-martini/martini"
	"net/http"
	"server/mikesnflpool/games"
	// "server/mikesnflpool/teams"
	"server/mikesnflpool/utils"
	"strconv"
)

type TeamStandings struct {
	Season  int            `json:"season"`
	Week    int            `json:"week"`
	TeamKey *datastore.Key `json:"teamKey"`
	Home    Totals         `json:"home"`
	Away    Totals         `json:"away"`
	Total   Totals         `json:"total"`
}

type Totals struct {
	Wins   int `json:"wins"`
	Losses int `json:"losses"`
	Ties   int `json:"ties"`
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

func getCacheKey(season int, week int) string {
	return "teamStandings" + strconv.Itoa(season) + strconv.Itoa(week)
}

func GetTeamStandings(season int, week int, c appengine.Context) (teamStandings []TeamStandings, err error) {
	var cachedTeamStandings []TeamStandings
	var cacheKey = getCacheKey(season, week)
	if _, err := memcache.JSON.Get(c, cacheKey, &cachedTeamStandings); err != nil {
		// Not in cache, so fetch item
		var teamStandings []TeamStandings
		q := datastore.NewQuery("TeamStandings").
			Filter("Season =", season).
			Filter("Week =", week).
			Order("Week")
		_, err := q.GetAll(c, &teamStandings)
		if err != nil {
			panic(err.Error)
		}
		c.Infof("%v", teamStandings)

		// for i, _ := range teamStandings {
		//   teamStandings[i].TeamKey = keys[i]
		// }

		// Add to memcache
		// item := &memcache.Item {
		//    Key: cacheKey,
		//    Object: teamStandings,
		// }
		// err = memcache.JSON.Add(c, item)
		return teamStandings, nil
	} else {
		// Found in cache
		c.Infof("TeamStandings successfully retrieved from cache.")
		return cachedTeamStandings, nil
	}
}

func TeamStandingsHandler(params martini.Params, w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	season, _ := strconv.Atoi(params["s"])
	week, _ := strconv.Atoi(r.URL.Query().Get("week"))

	UpdateTeamStandingsHandler(params, w, r)
	// clearStandings(w, r)

	standings, err := GetTeamStandings(season, week, c)
	if err != nil {
		panic(err.Error)
	}
	utils.ServeJson(w, &standings)
}

func UpdateTeamStandingsHandler(params martini.Params, w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	season, _ := strconv.Atoi(params["s"])
	// team,_ := params["t"]

	// if team == "ALL" {
	// var teams []teams.Team
	// q := datastore.NewQuery("Team")
	// if _, err := q.GetAll(c, &teams); err != nil {
	//   panic(err.Error)
	// }

	// Get teams
	// _, err := teams.GetTeams(c)
	// if err != nil {
	//   panic(err.Error)
	// }

	// var teamStandings []TeamStandings
	// q := datastore.NewQuery("TeamStandings").
	//        Filter("Season =", season).
	//        Filter("Week <", week).
	//        Order("Week")
	// _, err := q.GetAll(c, &teamStandings)
	// if err != nil {
	//   panic(err.Error)
	// }
	// c.Infof("%v", teamStandings)

	for week := 1; week < 18; week++ {
		// Get games for each week
		var allGames []games.Game
		allGames, err := games.GetGames(season, week, c)
		if err != nil {
			panic(err.Error)
		}

		// Arrays for PutMulti
		var standings []TeamStandings
		var standingsKeys []*datastore.Key

		// Calculate standings
		for _, game := range allGames {
			if game.Ended {
				// Away team standings
				key, item := newStanding(season, week, game, "away")
				standings = append(standings, item)
				standingsKeys = append(standingsKeys, datastore.NewKey(c, "TeamStandings", key, 0, nil))
				// Home team standings
				key, item = newStanding(season, week, game, "home")
				standings = append(standings, item)
				standingsKeys = append(standingsKeys, datastore.NewKey(c, "TeamStandings", key, 0, nil))
			}
		}

		if _, err := datastore.PutMulti(c, standingsKeys, standings); err != nil {
			panic(err.Error)
		}
	}
}

func newStanding(season int, week int, game games.Game, playingAs string) (key string, teamStandings TeamStandings) {
	var standings TeamStandings
	standings.Season = season
	standings.Week = week

	if playingAs == "away" {
		if game.AwayTeamScore > game.HomeTeamScore {
			standings.Away.Wins += 1
			standings.Total.Wins += 1
		} else if game.AwayTeamScore < game.HomeTeamScore {
			standings.Away.Losses += 1
			standings.Total.Losses += 1
		} else {
			standings.Away.Ties += 1
			standings.Total.Ties += 1
		}
		standings.TeamKey = game.AwayTeamKey
	}

	if playingAs == "home" {
		if game.HomeTeamScore > game.AwayTeamScore {
			standings.Home.Wins += 1
			standings.Total.Wins += 1
		} else if game.HomeTeamScore < game.AwayTeamScore {
			standings.Home.Losses += 1
			standings.Total.Losses += 1
		} else {
			standings.Home.Ties += 1
			standings.Total.Ties += 1
		}
		standings.TeamKey = game.HomeTeamKey
	}

	key = standings.TeamKey.Encode() + "_" + strconv.Itoa(season) + strconv.Itoa(week)
	return key, standings
}

func clearStandings(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var teamStandings []TeamStandings
	q := datastore.NewQuery("TeamStandings")
	keys, _ := q.GetAll(c, &teamStandings)
	_ = datastore.DeleteMulti(c, keys)
}
