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
		UserTable 	string
		ListsTable 	string
		ListKey			string
}

var (
	DBName string = "list_api"
	UserTable string = "users"
	ListsTable string = "lists"
	ListKey string = "Owner"
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
	_, errz := re.DB(dbName).Table("users").Run(session)
	if errz != nil {
		fmt.Println("TABLE USERS DOES NOT EXIST, creating");
		re.DB(dbName).TableCreate("users").RunWrite(session)
	}
	_, errz2 := re.DB(dbName).Table("lists").Run(session)
	if errz2 != nil {
		fmt.Println("TABLE ITEMS DOES NOT EXIST, creating");
		re.DB(dbName).TableCreate("lists").RunWrite(session)
	}
	log.Printf("Database created : %d, with name: %s", resp.DBsCreated, dbName)

  env := &Env{
    DBSession: session,
		DBName: DBName,
		UserTable: UserTable,
		ListsTable: ListsTable,

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
