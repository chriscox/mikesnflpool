package teams

import (
  // "fmt"
  // "encoding/json"
  "net/http"
  "html/template"
  "time"
  // "io/ioutil"
  // "github.com/routes"

  "appengine"
  "appengine/datastore"
  "appengine/user"

  "server/utils"

  // "github.com/go-martini/martini"
)

type Team struct {
  Name  string
  NickName string
}

type Conference struct {
  name  string
}

// guestbookKey returns the key used for all guestbook entries.
// func teamkKey(c appengine.Context) *datastore.Key {
//   // The string "default_guestbook" here could be varied to have multiple guestbooks.
//   return datastore.NewKey(c, "Team", "default_team", 0, nil)
// }

func ConferenceHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  q := datastore.NewQuery("Conference").Limit(2)
  var conferences []Conference
  if _, err := q.GetAll(c, &conferences); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  if err := conferenceTemplate.Execute(w, conferences); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

// guestbookKey returns the key used for all guestbook entries.
// func teamKey(c appengine.Context) *datastore.Key {
//         // The string "default_guestbook" here could be varied to have multiple guestbooks.
//         return datastore.NewKey(c, "Team", "default_team", 0, nil)
// }

// type test_struct struct {
//     Name string
//     NickName string
// }


func AddTeamHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  var t Team
  if err := utils.ReadJson(r, &t); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  key := datastore.NewKey(c, "Team", t.Name, 0, nil)
  if _, err := datastore.Put(c, key, &t); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  utils.ServeJson(w, &t)
}

func TeamHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  q := datastore.NewQuery("Team").Limit(32)
  teams := make([]Team, 0, 32)
  if _, err := q.GetAll(c, &teams); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  utils.ServeJson(w, &teams)
  // w.Write([]byte("Hello"))
  // routes.ServeJson(w, &teams)
  // w.WriteJson(&teams)
  // gorca.WriteJSON(c, w, r, teams)
  // if err := teamTemplate.Execute(w, teams); err != nil {
  //   http.Error(w, err.Error(), http.StatusInternalServerError)
  // }
  // if err := teamTemplate.Execute(w, teams); err != nil {
  //   http.Error(w, err.Error(), http.StatusInternalServerError)
  // }
}

var conferenceTemplate = template.Must(template.New("conference").Parse(`
<html>
  <head>
    <title>Conferences</title>
  </head>
  <body>
    OK
  </body>
</html>
`))

var teamTemplate = template.Must(template.New("team").Parse(`
<html>
  <head>
    <title>Teams</title>
  </head>
  <body>
    {{range .}}
      <div>{{.Name}}</div>
    {{end}}
  </body>
</html>
`))


type Greeting struct {
        Author  string
        Content string
        Date    time.Time
}

// guestbookKey returns the key used for all guestbook entries.
func guestbookKey(c appengine.Context) *datastore.Key {
        // The string "default_guestbook" here could be varied to have multiple guestbooks.
        return datastore.NewKey(c, "Guestbook", "default_guestbook", 0, nil)
}

func Root(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        // Ancestor queries, as shown here, are strongly consistent with the High
        // Replication Datastore. Queries that span entity groups are eventually
        // consistent. If we omitted the .Ancestor from this query there would be
        // a slight chance that Greeting that had just been written would not
        // show up in a query.
        q := datastore.NewQuery("Greeting").Ancestor(guestbookKey(c)).Order("-Date").Limit(10)
        greetings := make([]Greeting, 0, 10)
        if _, err := q.GetAll(c, &greetings); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        if err := guestbookTemplate.Execute(w, greetings); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
        }
}


var guestbookTemplate = template.Must(template.New("book").Parse(`
<html>
  <head>
    <title>Go Guestbook</title>
  </head>
  <body>
    {{range .}}
      {{with .Author}}
        <p><b>{{.}}</b> wrote:</p>
      {{else}}
        <p>An anonymous person wrote:</p>
      {{end}}
      <pre>{{.Content}}</pre>
    {{end}}
    <form action="/sign" method="post">
      <div><textarea name="content" rows="3" cols="60"></textarea></div>
      <div><input type="submit" value="Sign Guestbook"></div>
    </form>
  </body>
</html>
`))


func Sign(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        g := Greeting{
                Content: r.FormValue("content"),
                Date:    time.Now(),
        }
        if u := user.Current(c); u != nil {
                g.Author = u.String()
        }
        // We set the same parent key on every Greeting entity to ensure each Greeting
        // is in the same entity group. Queries across the single entity group
        // will be consistent. However, the write rate to a single entity group
        // should be limited to ~1/second.
        key := datastore.NewIncompleteKey(c, "Greeting", guestbookKey(c))
        _, err := datastore.Put(c, key, &g)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        http.Redirect(w, r, "/", http.StatusFound)
}

