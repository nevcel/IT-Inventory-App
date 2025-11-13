package controllers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"mvc-inventary/models"
)

// InventoryController manages the operations connected to the inventory.
type InventoryController struct {
	Inventory *models.Inventory
}

// Load loads the inventory data from a file.
func (c *InventoryController) Load() {
	file, err := os.Open(models.FileName)
	if os.IsNotExist(err) {
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&c.Inventory.Items); err != nil {
		fmt.Println("Error decoding file:", err)
	}
}

// Save saves the current inventory data to a file.
func (c *InventoryController) Save() {
	file, err := os.Create(models.FileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(c.Inventory.Items); err != nil {
		fmt.Println("Error encoding file:", err)
	}
}

// Add adds a new item to the inventory.
func (c *InventoryController) Add() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter item type: ")
	typeInput, _ := reader.ReadString('\n')
	typeInput = strings.TrimSpace(typeInput)

	fmt.Print("Enter item name: ")
	nameInput, _ := reader.ReadString('\n')
	nameInput = strings.TrimSpace(nameInput)

	fmt.Print("Enter notes: ")
	notesInput, _ := reader.ReadString('\n')
	notesInput = strings.TrimSpace(notesInput)

	id := len(c.Inventory.Items) + 1
	dateAdded := time.Now().Format("2006-01-02")

	newItem := models.Item{
		ID:        id,
		Type:      typeInput,
		Name:      nameInput,
		DateAdded: dateAdded,
		Notes:     notesInput,
	}

	c.Inventory.Items = append(c.Inventory.Items, newItem)
	c.Save()
	fmt.Println("Item added successfully!")
}

// Edit edits an existing item in the inventory.
func (c *InventoryController) Edit() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter item ID to edit: ")
	idInput, _ := reader.ReadString('\n')
	idInput = strings.TrimSpace(idInput)
	id, _ := strconv.Atoi(idInput)

	for i, item := range c.Inventory.Items {
		if item.ID == id && item.DateRemoved == "" {
			fmt.Printf("Editing item: %d | %s | %s\n", item.ID, item.Type, item.Name)

			fmt.Print("Enter new type (leave blank to keep current): ")
			typeInput, _ := reader.ReadString('\n')
			typeInput = strings.TrimSpace(typeInput)
			if typeInput != "" {
				c.Inventory.Items[i].Type = typeInput
			}

			fmt.Print("Enter new name (leave blank to keep current): ")
			nameInput, _ := reader.ReadString('\n')
			nameInput = strings.TrimSpace(nameInput)
			if nameInput != "" {
				c.Inventory.Items[i].Name = nameInput
			}

			fmt.Print("Enter new notes (leave blank to keep current): ")
			notesInput, _ := reader.ReadString('\n')
			notesInput = strings.TrimSpace(notesInput)
			if notesInput != "" {
				c.Inventory.Items[i].Notes = notesInput
			}

			c.Save()
			fmt.Println("Item updated successfully!")
			return
		}
	}
	fmt.Println("Item not found or already removed.")
}

// Remove removes an item from the inventory.
func (c *InventoryController) Remove() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter item ID to remove: ")
	idInput, _ := reader.ReadString('\n')
	idInput = strings.TrimSpace(idInput)
	id, _ := strconv.Atoi(idInput)

	for i, item := range c.Inventory.Items {
		if item.ID == id && item.DateRemoved == "" {
			c.Inventory.Items[i].DateRemoved = time.Now().Format("2006-01-02")
			c.Save()
			fmt.Println("Item removed successfully!")
			return
		}
	}
	fmt.Println("Item not found or already removed.")
}
