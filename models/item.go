package models

import(
  "time"
)

type Item struct {
	Item      string `json:"Item"`
	Accuired  bool `json:"Accuired,omitempty"`
	Added     time.Time
}
