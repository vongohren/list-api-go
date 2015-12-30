package handlers

import(
  "net/http"
  "encoding/json"
  "github.com/Snorlock/shoppingApi/db"
  "github.com/Snorlock/mux"
)

type Apis struct {
  Paths []Route
}

type Route struct {
  Path  string
  Methods []string
}

func IndexHandler(env *db.Env, routes []*mux.Route, w http.ResponseWriter, r *http.Request) error {
  apiPaths := []Route{}
  apis := Apis{apiPaths}
  for _, t := range routes {
    route := Route{t.GetPath(), t.GetMethods()}
		apiPaths = append(apiPaths, route)
	}
  apis.Paths = apiPaths
  js, err := json.Marshal(apis)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return nil
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
  return nil
}
