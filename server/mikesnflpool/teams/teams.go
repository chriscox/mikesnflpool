package teams

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"net/http"
	m "server/mikesnflpool/models"
	"server/mikesnflpool/utils"
	"strconv"
)

func getStandingsCacheKey(season int) string {
	return "teamStandings" + strconv.Itoa(season)
}

func GetTeams(c appengine.Context) (teams []m.Team, err error) {
	var cachedTeams []m.Team
	if _, err := memcache.JSON.Get(c, "teams", &cachedTeams); err != nil {
		// Get team standings
		seasonStandings, err := getStandings(2016, c)
		if err != nil {
			panic(err.Error())
		}

		// Not in cache, so fetch item
		var teams []m.Team
		q := datastore.NewQuery("Team")
		keys, err := q.GetAll(c, &teams)
		if err != nil {
			return nil, err
		}

		// Associate teamKey and standings with teams
		for i := range teams {
			team := &teams[i]
			team.TeamKey = keys[i]
			team.TeamStandings = cumulateStandings(teams[i], seasonStandings)
		}
		// Add to memcache
		item := &memcache.Item{
			Key:    "teams",
			Object: teams,
		}
		memcache.JSON.Add(c, item)
		return teams, nil
	} else {
		// Found in cache
		c.Infof("Teams successfully retrieved from cache.")
		return cachedTeams, nil
	}
}

func getStandings(season int, c appengine.Context) (teamStandings []m.TeamStandings, err error) {
	var cachedTeamStandings []m.TeamStandings
	var cacheKey = getStandingsCacheKey(season)
	if _, err := memcache.JSON.Get(c, cacheKey, &cachedTeamStandings); err != nil {
		// Not in cache, so fetch item
		var teamStandings []m.TeamStandings
		q := datastore.NewQuery("TeamStandings").Filter("Season =", season)
		_, err := q.GetAll(c, &teamStandings)
		if err != nil {
			panic(err.Error())
		}

		// Add to memcache
		item := &memcache.Item{
			Key:    cacheKey,
			Object: teamStandings,
		}
		err = memcache.JSON.Add(c, item)
		return teamStandings, nil
	} else {
		// Found in cache
		c.Infof("TeamStandings successfully retrieved from cache.")
		return cachedTeamStandings, nil
	}
}

func cumulateStandings(team m.Team, seasonStandings []m.TeamStandings) (cumulateStandings m.TeamStandings) {
	var cumulativeStandings m.TeamStandings
	for _, standing := range seasonStandings {
		if standing.TeamKey.Equal(team.TeamKey) {
			cumulativeStandings.Home.Wins += standing.Home.Wins
			cumulativeStandings.Away.Wins += standing.Away.Wins
			cumulativeStandings.Total.Wins += (standing.Home.Wins + standing.Away.Wins)
			cumulativeStandings.Home.Losses += standing.Home.Losses
			cumulativeStandings.Away.Losses += standing.Away.Losses
			cumulativeStandings.Total.Losses += (standing.Home.Losses + standing.Away.Losses)
			cumulativeStandings.Home.Ties += standing.Home.Ties
			cumulativeStandings.Away.Ties += standing.Away.Ties
			cumulativeStandings.Total.Ties += (standing.Home.Ties + standing.Away.Ties)
		}
	}
	return cumulativeStandings
}

func TeamHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	teams, err := GetTeams(c)
	if err != nil {
		panic(err.Error())
	}
	utils.ServeJson(w, &teams)
}

func AddTeamHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var t m.Team
	if err := utils.ReadJson(r, &t); err != nil {
		panic(err.Error())
	}

	key := datastore.NewKey(c, "Team", t.Abbr, 0, nil)
	if _, err := datastore.Put(c, key, &t); err != nil {
		panic(err.Error())
	}
	utils.ServeJson(w, &t)
}
