package models

import(
  "time"
)

type List struct {
  Id      string `gorethink:"id,omitempty"`
	Items   []string
  Owners   []string `json:"Owners" binding:"required"`
  Title   string `json:"Title" binding:"required"`
  Updated time.Time
}

type DetailedList struct {
  Id      string `gorethink:"id,omitempty"`
	Items   []Item
  Owners   []string `json:"Owners" binding:"required"`
  Title   string `json:"Title" binding:"required"`
  Updated time.Time
}
