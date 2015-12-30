package models

type ShoppingList struct {
  Id      string `gorethink:"id,omitempty"`
	Items   []Item
  Owner   string `json:"Owner" binding:"required"`
}

func NewShoppingList() *ShoppingList {
	return &ShoppingList{
		Items: make([]Item, 0),
	}
}
