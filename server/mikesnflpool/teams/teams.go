package teams

import (
  "appengine"
  "appengine/datastore"
  "appengine/memcache"
  "net/http"
  "server/mikesnflpool/utils"
)

type Team struct {
  Abbr      string          `json:"abbr"`
  Name      string          `json:"name"`
  NickName  string          `json:"nickName"`
  Division  string          `json:"division"`
  Selected  bool            `json:"selected" datastore:"-"`
  TeamKey   *datastore.Key  `json:"teamKey" datastore:"-"`
}

func GetTeams(c appengine.Context) (teams []Team, err error) {
  var cachedTeams []Team
  if _, err := memcache.JSON.Get(c, "teams", &cachedTeams); err != nil {
    // Not in cache, so fetch item
    var teams []Team
    q := datastore.NewQuery("Team")
    keys, err := q.GetAll(c, &teams)
    if err != nil {
      return nil, err
    }
    for i, _ := range teams {
      teams[i].TeamKey = keys[i]
    }
    
    // Add to memcache
    item := &memcache.Item {
       Key: "teams",
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

func TeamHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  teams, err := GetTeams(c)
  if err != nil {
    panic(err.Error)
  }
  utils.ServeJson(w, &teams)
}

func AddTeamHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  var t Team
  if err := utils.ReadJson(r, &t); err != nil {
    panic(err.Error)
  }

  key := datastore.NewKey(c, "Team", t.Abbr, 0, nil)
  if _, err := datastore.Put(c, key, &t); err != nil {
    panic(err.Error)
  }
  utils.ServeJson(w, &t)
}

