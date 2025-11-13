package main

import (
	"mvc-inventary/controllers"
	"mvc-inventary/models"
	"mvc-inventary/views/console"
)

func main() {
	inventory := &controllers.InventoryController{
		Inventory: &models.Inventory{},
	}
	inventory.Load()
	console.MainMenu(inventory)
}
