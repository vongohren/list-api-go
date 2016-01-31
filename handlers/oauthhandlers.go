package handlers

import (
  "log"
  "fmt"
  "os"
  "strings"
  "errors"
  "net/http"
  "time"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
  "github.com/Snorlock/shoppingApi/db"
  "github.com/Snorlock/mux"
  "github.com/pborman/uuid"
  "text/template"
  re "github.com/dancannon/gorethink"
  jwt "github.com/dgrijalva/jwt-go"
)

func init () {
  goth.UseProviders(
    facebook.New(os.Getenv("FACEBOOK_SECRET"), os.Getenv("FACEBOOK_APP_SECRET"), fmt.Sprintf("http://%s/auth/facebook/callback", os.Getenv("DNS_HOSTNAME"))),
  )
  gothic.GetProviderName = getProviderName
}

type Identifier struct {
  UUID string
  Url string
}

func BeginAuthHandler(env *db.Env,  w http.ResponseWriter, r *http.Request) error {
  gothic.SetState = func(req *http.Request) string {
    _ = "breakpoint"
    state := r.Header.Get("state")
	  return state
  }
  url := r.URL.Query().Get("url")
  uuid := uuid.New();

  r.Header.Set("state", fmt.Sprintf("%s!%s", url, uuid))
  gothic.BeginAuthHandler(w, r)

  identifier := Identifier{uuid, url}
  _, err4 := re.DB("list_api").Table(env.AuthSessionTable).Insert(identifier).RunWrite(env.DBSession)
  if err4 != nil {
      return err4
  }
  return nil;
}

func CallBack(env *db.Env, w http.ResponseWriter, r *http.Request) error {
	// print our state string to the console. Ideally, you should verify
	// that it's the same string as the one you set in `setState`
  _ = "breakpoint"
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		return err
	}
  users, err2 := re.DB("list_api").Table("users").Filter(map[string]interface{}{
    "Email":user.Email,
  }).Run(env.DBSession)
  if err2 != nil {
    return err2
  }
  var usersArr []interface{}
  err3 := users.All(&usersArr)
  if err3 != nil {
      return err3
  }
  if len(usersArr) > 0 {
    log.Printf("User exists")
  } else {
    log.Printf("User does not exist, creating")
    _, err4 := re.DB("list_api").Table("users").Insert(user).RunWrite(env.DBSession)
    if err4 != nil {
        return err4
    }
  }

  token := jwt.New(jwt.SigningMethodHS256)
  token.Claims["id"] = user.Email
  token.Claims["iat"] = time.Now().Unix()
  token.Claims["exp"] = time.Now().Add(time.Second * 3600 * 24).Unix()
  jwtString, err5 := token.SignedString([]byte("mysupersecretkey"))
  if err5 != nil {
      return err5
  }

  state := gothic.GetState(r)
  uuid := strings.Split(state, "!")[1]
  sessions := []Identifier{}
  seshz, errSesh := re.DB("list_api").Table(env.AuthSessionTable).Filter(map[string]interface{}{
    "UUID":uuid,
  }).Run(env.DBSession)

  if errSesh != nil {
    return errSesh
  }
  errArr := seshz.All(&sessions)
  if errArr != nil {
      return errArr
  }
  referer := ""
  if len(sessions) > 0 {
    //read out the url this session belongs to
    log.Printf("SeshExists")
    referer = sessions[0].Url
  } else {
    //say something went wrong
    log.Printf("Sesh does not exist")
  }

  jsonToken := Token{jwtString, referer}
  pwd, _ := os.Getwd()
  fmt.Println(pwd)
  tmpl := getHtmlReponse()
  err = tmpl.Execute(w, jsonToken)
  if err != nil {
      log.Print("template executing error: ", err)
  }

  return nil
}

type Token struct {
  Bearer string
  Url string
}

func getProviderName(req *http.Request) (string, error) {
	provider := req.URL.Query().Get("provider")
	if provider == "" {
		params := mux.Vars(req)
		provider = params["provider"]
	}
	if provider == "" {
		return provider, errors.New("you must select a provider")
	}
	return provider, nil
}

func getHtmlReponse() *template.Template {
  return template.Must(template.New("html").Parse(`
    <!DOCTYPE html>
    <html lang="en">
      <head>
        <meta charset="utf-8">
        <title>Golang Simple HTML</title>
      </head>
      <body>
        <h1>Success Auth</h1>
        <div>{{.Bearer}}</div>
      </body>
      <script>
        window.opener.postMessage({{.Bearer}},{{.Url}});
        self.close()
      </script>
    </html>
    `))
}
