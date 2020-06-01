package httphandlers

import (
	"encoding/json"
	"net/http"

	// db "quots/database"
	"euclia.xyz/accounts-api/models"

	"strings"

	"strconv"

	auth "euclia.xyz/accounts-api/authentication"
	utils "euclia.xyz/accounts-api/utils"
	"github.com/gorilla/mux"
)

func CreateOrganization(w http.ResponseWriter, r *http.Request) {
	var org models.Organization
	var userF models.User
	var userOrgs []string
	var orgUsers []string
	json.NewDecoder(r.Body).Decode(&org)
	org.Id = utils.RandStringBytesMaskImprSrcUnsafe(8)
	apiAll := r.Header.Get("Authorization")
	apiKeyAr := strings.Split(apiAll, " ")
	token := apiKeyAr[1]
	claims, err := auth.GetClaims(token)
	org.Creator = claims.Subject
	org.Admins = append(org.Admins, claims.Subject)
	orgUsers = append(orgUsers, claims.Subject)
	org.Users = orgUsers
	org, err = orgHandler.CreateOrganization(org)
	userF, err = usersHandler.GetUserById(claims.Subject)
	userOrgs = userF.Organizations
	userOrgs = append(userOrgs, org.Id)
	userF.Organizations = userOrgs
	userF, err = usersHandler.UpdateUserOrganizations(userF)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(org)
}

func UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	var org models.Organization
	json.NewDecoder(r.Body).Decode(&org)
	apiAll := r.Header.Get("Authorization")
	apiKeyAr := strings.Split(apiAll, " ")
	token := apiKeyAr[1]
	claims, err := auth.GetClaims(token)
	org.Creator = claims.Subject
	orgInDb, err := orgHandler.GetOrgById(org.Id)
	if indexOf(claims.Subject, orgInDb.Admins) == -1 {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusForbidden,
		}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(err)
		return
	}
	if len(org.Users) == 0 {
		org.Users = append(org.Users, org.Creator)
	}
	org, err = orgHandler.UpdateOrg(org)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(org)
}

func GetOrganizations(w http.ResponseWriter, r *http.Request) {
	apiAll := r.Header.Get("Authorization")
	apiKeyAr := strings.Split(apiAll, " ")
	token := apiKeyAr[1]
	claims, err := auth.GetClaims(token)
	count, orgsInDb, err := orgHandler.GetOrgByAdmin(claims.Subject)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	countstr := strconv.Itoa(count)
	w.Header().Add("total", countstr)
	json.NewEncoder(w).Encode(orgsInDb)
}

func GetOrganizationById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	apiAll := r.Header.Get("Authorization")
	apiKeyAr := strings.Split(apiAll, " ")
	token := apiKeyAr[1]
	_, err := auth.GetClaims(token)
	orgsInDb, err := orgHandler.GetOrgById(id)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	// countstr := strconv.Itoa(count)
	// w.Header().Add("total", countstr)
	json.NewEncoder(w).Encode(orgsInDb)
}

func DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
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
	org, err := orgHandler.GetOrgById(id)
	if claims.Subject != org.Creator {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	} else if claims.Subject == org.Creator {
		orgsDeleted, err := orgHandler.DeleteOrg(org.Id)
		_, users, err := usersHandler.FindOrgsUsers(org.Id)
		for _, user := range users {
			i := indexOf(org.Id, user.Organizations)
			user.Organizations[i] = user.Organizations[len(user.Organizations)-1] // Copy last element to index i.
			user.Organizations[len(user.Organizations)-1] = ""                    // Erase last element (write zero value).
			user.Organizations = user.Organizations[:len(user.Organizations)-1]
			var userToUp models.User
			userToUp.Organizations = user.Organizations
			userToUp.Id = user.Id
			usersHandler.UpdateUserOrganizations(userToUp)
		}
		if err != nil {
			err := models.ErrorReport{
				Message: err.Error(),
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		} else {
			json.NewEncoder(w).Encode(orgsDeleted)
		}
	}

}
