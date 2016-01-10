package models

import(
  "time"
)

type Item struct {
  Id        string `gorethink:"id,omitempty"`
	Text      string `json:"Text"`
	Done      string `json:"Done,omitempty"`
	Added     time.Time
}
