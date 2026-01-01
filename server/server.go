package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"crypto/rand"
	"encoding/hex"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"mvc-inventary/controllers"
	"mvc-inventary/models"
)

var inventory *controllers.InventoryController
var users *controllers.UserController
var sessions = map[string]int{} // sessionID -> userID

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// StartServer startet den HTTP-Server
func StartServer() {
	inventory = &controllers.InventoryController{
		Inventory: &models.Inventory{},
	}
	inventory.Load()

	users = &controllers.UserController{}
	users.Load()

	r := mux.NewRouter()

	// Routen definieren
	r.HandleFunc("/inventory", GetInventory).Methods("GET")
	r.HandleFunc("/inventory/{id}", GetItem).Methods("GET")
	r.HandleFunc("/inventory", AddItem).Methods("POST")
	r.HandleFunc("/inventory/{id}", EditItem).Methods("PUT")
	r.HandleFunc("/inventory/{id}", RemoveItem).Methods("DELETE")
	r.HandleFunc("/items/search", SearchItems).Methods("GET")
	r.HandleFunc("/login", Login).Methods("POST")
	r.HandleFunc("/logout", Logout).Methods("POST")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./views/static/"))))
	r.Handle("/", http.FileServer(http.Dir("./views/static/")))

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

	w.Header().Set("Content-Type", "application/json")
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

// SearchItems sucht Atrikel in der Liste
func SearchItems(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	query = strings.ToLower(query)

	var result []models.Item

	for _, item := range inventory.Inventory.Items {
		if strings.Contains(strings.ToLower(item.Name), query) {
			result = append(result, item)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Helper Session-ID generieren
func newSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Login Handler
func Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungueltiges JSON", http.StatusBadRequest)
		return
	}

	user := users.FindByUsername(req.Username)
	if user == nil {
		http.Error(w, "Login fehlgeschlagen", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Login fehlgeschlagen", http.StatusUnauthorized)
		return
	}

	sid, err := newSessionID()
	if err != nil {
		http.Error(w, "Serverfehler", http.StatusInternalServerError)
		return
	}

	sessions[sid] = user.ID

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		// Secure: true, // erst aktivieren wenn HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"username": user.Username,
		"role":     user.Role,
	})
}

// Logout Handler
func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		delete(sessions, cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusNoContent)
}
