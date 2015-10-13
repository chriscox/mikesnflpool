package games

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"github.com/go-martini/martini"
	"net/http"
	"server/mikesnflpool/teams"
	"server/mikesnflpool/utils"
	m "server/mikesnflpool/models"
	"strconv"
)

func getCacheKey(season int, week int) string {
	return "games" + strconv.Itoa(season) + strconv.Itoa(week)
}

func GetGames(season int, week int, c appengine.Context) (games []m.Game, err error) {
	var cachedGames []m.Game
	var cacheKey = getCacheKey(season, week)
	if _, err := memcache.JSON.Get(c, cacheKey, &cachedGames); err != nil {
		// Not in cache, so fetch item
		games := make([]m.Game, 0)
		q := datastore.NewQuery("Game").
			Filter("Season =", season).
			Filter("Week =", week).
			Order("Date")
		keys, err := q.GetAll(c, &games)
		if err != nil {
			panic(err.Error)
		}

		for i, _ := range games {
			games[i].GameKey = keys[i]
		}

		// Add to memcache
		item := &memcache.Item{
			Key:    cacheKey,
			Object: games,
		}
		err = memcache.JSON.Add(c, item)
		return games, nil
	} else {
		// Found in cache
		c.Infof("Games successfully retrieved from cache.")
		return cachedGames, nil
	}
}

func GameHandler(params martini.Params, w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	season, _ := strconv.Atoi(params["s"])
	week, _ := strconv.Atoi(r.URL.Query().Get("week"))

	// Get teams
	teams, err := teams.GetTeams(c)
	if err != nil {
		panic(err.Error)
	}

	// Get games
	games, err := GetGames(season, week, c)
	if err != nil {
		panic(err.Error)
	}

	// Associate team with game
	for i := range games {
		game := &games[i]
		for j, t := range teams {
			if game.AwayTeamKey.Equal(t.TeamKey) {
				game.AwayTeam = teams[j]
				continue
			}
			if game.HomeTeamKey.Equal(t.TeamKey) {
				game.HomeTeam = teams[j]
				continue
			}
		}
	}

	utils.ServeJson(w, &games)
}

func AddOrUpdateGameHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var g m.Game
	if err := utils.ReadJson(r, &g); err != nil {
		panic(err.Error())
	}

	// Assign team keys
	g.AwayTeamKey = datastore.NewKey(c, "Team", g.AwayTeamAbbr, 0, nil)
	g.HomeTeamKey = datastore.NewKey(c, "Team", g.HomeTeamAbbr, 0, nil)

	// Get teams
	teams, err := teams.GetTeams(c)
	if err != nil {
		panic(err.Error)
	}

	// Associate team with game
	for i, t := range teams {
		if g.AwayTeamKey.Equal(t.TeamKey) {
			g.AwayTeam = teams[i]
			continue
		}
		if g.HomeTeamKey.Equal(t.TeamKey) {
			g.HomeTeam = teams[i]
			continue
		}
	}

	// Clear games cache
	memcache.Delete(c, getCacheKey(g.Season, g.Week))

	// Check if existing game
	var existingGame m.Game
	err = datastore.Get(c, g.GameKey, &existingGame)
	isNew := err != nil

	if isNew {
		// Add game
		key := datastore.NewIncompleteKey(c, "Game", nil)
		gameKey, err := datastore.Put(c, key, &g)
		if err != nil {
			panic(err.Error)
		}
		g.GameKey = gameKey
		utils.ServeJson(w, &g)
	} else {
		// Update existing game
		existingGame.AwayTeamScore = g.AwayTeamScore
		existingGame.HomeTeamScore = g.HomeTeamScore
		existingGame.AwayTeamSpread = g.AwayTeamSpread
		existingGame.HomeTeamSpread = g.HomeTeamSpread
		existingGame.Ended = g.Ended
		if _, err := datastore.Put(c, g.GameKey, &existingGame); err != nil {
			panic(err.Error)
		}
	}
}

func DeleteGameHandler(params martini.Params, w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	season, _ := strconv.Atoi(params["s"])
	week, _ := strconv.Atoi(params["w"])
	gameKey, err := datastore.DecodeKey(params["g"])
	if err != nil {
		panic(err.Error)
	}

	// Get user pick keys for this game
	q := datastore.NewQuery("UserPick").Filter("GameKey =", gameKey).KeysOnly()
	keys, err := q.GetAll(c, nil)
	if err != nil {
		panic(err.Error)
	}

	// Clear games cache
	memcache.Delete(c, getCacheKey(season, week))

	// Delete keys
	keys = append(keys, gameKey)
	if err := datastore.DeleteMulti(c, keys); err != nil {
		panic(err.Error)
	}
}
