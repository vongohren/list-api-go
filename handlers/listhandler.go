package handlers

import(
  "fmt"
  "time"
  "net/http"
  "encoding/json"
  "github.com/Snorlock/shoppingApi/models"
  "github.com/Snorlock/shoppingApi/db"
  re "github.com/dancannon/gorethink"
)

func AddHandler(env *db.Env, token interface{}, w http.ResponseWriter, r *http.Request) error {
  dec := json.NewDecoder(r.Body)

  var item models.Item
  error2 := dec.Decode(&item)
  if error2 != nil {
    fmt.Println(error2);
    return error2
  }
  item.Added = time.Now()
  item.Accuired = false
  items := []models.Item{}
  items = append(items,item)
  list := models.ShoppingList{"",items, token.(string)}
  _ = "breakpoint"
  res, err2 := re.DB(env.DBName).Table(env.ListsTable).Filter(map[string]interface{}{"Owner":token.(string) ,}).Run(env.DBSession)
  defer res.Close()

  // Scan query result into the person variable
  lists := []models.ShoppingList{}
  err2 = res.All(&lists)
  if err2 != nil {
      fmt.Printf("Error scanning database result: %s", err2)
      return err2
  }

  if len(lists) > 0 {
    existingItems := lists[0].Items
    existingItems = append(existingItems, item)
    _, err3 := re.DB(env.DBName).Table(env.ListsTable).Get(lists[0].Id).Update(map[string]interface{}{"Items":existingItems}).Run(env.DBSession)
    if err3 != nil {
      fmt.Println(err3);
      http.Error(w, err3.Error(), http.StatusInternalServerError)
      fmt.Fprint(w, "ERRORR!\n")
    }
  } else {
    _, err := re.DB(env.DBName).Table(env.ListsTable).Insert(list).RunWrite(env.DBSession)
    if err != nil {
      fmt.Println(err);
      http.Error(w, err.Error(), http.StatusInternalServerError)
      fmt.Fprint(w, "ERRORR!\n")
    }
  }


  return nil
}

func GetListHandler(env *db.Env, token interface{}, w http.ResponseWriter, r *http.Request) error {
  res, err2 := re.DB(env.DBName).Table(env.ListsTable).Filter(map[string]interface{}{"Owner":token.(string) ,}).Run(env.DBSession)
  defer res.Close()
  _ = "breakpoint"
  // Scan query result into the person variable
  lists := []models.ShoppingList{}
  err2 = res.All(&lists)
  if err2 != nil {
      fmt.Printf("Error scanning database result: %s", err2)
      return err2
  }
  if len(lists) > 0 {
    js, err := json.Marshal(lists[0])
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
