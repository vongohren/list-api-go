package models

import(
  "time"
)

type List struct {
  Id      string `gorethink:"id,omitempty"`
	Items   []string
  Owners   []string `json:"Owners" binding:"required"`
  Titel   string `json:"Title" binding:"required"`
  Updated time.Time
}
