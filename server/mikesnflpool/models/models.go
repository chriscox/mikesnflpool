package models

import (
	"appengine/datastore"
	"time"
)

type Tournament struct {
	Name string `json:"name"`
}

type GameEvent struct {
	TournamentKey *datastore.Key `json:"tournament"`
	Season        int            `json:"season"`
	Week          int            `json:"week"`
}

type Game struct {
	Season  int            `json:"season"`
	Week    int            `json:"week"`
	Date    time.Time      `json:"date"`
	GameKey *datastore.Key `json:"gameKey" datastore:"-"`
	Ended   bool           `json:"ended"`

	AwayTeamKey    *datastore.Key 	`json:"awayTeamKey"`
	AwayTeamAbbr   string         	`json:"awayTeamAbbr" datastore:"-"`
	AwayTeam       Team     		`json:"awayTeam" datastore:"-"`
	AwayTeamScore  int          	`json:"awayTeamScore"`
	AwayTeamSpread float32        	`json:"awayTeamSpread"`

	HomeTeamKey    *datastore.Key 	`json:"homeTeamKey"`
	HomeTeamAbbr   string         	`json:"homeTeamAbbr" datastore:"-"`
	HomeTeam       Team     		`json:"homeTeam" datastore:"-"`
	HomeTeamScore  int            	`json:"homeTeamScore"`
	HomeTeamSpread float32        	`json:"homeTeamSpread"`
}


type Team struct {
	Abbr     string         	`json:"abbr"`
	Name     string         	`json:"name"`
	NickName string         	`json:"nickName"`
	Division string         	`json:"division"`
	Selected bool           	`json:"selected" datastore:"-"`
	TeamKey  *datastore.Key 	`json:"teamKey" datastore:"-"`
	TeamStandings TeamStandings `json:"standings" datastore:"-"`
}

type TeamStandings struct {
	Season  	int            `json:"season"`
	Week    	int            `json:"week"`
	TeamKey 	*datastore.Key `json:"teamKey"`
	Home    	StandingTotals     `json:"home"`
	Away    	StandingTotals     `json:"away"`
	Total 		StandingTotals     `json:"total"`
}

type StandingTotals struct {
	Wins   int `json:"wins"`
	Losses int `json:"losses"`
	Ties   int `json:"ties"`
}