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
  "html/template"
  re "github.com/dancannon/gorethink"
  jwt "github.com/dgrijalva/jwt-go"
)

func init () {
  goth.UseProviders(
    facebook.New(os.Getenv("FACEBOOK_SECRET"), os.Getenv("FACEBOOK_APP_SECRET"), fmt.Sprintf("http://%s/auth/facebook/callback", os.Getenv("DNS_HOSTNAME"))),
  )
  gothic.GetProviderName = getProviderName
}

type CookieUrl struct {
  Cookie string
  Url string
}

func BeginAuthHandler(env *db.Env,  w http.ResponseWriter, r *http.Request) error {
  gothic.BeginAuthHandler(w, r)
  name := gothic.SessionName
  fmt.Println(name)
  _ = "breakpoint"
  sesh, erhor := gothic.Store.Get(r, gothic.SessionName)
  fmt.Println(erhor)
  ertor := sesh.Save(r,w)
  fmt.Println(ertor)
  cookie := w.Header().Get("Set-Cookie")
  cookie = strings.TrimPrefix(strings.Split(cookie, ";")[0], fmt.Sprintf("%s=", gothic.SessionName))
  cookieUrl := CookieUrl{cookie, r.Header.Get("Referer")}
  _, err4 := re.DB("list_api").Table("auth_sessions").Insert(cookieUrl).RunWrite(env.DBSession)
  if err4 != nil {
      return err4
  }
  log.Printf(gothic.GetState(r))
  return nil;
}

func CallBack(env *db.Env, w http.ResponseWriter, r *http.Request) error {
	// print our state string to the console. Ideally, you should verify
	// that it's the same string as the one you set in `setState`
  cookie, _ := r.Cookie(gothic.SessionName)
  fmt.Println(cookie)
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
  log.Printf(jwtString)
  if err5 != nil {
      return err5
  }
  sessions := []CookieUrl{}
  seshz, errSesh := re.DB("list_api").Table("auth_sessions").Filter(map[string]interface{}{
    "Cookie":cookie.Value,
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

  tmpl := fmt.Sprintf("templates/successAuth.html")
  t, err := template.ParseFiles(tmpl)
  if err != nil {
      log.Print("template parsing error: ", err)
  }
  err = t.Execute(w, jsonToken)
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
