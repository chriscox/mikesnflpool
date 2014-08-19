package utils

import (
  // "fmt"
  "encoding/json"
  "net/http"
  "strconv"
  "io/ioutil"
)

// ServeJson replies to the request with a JSON
// representation of resource v.
func ServeJson(w http.ResponseWriter, v interface{}) {
  content, err := json.MarshalIndent(v, "", "  ")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Length", strconv.Itoa(len(content)))
  w.Header().Set("Content-Type", "application/json")
  w.Write(content)
}

// ReadJson will parses the JSON-encoded data in the http
// Request object and stores the result in the value
// pointed to by v.
func ReadJson(r *http.Request, v interface{}) error {
  body, err := ioutil.ReadAll(r.Body)
  r.Body.Close()
  if err != nil {
    return err
  }
  return json.Unmarshal(body, v)
}