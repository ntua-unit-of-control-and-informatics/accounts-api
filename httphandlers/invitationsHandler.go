package httphandlers

import (
	"encoding/json"
	"net/http"

	// "net/smtp"

	email "euclia.xyz/accounts-api/emails"
	"euclia.xyz/accounts-api/models"

	// db "quots/database"
	"strconv"
	"strings"

	// "strconv"
	auth "euclia.xyz/accounts-api/authentication"
	utils "euclia.xyz/accounts-api/utils"
	// "github.com/gorilla/mux"
)

func CreateInvitation(w http.ResponseWriter, r *http.Request) {
	var invit models.Invitation
	json.NewDecoder(r.Body).Decode(&invit)
	apiAll := r.Header.Get("Authorization")
	apiKeyAr := strings.Split(apiAll, " ")
	token := apiKeyAr[1]
	invit.Id = utils.RandStringBytesMaskImprSrcUnsafe(12)
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
			return
		}
	}
	not, err := invitHandler.CreateInvitation(invit)
	go email.SendInvitation(claims.GivenName, invit.EmailTo, invit.Body, invit.InvitationTo)
	json.NewEncoder(w).Encode(not)
}

func GetInvitation(w http.ResponseWriter, r *http.Request) {
	min := r.FormValue("min")
	max := r.FormValue("max")
	email := r.FormValue("email")
	viewed := r.FormValue("viewed")
	var skip int64
	var maxim int64
	var errC error
	if min == "" {
		skip = 0
	} else {
		skip, errC = strconv.ParseInt(min, 10, 64)
	}
	if max == "" {
		maxim = 10
	} else {
		maxim, errC = strconv.ParseInt(max, 10, 64)
	}

	viewedB, _ := strconv.ParseBool(viewed)
	if errC != nil {
		err := models.ErrorReport{
			Message: errC.Error(),
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	invInDb, err := invitHandler.GetByEmail(email, skip, maxim, viewedB)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	// var count int64
	counted, err := invitHandler.CountWithEmail(email, viewedB)
	countstr := strconv.Itoa(int(counted))
	w.Header().Add("total", countstr)
	json.NewEncoder(w).Encode(invInDb)

}

func UpdateInvitation(w http.ResponseWriter, r *http.Request) {
	answer := r.FormValue("answer")
	var invit models.Invitation
	json.NewDecoder(r.Body).Decode(&invit)
	if answer == "accept" {
		orgId := invit.InvitationToId
		org, err := orgHandler.GetOrgById(orgId)
		user, err := usersHandler.FindUserByEmail(invit.EmailTo)
		if err != nil {
			err := models.ErrorReport{
				Message: err.Error(),
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}
		if !utils.Contains(org.Users, user.Id) {
			org.Users = append(org.Users, user.Id)
			org, err = orgHandler.UpdateOrg(org)
		}
		if !utils.Contains(user.Organizations, org.Id) {
			user.Organizations = append(user.Organizations, org.Id)
			user, err = usersHandler.UpdateUser(user)
		}
		if err != nil {
			err := models.ErrorReport{
				Message: err.Error(),
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}
		invit.Viewed = true
		invit.Action = answer
		invit, err = invitHandler.UpdateInvitation(invit)
		json.NewEncoder(w).Encode(invit)
	}
	if answer == "decline" {
		invit.Action = answer
		invit.Viewed = true
		var err error
		invit, err = invitHandler.UpdateInvitation(invit)
		if err != nil {
			err := models.ErrorReport{
				Message: err.Error(),
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}
		json.NewEncoder(w).Encode(invit)
	}
}
