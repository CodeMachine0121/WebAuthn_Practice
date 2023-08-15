package Utils

import (
	"encoding/json"
	"github.com/duo-labs/webauthn/protocol"
	"log"
	"net/http"
)

func ErrorHandle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func JSONResponse(w http.ResponseWriter, data *protocol.CredentialCreation, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}
