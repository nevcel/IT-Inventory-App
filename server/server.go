package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"mvc-inventary/controllers"
	"mvc-inventary/models"
)

var inventory *controllers.InventoryController

// StartServer startet den HTTP-Server
func StartServer() {
	inventory = &controllers.InventoryController{
		Inventory: &models.Inventory{},
	}
	inventory.Load()

	r := mux.NewRouter()

	// Routen definieren
	r.HandleFunc("/inventory", GetInventory).Methods("GET")
	r.HandleFunc("/inventory/{id}", GetItem).Methods("GET")
	r.HandleFunc("/inventory", AddItem).Methods("POST")
	r.HandleFunc("/inventory/{id}", EditItem).Methods("PUT")
	r.HandleFunc("/inventory/{id}", RemoveItem).Methods("DELETE")

	http.ListenAndServe(":8080", r)
}

// --- API Handler ---

// GetInventory gibt alle Items zurück
func GetInventory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventory.Inventory.Items)
}

// GetItem gibt ein einzelnes Item zurück
func GetItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	for _, item := range inventory.Inventory.Items {
		if strconv.Itoa(item.ID) == id {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w, "Item not found", http.StatusNotFound)
}

// AddItem fügt ein neues Item hinzu
func AddItem(w http.ResponseWriter, r *http.Request) {
	var newItem models.Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newItem.ID = len(inventory.Inventory.Items) + 1
	newItem.DateAdded = time.Now().Format("2006-01-02")
	inventory.Inventory.Items = append(inventory.Inventory.Items, newItem)
	inventory.Save()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
}

// EditItem bearbeitet ein bestehendes Item
func EditItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updatedItem models.Item
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, item := range inventory.Inventory.Items {
		if strconv.Itoa(item.ID) == id && item.DateRemoved == "" {
			if updatedItem.Type != "" {
				inventory.Inventory.Items[i].Type = updatedItem.Type
			}
			if updatedItem.Name != "" {
				inventory.Inventory.Items[i].Name = updatedItem.Name
			}
			if updatedItem.Notes != "" {
				inventory.Inventory.Items[i].Notes = updatedItem.Notes
			}
			inventory.Save()
			json.NewEncoder(w).Encode(inventory.Inventory.Items[i])
			return
		}
	}
	http.Error(w, "Item not found or removed", http.StatusNotFound)
}

// RemoveItem markiert ein Item als entfernt
func RemoveItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for i, item := range inventory.Inventory.Items {
		if strconv.Itoa(item.ID) == id {
			// Entferne das Item aus der Slice
			inventory.Inventory.Items = append(inventory.Inventory.Items[:i], inventory.Inventory.Items[i+1:]...)
			inventory.Save()
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Item deleted"})
			return
		}
	}
	http.Error(w, "Item not found", http.StatusNotFound)
}
