package tournaments

import (
	"appengine"
	"appengine/datastore"
	"github.com/go-martini/martini"
	"net/http"
	"server/mikesnflpool/utils"
	"strconv"
	"time"
)

type Tournament struct {
	Name          string         `json:"name"`
	Season        int            `json:"season"`
	TournamentKey *datastore.Key `json:"tournamentKey" datastore:"-"`
}

type TournamentUser struct {
	UserKey *datastore.Key `json:"userKey"`
	Admin   bool           `json:"-"`
	Season  int            `json:"season"`
}

type GameEvent struct {
	Season int `json:"season"`
	Week   int `json:"week"`
}

/**
 * Todo: THIS IS TEMPORARY. Remove when tournament users is fixed.
 */
type User struct {
	FirstName       string         `json:"firstName"`
	LastName        string         `json:"lastName"`
	Email           string         `json:"email"`
	Password        string         `json:"password,omitempty" datastore:"-"`
	SecurePassword  []byte         `json:",omitempty"`
	Token           string         `json:"token,omitempty"`
	TokenExpiration time.Time      `json:",omitempty"`
	UserKey         *datastore.Key `json:"userKey" datastore:"-"`
	TournamentKey   *datastore.Key `json:"tournamentKey" datastore:"-"`
	Admin           bool           `json:"admin,omitempty" datastore:"-"`
	Bot             bool           `json:"bot,omitempty"`
	BotType         string         `json:"botType,omitempty"`
}

func AllTournamentHandler(parms martini.Params, w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// Get all tournaments
	q := datastore.NewQuery("Tournament")
	var tournaments []Tournament
	keys, err := q.GetAll(c, &tournaments)
	if err != nil {
		panic(err.Error())
	}

	// Associate keys with tournament
	for i := range keys {
		t := &tournaments[i]
		t.TournamentKey = keys[i]
	}
	utils.ServeJson(w, &tournaments)
}

func TournamentHandler(parms martini.Params, w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	season, _ := strconv.Atoi(parms["season"])

	// Get tournaments
	q := datastore.NewQuery("Tournament").Filter("Season =", season)
	var tournaments []Tournament
	keys, err := q.GetAll(c, &tournaments)
	if err != nil {
		panic(err.Error())
	}

	// Associate keys with tournament
	for i := range keys {
		t := &tournaments[i]
		t.TournamentKey = keys[i]
	}
	utils.ServeJson(w, &tournaments)
}

func AddTournamentHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	var t Tournament
	if err := utils.ReadJson(r, &t); err != nil {
		panic(err.Error())
	}

	// Add tournament
	key := datastore.NewIncompleteKey(c, "Tournament", nil)
	tournyKey, err := datastore.Put(c, key, &t)
	if err != nil {
		panic(err.Error())
	}

	// Add Game Events
	for i := 1; i < 18; i++ {
		var gameEvent GameEvent
		gameEvent.Season = t.Season
		gameEvent.Week = i
		key := datastore.NewIncompleteKey(c, "GameEvent", tournyKey)
		if _, err := datastore.Put(c, key, &gameEvent); err != nil {
			panic(err.Error())
		}
	}

	// TODO: Remove adding users here. This is temporary.
	// Add tournament users
	q := datastore.NewQuery("User")
	var users []User
	userKeys, err := q.GetAll(c, &users)
	if err != nil {
		panic(err.Error())
	}
	for i := range userKeys {
		var tu TournamentUser
		tu.UserKey = userKeys[i]
		tu.Season = t.Season
		var tUserKey = datastore.NewIncompleteKey(c, "TournamentUser", tournyKey)
		_, err := datastore.Put(c, tUserKey, &tu)
		if err != nil {
			panic(err.Error())
		}
	}

	utils.ServeJson(w, &t)
}
