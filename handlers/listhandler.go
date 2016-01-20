package handlers

import(
  "fmt"
  "time"
  "net/http"
  "encoding/json"
  "github.com/Snorlock/shoppingApi/models"
  "github.com/Snorlock/shoppingApi/db"
  "github.com/Snorlock/mux"
  re "github.com/dancannon/gorethink"
)

type ItemRequest struct {
  ListId    string `json:"Id"`
  Text      string `json:"Text"`
}

func AddItemHandler(env *db.Env, id interface{}, w http.ResponseWriter, r *http.Request) error {
  //TODO Dont allow items without valid list id
  dec := json.NewDecoder(r.Body)
  _ = "breakpoint"
  var request ItemRequest
  error2 := dec.Decode(&request)
  if error2 != nil {
    fmt.Println(error2);
    return error2
  }
  item := models.Item{"", request.Text, "false", time.Now()}

  itemResp, err := re.DB(env.DBName).Table(env.ItemsTable).Insert(item).RunWrite(env.DBSession)
  if err != nil {
    return err
  }
  itemId := itemResp.GeneratedKeys[0]
  _, err3 := re.DB(env.DBName).Table(env.ListsTable).Get(request.ListId).Update(map[string]interface{}{"Items":re.Row.Field("Items").Append(itemId), "Update":time.Now()}).Run(env.DBSession)
  fmt.Println(itemId)
  fmt.Println(err3)
  return nil
}

type ListRequest struct {
  Title      string `json:"Title"`
}

func AddListHandler(env *db.Env, id interface{}, w http.ResponseWriter, r *http.Request) error {
  dec := json.NewDecoder(r.Body)
  var request ListRequest
  error2 := dec.Decode(&request)
  if error2 != nil {
    fmt.Println(error2);
    return error2
  }
  _ = "breakpoint"
  email, _ := id.(string)
  owners := []string{email}
  list := models.List{"", []string{}, owners, request.Title, time.Now(),}
  insertedList, err := re.DB(env.DBName).Table(env.ListsTable).Insert(list).RunWrite(env.DBSession)
  if err != nil {
    return err
  }
  list.Id = insertedList.GeneratedKeys[0]
  js, err := json.Marshal(list)
  if err != nil {
    return err
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
  return nil
}

func GetListsHandler(env *db.Env, id interface{}, w http.ResponseWriter, r *http.Request) error {
  email, _ := id.(string)
  res, err2 := re.DB(env.DBName).Table(env.ListsTable).GetAllByIndex("Owners",email).OrderBy(re.Desc("Updated")).Run(env.DBSession)
  defer res.Close()
  // Scan query result into the person variable
  lists := []models.List{}
  err2 = res.All(&lists)
  if err2 != nil {
      fmt.Printf("Error scanning database result: %s", err2)
      return err2
  }
  if len(lists) > 0 {
    js, err := json.Marshal(lists)
    if err != nil {
      return err
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
    return nil
  }
  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  w.WriteHeader(http.StatusNoContent)

  return nil
}

func GetListDetailHandler(env *db.Env, id interface{}, w http.ResponseWriter, r *http.Request) error {
  params := mux.Vars(r)
  id = params["id"]
  res, err2 := re.DB(env.DBName).Table(env.ListsTable).Get(id).Run(env.DBSession)
  defer res.Close()
  _ = "breakpoint"
  // Scan query result into the person variable
  list := models.List{}
  err2 = res.One(&list)
  if err2 != nil {
      fmt.Printf("Error scanning database result: %s", err2)
      return err2
  }
  detailedList := models.DetailedList{list.Id, []models.Item{}, list.Owners, list.Title, list.Updated,}
  b := make([]interface{}, len(list.Items))
  for i := range list.Items {
      b[i] = list.Items[i]
  }
  itemsRes, err := re.DB(env.DBName).Table(env.ItemsTable).GetAll(b...).OrderBy(re.Desc("Updated")).Run(env.DBSession)
  if err !=nil {
    return err
  }
  items := []models.Item{}
  defer itemsRes.Close()
  err = itemsRes.All(&items)
  if err !=nil {
    return err
  }
  detailedList.Items = items

  js, err := json.Marshal(detailedList)
  if err != nil {
    return err
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
  return nil
}
