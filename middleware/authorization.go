package middleware

import (
    "net/http"
    "log"
    jwt "github.com/dgrijalva/jwt-go"
)

func checkAuthorization(h HandlerMethods, w http.ResponseWriter, r *http.Request) (interface{}, bool) {
  var null interface{} = nil
  if h.shallAuthorize() {
    _ = "breakpoint"
    token, err := jwt.ParseFromRequest(r, keyFn2)
    if err == nil && token.Valid {
      return token.Claims["id"], false
    } else {
      log.Printf("HTTP %d - %s", http.StatusUnauthorized, token)
      http.Error(w, http.StatusText(http.StatusUnauthorized),
        http.StatusUnauthorized)
      return null, true
    }
  }
  return null, false
}

func keyFn2(token *jwt.Token) (interface{}, error) {
    return []byte("mysupersecretkey"), nil
}
