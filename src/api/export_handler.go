package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/axentx/social-assembly/pkg/config"
	"github.com/axentx/social-assembly/pkg/db"
	"github.com/axentx/social-assembly/pkg/encryption"
	"github.com/gorilla/mux"
)

// ExportHandler handles GDPR-compliant user data export requests.
func ExportHandler(w http.ResponseWriter, r *http.Request) {
	// Validate user ID (alphanumeric only)
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if !isValidUserID(userID) {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	// Fetch user data
	userData, err := db.GetUser(userID)
	if err != nil {
		if err == db.ErrUserNotFound {
			http.Error(w, "user not found", http.StatusNotFound)
		} else {
			http.Error(w, "failed to fetch user data", http.StatusInternalServerError)
		}
		return
	}

	// Encrypt sensitive data
	encryptedData, err := encryption.EncryptUserData(userData)
	if err != nil {
		http.Error(w, "failed to encrypt data", http.StatusInternalServerError)
		return
	}

	// Return encrypted data (GDPR-compliant JSON)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(encryptedData)
}

// isValidUserID checks if the user ID is alphanumeric (basic sanitization).
func isValidUserID(userID string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, userID)
	return matched
}

func main() {
	// Load config (AES key from env or secure defaults)
	if err := config.LoadConfig(); err != nil {
		panic("failed to load config: " + err.Error())
	}

	// Initialize router
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/export/{user_id}", ExportHandler).Methods("GET")

	// Start server
	http.ListenAndServe(":8080", r)
}