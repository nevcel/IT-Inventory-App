package console

import (
	"bufio"
	"fmt"
	"mvc-inventary/controllers"
	"os"
	"strings"
)

// FileName is the name of the file where the inventory data is stored.
func ShowInventory(inv *controllers.InventoryController) {
	fmt.Println("\nCurrent Inventory:")
	fmt.Println("ID  | Type            | Name            | Date Added      | Notes")
	for _, item := range inv.Inventory.Items {
		if item.DateRemoved == "" {
			fmt.Printf("%-3d | %-15s | %-15s | %-15s | %-15s\n", item.ID, item.Type, item.Name, item.DateAdded, item.Notes)
		}
	}
}

// MainMenu displays the main menu and handles user input for various inventory operations.
func MainMenu(inv *controllers.InventoryController) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n--- IT Inventory Management ---")
		fmt.Println("1. Show inventory")
		fmt.Println("2. Add an item")
		fmt.Println("3. Edit an item")
		fmt.Println("4. Remove an item")
		fmt.Println("5. Exit")
		fmt.Print("Choose an option: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			ShowInventory(inv)
		case "2":
			inv.Add()
		case "3":
			inv.Edit()
		case "4":
			inv.Remove()
		case "5":
			fmt.Println("Exiting application.")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
