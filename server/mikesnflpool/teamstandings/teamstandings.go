package teamstandings

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"github.com/go-martini/martini"
	"net/http"
	"server/mikesnflpool/games"
	m "server/mikesnflpool/models"
	"strconv"
)

func getCacheKey(season int) string {
	return "teamStandings" + strconv.Itoa(season)
}

func UpdateTeamStandingsHandler(params martini.Params, w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	season, _ := strconv.Atoi(params["s"])

	clearStandings(season, c)

	for week := 1; week < 18; week++ {
		// Get games for each week
		var allGames []m.Game
		allGames, err := games.GetGames(season, week, c)
		if err != nil {
			panic(err.Error)
		}

		// Arrays for PutMulti
		var standings []m.TeamStandings
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

func newStanding(season int, week int, game m.Game, playingAs string) (key string, teamStandings m.TeamStandings) {
	var standings m.TeamStandings
	standings.Season = season
	standings.Week = week

	if playingAs == "away" {
		if game.AwayTeamScore > game.HomeTeamScore {
			standings.Away.Wins += 1
		} else if game.AwayTeamScore < game.HomeTeamScore {
			standings.Away.Losses += 1
		} else {
			standings.Away.Ties += 1
		}
		standings.TeamKey = game.AwayTeamKey
	}

	if playingAs == "home" {
		if game.HomeTeamScore > game.AwayTeamScore {
			standings.Home.Wins += 1
		} else if game.HomeTeamScore < game.AwayTeamScore {
			standings.Home.Losses += 1
		} else {
			standings.Home.Ties += 1
		}
		standings.TeamKey = game.HomeTeamKey
	}

	key = standings.TeamKey.Encode() + "_" + strconv.Itoa(season) + strconv.Itoa(week)
	return key, standings
}

func clearStandings(season int, c appengine.Context) {
	var teamStandings []m.TeamStandings
	q := datastore.NewQuery("TeamStandings")
	keys, _ := q.GetAll(c, &teamStandings)
	_ = datastore.DeleteMulti(c, keys)

	// Clear cache
	memcache.Delete(c, getCacheKey(season))
	memcache.Delete(c, "teams")
}
