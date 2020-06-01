package httphandlers

import (
	"encoding/json"
	"net/http"

	// "net/smtp"
	email "euclia.xyz/accounts-api/emails"
	"euclia.xyz/accounts-api/models"

	// db "quots/database"
	"strings"
	// "strconv"
	auth "euclia.xyz/accounts-api/authentication"
	// "github.com/gorilla/mux"
)

func CreateRequest(w http.ResponseWriter, r *http.Request) {
	var req models.Request
	json.NewDecoder(r.Body).Decode(&req)
	apiAll := r.Header.Get("Authorization")
	apiKeyAr := strings.Split(apiAll, " ")
	token := apiKeyAr[1]
	claims, err := auth.GetClaims(token)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
	}
	if claims.Subject != req.FromUser {
		err := models.ErrorReport{
			Message: "You cannot request for creadits on behalf of other user",
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
	}
	reqN, err := reqHandler.CreateRequest(req)
	go email.SendCreditRequest(req.FromUser, req.UserEmail, req.Message)
	json.NewEncoder(w).Encode(reqN)
}
