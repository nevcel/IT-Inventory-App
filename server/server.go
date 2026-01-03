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

type registerRequest struct {
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
	r.HandleFunc("/inventory", RequireAuth(GetInventory)).Methods("GET")
	r.HandleFunc("/inventory/{id}", RequireAuth(GetItem)).Methods("GET")
	r.HandleFunc("/inventory", RequireAuth(AddItem)).Methods("POST")
	r.HandleFunc("/inventory/{id}", RequireAuth(EditItem)).Methods("PUT")
	r.HandleFunc("/inventory/{id}", RequireAuth(RemoveItem)).Methods("DELETE")
	r.HandleFunc("/items/search", SearchItems).Methods("GET")
	r.HandleFunc("/login", Login).Methods("POST")
	r.HandleFunc("/logout", Logout).Methods("POST")
	r.HandleFunc("/register", Register).Methods("POST")

	//Test Router für Login:
	r.HandleFunc("/me", RequireAuth(Me)).Methods("GET")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./views/static/"))))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./views/static/")))

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
	u := GetCurrentUser(r)
	if u == nil {
		http.Error(w, "Nicht eingeloggt", http.StatusUnauthorized)
		return
	}

	items := inventory.Inventory.Items

	if u.Role != "admin" {
		filtered := []models.Item{}
		for _, it := range items {
			if it.OwnerID == u.ID {
				filtered = append(filtered, it)
			}
		}
		items = filtered
	}

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
	u := GetCurrentUser(r)
	if u == nil {
		http.Error(w, "Nicht eingeloggt", http.StatusUnauthorized)
		return
	}

	newItem.OwnerID = u.ID
	newItem.ID = len(inventory.Inventory.Items) + 1
	newItem.DateAdded = time.Now().Format("2006-01-02")
	newItem.DateEdited = newItem.DateAdded
	inventory.Inventory.Items = append(inventory.Inventory.Items, newItem)
	inventory.Save()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
}

// EditItem bearbeitet ein bestehendes Item
func EditItem(w http.ResponseWriter, r *http.Request) {
	// Current User holen
	u := GetCurrentUser(r)
	if u == nil {
		http.Error(w, "Nicht eingeloggt", http.StatusUnauthorized)
		return
	}

	// ID aus URL holen
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Ungueltige ID", http.StatusBadRequest)
		return
	}

	// Body einlesen (nur Felder, die editierbar sind)
	type updateRequest struct {
		Type  string `json:"type"`
		Name  string `json:"name"`
		Notes string `json:"notes"`
	}

	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungueltiges JSON", http.StatusBadRequest)
		return
	}

	// Item suchen
	for i := range inventory.Inventory.Items {
		if inventory.Inventory.Items[i].ID == id {

			// Berechtigung pruefen: Admin darf alles, User nur eigene Items
			if u.Role != "admin" && inventory.Inventory.Items[i].OwnerID != u.ID {
				http.Error(w, "Keine Berechtigung", http.StatusForbidden)
				return
			}

			// Felder aktualisieren
			inventory.Inventory.Items[i].Type = req.Type
			inventory.Inventory.Items[i].Name = req.Name
			inventory.Inventory.Items[i].Notes = req.Notes

			// date_edited setzen
			inventory.Inventory.Items[i].DateEdited = time.Now().Format("2006-01-02")

			// Speichern
			inventory.Save()

			// Antwort
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(inventory.Inventory.Items[i])
			return
		}
	}

	http.Error(w, "Item nicht gefunden", http.StatusNotFound)
}

// RemoveItem markiert ein Item als entfernt
func RemoveItem(w http.ResponseWriter, r *http.Request) {
	// Current User holen
	u := GetCurrentUser(r)
	if u == nil {
		http.Error(w, "Nicht eingeloggt", http.StatusUnauthorized)
		return
	}

	// ID aus URL holen und parsen
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Ungueltige ID", http.StatusBadRequest)
		return
	}

	// Item suchen (Index + Item)
	idx := -1
	var found models.Item

	for i, it := range inventory.Inventory.Items {
		if it.ID == id {
			idx = i
			found = it
			break
		}
	}

	if idx == -1 {
		http.Error(w, "Item nicht gefunden", http.StatusNotFound)
		return
	}

	// Berechtigung pruefen
	if u.Role != "admin" && found.OwnerID != u.ID {
		http.Error(w, "Keine Berechtigung", http.StatusForbidden)
		return
	}

	// Loeschen
	inventory.Inventory.Items = append(inventory.Inventory.Items[:idx], inventory.Inventory.Items[idx+1:]...)
	inventory.Save()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item deleted"})
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

// Register: Handler zum registrieren von User
func Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungueltiges JSON", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username und Passwort sind erforderlich", http.StatusBadRequest)
		return
	}

	// Username bereits vergeben?
	if users.FindByUsername(req.Username) != nil {
		http.Error(w, "Username bereits vergeben", http.StatusConflict)
		return
	}

	// Passwort hashen
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Serverfehler", http.StatusInternalServerError)
		return
	}

	// Neue ID bestimmen (max+1)
	nextID := 1
	for _, u := range users.Users {
		if u.ID >= nextID {
			nextID = u.ID + 1
		}
	}

	newUser := models.User{
		ID:           nextID,
		Username:     req.Username,
		PasswordHash: string(hash),
		Role:         "user",
	}

	users.Users = append(users.Users, newUser)
	users.Save()

	w.WriteHeader(http.StatusCreated)
}

// Test Handler für Login
func Me(w http.ResponseWriter, r *http.Request) {
	u := GetCurrentUser(r)
	if u == nil {
		http.Error(w, "Nicht eingeloggt", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":       u.ID,
		"username": u.Username,
		"role":     u.Role,
	})
}

// --- Helper ---

// Helper: Current User hole
func GetCurrentUser(r *http.Request) *models.User {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil
	}

	userID, ok := sessions[cookie.Value]
	if !ok {
		return nil
	}

	// User anhand ID finden
	for i := range users.Users {
		if users.Users[i].ID == userID {
			return &users.Users[i]
		}
	}

	return nil
}

// Helper: Login erforderlich
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := GetCurrentUser(r)
		if u == nil {
			http.Error(w, "Nicht eingeloggt", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// Helper: Admin erforderlich
func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := GetCurrentUser(r)
		if u == nil {
			http.Error(w, "Nicht eingeloggt", http.StatusUnauthorized)
			return
		}
		if u.Role != "admin" {
			http.Error(w, "Keine Berechtigung", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
