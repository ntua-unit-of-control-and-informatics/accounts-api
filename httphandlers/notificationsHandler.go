package httphandlers

import (
	"encoding/json"
	"net/http"

	// "net/smtp"
	"euclia.xyz/accounts-api/models"
	// db "quots/database"
	"strings"
	// "strconv"
	auth "euclia.xyz/accounts-api/authentication"
	// "github.com/gorilla/mux"
)

func CreateInvitationNotification(w http.ResponseWriter, r *http.Request) {
	var notif models.Notification
	json.NewDecoder(r.Body).Decode(&notif)
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
		return
	}
	_, orgs, err := orgHandler.GetOrgByAdmin(claims.Subject)
	for _, orga := range orgs {
		i := indexOf(claims.Subject, orga.Admins)
		if i == -1 {
			err := models.ErrorReport{
				Message: "Cannot Invite user. Only Admins can",
				Status:  http.StatusForbidden,
			}
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(err)
		}
		return
	}
	not, err := notHandler.CreateNotification(notif)
	// go email.SendInvitation(not.To, claims.Name, notif.InvitationTo)
	json.NewEncoder(w).Encode(not)
}

// func GetInvitation(w http.ResponseWriter, r *http.Request) {
// 	email := r.FormValue("email")
// 	// id := r.FormValue("_id")
// 	apiAll := r.Header.Get("Authorization")
// 	apiKeyAr := strings.Split(apiAll, " ")
// 	token := apiKeyAr[1]
// 	claims, err := auth.GetClaims(token)
// 	if err != nil {
// 		err := models.ErrorReport{
// 			Message: err.Error(),
// 			Status:  http.StatusBadRequest,
// 		}
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(err)
// 		return
// 	}
// 	_, orgs, err := orgHandler.GetOrgByAdmin(claims.Subject)
// 	for _, orga := range orgs {
// 		i := indexOf(claims.Subject, orga.Admins)
// 		if i == -1 {
// 			err := models.ErrorReport{
// 				Message: "Cannot Invite user. Only Admins can",
// 				Status:  http.StatusForbidden,
// 			}
// 			w.WriteHeader(http.StatusForbidden)
// 			json.NewEncoder(w).Encode(err)
// 		}
// 		return
// 	}
// 	not, err := notHandler.CreateNotification(notif)
// 	// go email.SendInvitation(not.To, claims.Name, notif.InvitationTo)
// 	json.NewEncoder(w).Encode(not)
// }
