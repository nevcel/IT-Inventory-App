package models

// Item represents an inventory item with various attributes.
type Item struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	DateAdded   string `json:"date_added"`
	DateRemoved string `json:"date_removed"`
	Notes       string `json:"notes"`
}

// Inventory represents a collection of items.
type Inventory struct {
	Items []Item
}

// FileName is the name of the file where the inventory data is stored.
const FileName = "inventory.json"
