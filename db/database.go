package db

import (
	re "github.com/dancannon/gorethink"
	"log"
  "fmt"
	"os"
	"time"
)

type Env struct {
    DBSession   *re.Session
		DBName 			string
		UsersTable 	string
		ListsTable 	string
		ItemsTable 	string
		ListKey			string
}

var (
	DBName string = "list_api"
	UsersTable string = "users"
	ListsTable string = "lists"
	ItemsTable string = "items"
	ListKey string = "Owners"
	session *re.Session
)

func StartDatabase() *Env{
  fmt.Println("db")
	for {
		connected := connectToDB()
		if(connected) {
			break;
		}
		fmt.Println("waiting for rethink %s:%s",os.Getenv("DB_PORT_28015_TCP_ADDR"), os.Getenv("DB_PORT_28015_TCP_PORT"))
		time.Sleep(2000 * time.Millisecond)
	}
  if session != nil {
    fmt.Println(session)
    fmt.Println("connectedzozozozoz");
  }
	var dbName = "list_api"
	resp, error := re.DBCreate(dbName).RunWrite(session)
	if error != nil {
		fmt.Println("DB creation either failed or DB exists already")
	}
	_, errz := re.DB(dbName).Table(UsersTable).Run(session)
	if errz != nil {
		fmt.Println("TABLE USERS DOES NOT EXIST, creating");
		re.DB(dbName).TableCreate(UsersTable).RunWrite(session)
	}
	_, errz2 := re.DB(dbName).Table(ListsTable).Run(session)
	if errz2 != nil {
		fmt.Println("TABLE LISTS DOES NOT EXIST, creating");
		re.DB(dbName).TableCreate(ListsTable).RunWrite(session)
		re.DB(dbName).Table(ListsTable).IndexCreate("Owners",re.IndexCreateOpts{Multi:true}).RunWrite(session)
	}

	_, errz3 := re.DB(dbName).Table(ItemsTable).Run(session)
	if errz3 != nil {
		fmt.Println("TABLE ITEMS DOES NOT EXIST, creating");
		re.DB(dbName).TableCreate(ItemsTable).RunWrite(session)
	}

	log.Printf("Database created : %d, with name: %s", resp.DBsCreated, dbName)

  env := &Env{
    DBSession: session,
		DBName: DBName,
		UsersTable: UsersTable,
		ListsTable: ListsTable,
		ItemsTable: ItemsTable,

  }
  return env;
}

func connectToDB() bool {
  sesh, err := re.Connect(re.ConnectOpts{
			Address: fmt.Sprintf("%s:%s",os.Getenv("DB_PORT_28015_TCP_ADDR"), os.Getenv("DB_PORT_28015_TCP_PORT")),
      MaxOpen:  40,
  })
	if err != nil {
    log.Printf(err.Error())
		return false
  }
	session = sesh
	return true
}
